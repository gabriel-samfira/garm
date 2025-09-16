package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	garmWs "github.com/cloudbase/garm-provider-common/util/websocket"
	"github.com/cloudbase/garm/cmd/garm-agent/config"
	"github.com/cloudbase/garm/workers/websocket/agent/messaging"
	"github.com/gorilla/websocket"
)

var closed = make(chan struct{})

func init() { close(closed) }

func NewService(ctx context.Context, cfg *config.Agent) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate agent config: %w", err)
	}

	return &Service{
		ctx:      ctx,
		cfg:      cfg,
		done:     closed,
		sessions: make(map[string]*ShellSession),
	}, nil
}

type Service struct {
	ctx context.Context
	cfg *config.Agent
	cli *garmWs.Reader

	mux     sync.Mutex
	running bool
	done    chan struct{}

	sessions map[string]*ShellSession
}

func (s *Service) Done() chan struct{} {
	return s.done
}

func (s *Service) writeMessage(msg []byte) error {
	if err := s.cli.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		return err
	}
	return nil
}

func (s *Service) handleMessage(msgType int, msg []byte) (err error) {
	if msgType != websocket.BinaryMessage && msgType != websocket.TextMessage {
		slog.InfoContext(s.ctx, "ignoring invalid message type", "message_type", msgType)
		return nil
	}

	agentMsg, err := messaging.UnmarshalAgentMessage(msg)
	if err != nil {
		slog.ErrorContext(s.ctx, "failed to unmarshal agent message", "error", err)
		return fmt.Errorf("failed to unmarshal agent message")
	}

	slog.InfoContext(s.ctx, "handling message", "message_type", agentMsg.Type)
	switch agentMsg.Type {
	case messaging.MessageTypeCreateShell:
		slog.InfoContext(s.ctx, "handling create shell message")
		createShell, err := messaging.Unmarshal[messaging.CreateShellMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshall create shell message: %w", err)
		}
		defer func() {
			if err != nil {
				shellReadyMsg := messaging.ShellReadyMessage{
					SessionID: createShell.SessionID,
					IsError:   1,
					Message:   []byte(fmt.Sprintf("failed to create shell: %q", err)),
				}
				if innerErr := s.writeMessage(shellReadyMsg.Marshal()); innerErr != nil {
					slog.ErrorContext(s.ctx, "failed to send error message", "error", innerErr)
				}
			}
		}()

		sessionID := createShell.ID()
		if sessionID == "" {
			return fmt.Errorf("failed to parse session ID")
		}
		s.mux.Lock()
		if _, ok := s.sessions[sessionID]; ok {
			s.mux.Unlock()
			return fmt.Errorf("session ID %s already exists", sessionID)
		}
		session, err := NewShellSession(s.ctx, createShell, s.writeMessage)
		if err != nil {
			s.mux.Unlock()
			return fmt.Errorf("failed to create session: %w", err)
		}
		if err := session.Start(); err != nil {
			s.mux.Unlock()
			return fmt.Errorf("failed to start session: %w", err)
		}
		s.sessions[sessionID] = session
		go func(sessionID string) {
			select {
			case <-s.ctx.Done():
			case <-s.done:
			case <-session.Done():
			}
			s.mux.Lock()
			delete(s.sessions, sessionID)
			s.mux.Unlock()
		}(sessionID)
		s.mux.Unlock()
	case messaging.MessageTypeShellResize:
		resizeMsg, err := messaging.Unmarshal[messaging.ShellResizeMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshall shell resize message: %w", err)
		}

		s.mux.Lock()
		session, ok := s.sessions[resizeMsg.ID()]
		if !ok {
			s.mux.Unlock()
			return nil
		}
		if err := session.shell.Resize(resizeMsg.Cols, resizeMsg.Rows); err != nil {
			s.mux.Unlock()
			return fmt.Errorf("failed to resize shell: %w", err)
		}
		s.mux.Unlock()
	case messaging.MessageTypeClientShellClosed:
		closedMsg, err := messaging.Unmarshal[messaging.ClientShellClosedMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshall shell closed message: %w", err)
		}
		s.mux.Lock()
		session, ok := s.sessions[closedMsg.ID()]
		if !ok {
			s.mux.Unlock()
			return nil
		}
		if err := session.Stop(); err != nil {
			s.mux.Unlock()
			return fmt.Errorf("failed to close session: %w", err)
		}
		s.mux.Unlock()
	case messaging.MessageTypeShellData:
		slog.InfoContext(s.ctx, "received shell data message")
		shellData, err := messaging.Unmarshal[messaging.ShellDataMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshall shell data message: %w", err)
		}
		s.mux.Lock()
		session, ok := s.sessions[shellData.ID()]
		if !ok {
			s.mux.Unlock()
			return nil
		}
		slog.InfoContext(s.ctx, "found shell session", "session_id", shellData.ID())
		if _, err := session.shell.Write(shellData.Data); err != nil {
			slog.ErrorContext(s.ctx, "failed to write shell data; stopping session", "error", err, "session_id", shellData.ID())
			if err := session.Stop(); err != nil {
				s.mux.Unlock()
				return fmt.Errorf("failed to stop session %s", shellData.ID())
			}
		}
		s.mux.Unlock()
	}
	return nil
}

func (s *Service) Start() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.running {
		return nil
	}

	cli, err := garmWs.NewReader(s.ctx, s.cfg.ServerURL, "/agent/", s.cfg.Token, s.handleMessage)
	if err != nil {
		return fmt.Errorf("failed to create websocket client: %w", err)
	}
	s.cli = cli

	if err := s.cli.Start(); err != nil {
		return fmt.Errorf("failed to start websocket connection: %w", err)
	}
	s.running = true
	s.done = make(chan struct{})
	go s.loop()

	return nil
}

func (s *Service) Stop() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.running {
		return nil
	}

	close(s.done)
	s.running = false
	if s.cli != nil {
		s.cli.Stop()
	}
	return nil
}

func (s *Service) loop() {
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer func() {
		s.Stop()
		heartbeatTicker.Stop()
	}()

	for {
		select {
		case <-s.done:
			return
		case <-s.ctx.Done():
			slog.InfoContext(s.ctx, "daemon is shutting down")
			return
		case <-s.cli.Done():
			slog.InfoContext(s.ctx, "remote host closed WS connection")
			return
		case <-heartbeatTicker.C:
			// send heartbeat
			msg := messaging.AgentMessage{
				Type: messaging.MessageTypeHeartbeat,
				Data: []byte{},
			}
			if err := s.cli.WriteMessage(websocket.BinaryMessage, msg.Marshal()); err != nil {
				slog.ErrorContext(s.ctx, "failed to send heartbeat", "error", err)
			}
		}
	}
}
