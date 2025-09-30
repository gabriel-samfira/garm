package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	garmWs "github.com/cloudbase/garm-provider-common/util/websocket"
	"github.com/cloudbase/garm/cmd/garm-agent/config"
	"github.com/cloudbase/garm/cmd/garm-agent/service/runner"
	"github.com/cloudbase/garm/cmd/garm-agent/state"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/workers/websocket/agent/messaging"
	"github.com/gorilla/websocket"
)

var closed = make(chan struct{})

func init() { close(closed) }

func NewService(ctx context.Context, cfg *config.Agent) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate agent config: %w", err)
	}

	forgeType, err := cfg.ForgeType()
	if err != nil {
		return nil, fmt.Errorf("failed to get forge type for agent: %w", err)
	}

	agentState, err := state.NewStateManager(cfg.StateDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create state manager: %w", err)
	}
	return &Service{
		ctx:        ctx,
		cfg:        cfg,
		done:       closed,
		connecting: make(chan struct{}),
		connected:  closed,
		forgeType:  forgeType,
		sessions:   make(map[string]*ShellSession),
		agentState: agentState,
	}, nil
}

type Service struct {
	ctx         context.Context
	cfg         *config.Agent
	cli         *garmWs.Reader
	agentState  *state.StateManager
	runnerAlive bool
	runnerCmd   runner.Worker

	forgeType params.EndpointType

	mux     sync.Mutex
	cliMux  sync.Mutex
	running bool
	done    chan struct{}

	connecting chan struct{}
	connected  chan struct{}

	sessions map[string]*ShellSession
}

func (s *Service) Done() chan struct{} {
	return s.done
}

func (s *Service) getClient() (*garmWs.Reader, error) {
	s.cliMux.Lock()
	cli := s.cli
	s.cliMux.Unlock()

	if cli == nil {
		return nil, fmt.Errorf("websocket client not connected")
	}
	return cli, nil
}

func (s *Service) writeMessage(msg []byte) error {
	cli, err := s.getClient()
	if err != nil {
		return err
	}

	if err := cli.WriteMessage(websocket.BinaryMessage, msg); err != nil {
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

	switch agentMsg.Type {
	case messaging.MessageTypeCreateShell:
		createShell, err := messaging.Unmarshal[messaging.CreateShellMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshall create shell message: %w", err)
		}
		slog.InfoContext(s.ctx, "handling create shell message", "session_id", createShell.ID())
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
		session, err := NewShellSession(s.ctx, createShell, s.writeMessage, s.cfg)
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
		slog.InfoContext(s.ctx, "handling close shell message", "session_id", closedMsg.ID())
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

	if s.cfg.WorkDir != "" {
		if mode, err := os.Stat(s.cfg.WorkDir); err == nil {
			if mode.IsDir() {
				os.Chdir(s.cfg.WorkDir)
			}
		} else {
			slog.ErrorContext(s.ctx, "failed to access work_dir", "work_dir", s.cfg.WorkDir, "error", err)
			return err
		}
	}

	s.running = true
	s.done = make(chan struct{})
	go s.keepAliveLoop()
	go s.loop()
	go s.keepRunnerAlive()

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

func (s *Service) determineRunnerState() params.RunnerStatus {
	state := params.RunnerOffline
	if s.runnerAlive {
		state = params.RunnerIdle
	}

	st, err := s.agentState.GetState()
	if err != nil {
		slog.ErrorContext(s.ctx, "failed to get state", "error", err)
		return state
	}
	if st.JobStarted {
		if !s.runnerAlive {
			// We're comming back online and for some reason, we didn't record
			// that the job was finished, but we did record that the job was started.
			// If the job was started but the runner is offline, then the job was either
			// finished, or canceled.
			state = params.RunnerTerminated
		} else {
			state = params.RunnerActive
		}
	}

	if st.JobFinished {
		state = params.RunnerTerminated
	}

	return state
}

func (s *Service) sendRunnerStatus() {
	status := runner.RunnerState{
		RunnerStatus: s.determineRunnerState(),
	}
	if err := s.sendRunnerStatusMessage(status); err != nil {
		slog.ErrorContext(s.ctx, "failed to send status", "error", err)
	}
}

func (s *Service) sendRunnerStatusMessage(status runner.RunnerState) error {
	asJs, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := messaging.AgentMessage{
		Type: messaging.MessageTypeRunnerUpdate,
		Data: asJs,
	}

	cli, err := s.getClient()
	if err != nil {
		return err
	}

	if err := cli.WriteMessage(websocket.BinaryMessage, msg.Marshal()); err != nil {
		return fmt.Errorf("failed to send runner status: %w", err)
	}
	return nil
}

func (s *Service) SetRunnerStarted(st bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.runnerAlive = st
	s.sendRunnerStatus()
}

func (s *Service) SetJobStarted() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if err := s.agentState.SetJobStarted(); err != nil {
		slog.ErrorContext(s.ctx, "failed to set job started", "error", err)
	}
	// attempt to send message to GARM anyway
	status := runner.RunnerState{
		RunnerStatus: params.RunnerActive,
	}
	if err := s.sendRunnerStatusMessage(status); err != nil {
		slog.ErrorContext(s.ctx, "failed to send status", "error", err)
	}
}
func (s *Service) SetJobFinished() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if err := s.agentState.SetJobFinished(); err != nil {
		slog.ErrorContext(s.ctx, "failed to set job finished", "error", err)
	}
	// attempt to send message to GARM anyway
	status := runner.RunnerState{
		RunnerStatus: params.RunnerTerminated,
	}
	if err := s.sendRunnerStatusMessage(status); err != nil {
		slog.ErrorContext(s.ctx, "failed to send status", "error", err)
	}
}

func (s *Service) sendHeartbeat() error {
	msg := messaging.AgentMessage{
		Type: messaging.MessageTypeHeartbeat,
		Data: []byte{},
	}

	cli, err := s.getClient()
	if err != nil {
		return err
	}

	if err := cli.WriteMessage(websocket.BinaryMessage, msg.Marshal()); err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}
	return nil
}

func (s *Service) sleepWithCancel(d time.Duration) (shouldQuit bool) {
	sleepTicker := time.NewTicker(d)
	defer sleepTicker.Stop()

	select {
	case <-sleepTicker.C:
		return false
	case <-s.done:
	case <-s.ctx.Done():
	}
	return true
}

func (s *Service) keepRunnerAlive() {
retryCreate:
	state := s.determineRunnerState()
	if state == params.RunnerTerminated {
		// no need for this goroutine.
		return
	}
	runnerCommand, err := runner.NewRunnerCommand(s.ctx, s.cfg.RunnerExecArgs, s.cfg.WorkDir, s.forgeType, s)
	if err != nil {
		slog.ErrorContext(s.ctx, "failed to create runner command", "error", err)
		if s.sleepWithCancel(5 * time.Second) {
			return
		}
		goto retryCreate
	}
	s.mux.Lock()
	s.runnerCmd = runnerCommand
	s.mux.Unlock()
	defer s.runnerCmd.Stop()

	retryCount := 0

retryStart:
	if retryCount > 5 {
		slog.WarnContext(s.ctx, "max retry reached", "max_retries", 5)
		return
	}
	runnerState := s.determineRunnerState()
	if runnerState == params.RunnerTerminated {
		// we only attepmt to start the runner if we need to. A runner that has already run a job,
		// should not be started again, even if the agent is still online.
		return
	}
	if err := runnerCommand.Start(); err != nil {
		slog.ErrorContext(s.ctx, "failed to start runner", "error", err)
		retryCount++
		runnerState := s.determineRunnerState()
		if runnerState == params.RunnerOffline {
			// The runner did not run a job as far as we know, but it's failing to start. Send a failed message to GARM.
			s.sendRunnerStatusMessage(runner.RunnerState{RunnerStatus: params.RunnerFailed})
		}
		if s.sleepWithCancel(5 * time.Second) {
			return
		}
		goto retryStart
	}
	retryCount = 0

	for {
		select {
		case <-s.done:
			return
		case <-s.ctx.Done():
			return
		case <-s.runnerCmd.Wait():
			if s.determineRunnerState() == params.RunnerTerminated {
				return
			}
			if s.sleepWithCancel(5 * time.Second) {
				return
			}
			goto retryStart
		}
	}
}

func (s *Service) keepAliveLoop() {
	var sleepTime time.Duration
retryConnecting:
	if sleepTime > 0 {
		if s.sleepWithCancel(sleepTime) {
			return
		}
	}
	for {
		select {
		case <-s.done:
			return
		case <-s.ctx.Done():
			return
		case <-s.connected:
			slog.InfoContext(s.ctx, "attempting to connect to GARM server", "server", s.cfg.ServerURL)
			sleepTime = 5 * time.Second
			cli, err := garmWs.NewReader(s.ctx, s.cfg.ServerURL, "/agent/", s.cfg.Token, s.handleMessage)
			if err != nil {
				slog.WarnContext(s.ctx, "failed to create websocket client", "error", err)
				goto retryConnecting
			}

			s.cliMux.Lock()
			s.cli = cli
			s.cliMux.Unlock()

			if err := s.cli.Start(); err != nil {
				slog.WarnContext(s.ctx, "failed to start websocket connection", "error", err)
				goto retryConnecting
			}
			slog.InfoContext(s.ctx, "successfully connected to GARM", "server", s.cfg.ServerURL)
			s.connected = make(chan struct{})
			close(s.connecting)
		}
	}

}

func (s *Service) loop() {
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer func() {
		slog.InfoContext(s.ctx, "daemon is shutting down")
		s.Stop()
		heartbeatTicker.Stop()
	}()

connecting:
	select {
	case <-s.done:
		return
	case <-s.ctx.Done():
		return
	case <-s.connecting:
	}
	// send initial heartbeat
	if err := s.sendHeartbeat(); err != nil {
		slog.ErrorContext(s.ctx, "failed to send heartbeat", "error", err)
	}
	s.sendRunnerStatus()

	for {
		select {
		case <-s.done:
			return
		case <-s.ctx.Done():
			return
		case <-s.cli.Done():
			slog.InfoContext(s.ctx, "remote host closed WS connection")
			s.connecting = make(chan struct{})
			close(s.connected)
			goto connecting
		case <-heartbeatTicker.C:
			// send heartbeat
			if err := s.sendHeartbeat(); err != nil {
				slog.ErrorContext(s.ctx, "failed to send heartbeat", "error", err)
			}
		}
	}
}
