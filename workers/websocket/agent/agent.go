package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	runnerErrors "github.com/cloudbase/garm-provider-common/errors"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner"
	garmUtil "github.com/cloudbase/garm/util"
	"github.com/cloudbase/garm/workers/websocket/agent/messaging"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 16384 // 16 KB
)

func NewAgent(ctx context.Context, conn *websocket.Conn, instance params.Instance, store runner.AgentStoreOps) (*Agent, error) {
	if conn == nil {
		return nil, fmt.Errorf("missing connection for agent")
	}
	ctx = garmUtil.WithSlogContext(
		ctx,
		slog.Any("worker", "agent"),
		slog.Any("agent_name", instance.Name),
	)
	// waiting on a nil channel will block forever. Create a channel here and close it,
	// ensuring that even if we forget to call Start() before we call Done(), we never deadlock
	// when waiting on Done().
	deadChan := make(chan struct{})
	close(deadChan)

	return &Agent{
		ctx:           ctx,
		conn:          conn,
		instance:      instance,
		agentStore:    store,
		done:          deadChan,
		shellSessions: make(map[string]*ClientSession),
	}, nil
}

type Agent struct {
	ctx        context.Context
	instance   params.Instance
	mux        sync.Mutex
	writeMux   sync.Mutex
	conn       *websocket.Conn
	agentStore runner.AgentStoreOps

	running bool
	done    chan struct{}

	shellSessions map[string]*ClientSession
}

func (a *Agent) CreateShellSession(ctx context.Context, sessionID uuid.UUID, clientConn *websocket.Conn) (*ClientSession, error) {
	a.mux.Lock()
	defer a.mux.Unlock()

	_, ok := a.shellSessions[sessionID.String()]
	if ok {
		return nil, runnerErrors.NewConflictError("session ID %q already in use", sessionID)
	}
	sess, err := NewClientSession(ctx, clientConn, a.writeMessage, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new client session: %w", err)
	}

	if err := sess.Start(); err != nil {
		return nil, fmt.Errorf("failed to start client session: %w", err)
	}
	a.shellSessions[sessionID.String()] = sess
	return sess, nil
}

func (a *Agent) RemoveClientSession(sessionID uuid.UUID, safe bool) error {
	if !safe {
		a.mux.Lock()
		defer a.mux.Unlock()
	}
	sess, ok := a.shellSessions[sessionID.String()]
	if !ok {
		return nil
	}

	if err := sess.Stop(); err != nil {
		return fmt.Errorf("failed to stop session")
	}

	delete(a.shellSessions, sessionID.String())
	return nil
}

func (a *Agent) Done() <-chan struct{} {
	return a.done
}

func (a *Agent) IsRunning() bool {
	return a.running
}

func (a *Agent) Start() error {
	a.mux.Lock()
	defer a.mux.Unlock()

	if a.running {
		return nil
	}

	a.done = make(chan struct{})
	a.running = true
	go a.agentReader()
	go a.loop()
	return nil
}

func (a *Agent) Stop() error {
	a.mux.Lock()
	defer a.mux.Unlock()

	if !a.running {
		return nil
	}
	slog.InfoContext(a.ctx, "removing sessions")
	for _, val := range a.shellSessions {
		slog.InfoContext(a.ctx, "removing session", "session_id", val.sessionID)
		a.RemoveClientSession(val.sessionID, true)
	}

	a.running = false
	slog.InfoContext(a.ctx, "sending websocket close message")
	a.writeMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	slog.InfoContext(a.ctx, "closing connection")
	a.conn.Close()
	slog.InfoContext(a.ctx, "closing done channel")
	close(a.done)
	return nil
}

func (a *Agent) writeMessage(messageType int, message []byte) error {
	a.writeMux.Lock()
	defer a.writeMux.Unlock()
	if err := a.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}
	if err := a.conn.WriteMessage(messageType, message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

// agentReader listens for messages sent by the garm-agent. It unmarshals the message and
// routes it to appropriate functions.
func (a *Agent) agentReader() {
	defer func() {
		slog.InfoContext(a.ctx, ">>> stopping agent")
		a.Stop()
	}()
	a.conn.SetReadLimit(maxMessageSize)
	a.conn.SetPongHandler(func(string) error {
		if err := a.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return err
		}
		return nil
	})
	for {
		if err := a.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			slog.With(slog.Any("error", err)).Error("failed to set read deadline")
		}
		mt, data, err := a.conn.ReadMessage()
		if err != nil {
			if IsErrorOfInterest(err) {
				slog.ErrorContext(a.ctx, "error reading websocket message", slog.Any("error", err))
			}
			return
		}

		slog.InfoContext(a.ctx, "got websocket message", "message_type", mt, "message_data", data)
		if mt == websocket.CloseMessage {
			return
		}

		if err := a.messageHandler(data); err != nil {
			slog.ErrorContext(a.ctx, "error handling message", slog.Any("error", err))
		}
	}
}

func (a *Agent) messageHandler(msg []byte) (err error) {
	if len(msg) < 1 {
		return fmt.Errorf("mesage is too short")
	}
	agentMsg, err := messaging.UnmarshalAgentMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal agetne message")
	}

	switch agentMsg.Type {
	case messaging.MessageTypeHeartbeat:
		slog.DebugContext(a.ctx, "received heartbeat message from agent")
		err = a.agentStore.RecordAgentHeartbeat(a.ctx)
	case messaging.MessageTypeStatusMessage:
		// record status message
	case messaging.MessageTypeShellReady:
		shellReady, err := messaging.Unmarshal[messaging.ShellReadyMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal shell ready message: %w", err)
		}
		session, ok := a.shellSessions[shellReady.ID()]
		if !ok {
			return nil
		}
		if err := session.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	case messaging.MessageTypeShellExit:
		shellExit, err := messaging.Unmarshal[messaging.ShellDataMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal shell exit message: %w", err)
		}
		session, ok := a.shellSessions[shellExit.ID()]
		if !ok {
			return nil
		}
		if err := a.RemoveClientSession(session.sessionID, false); err != nil {
			return fmt.Errorf("failed to remove session: %w", err)
		}
	case messaging.MessageTypeShellData:
		shellData, err := messaging.Unmarshal[messaging.ShellDataMessage](agentMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal shell data message: %w", err)
		}
		session, ok := a.shellSessions[shellData.ID()]
		if !ok {
			return nil
		}
		if err := session.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	}
	return
}

func (a *Agent) loop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		a.Stop()
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			if err := a.writeMessage(websocket.PingMessage, nil); err != nil {
				if IsErrorOfInterest(err) {
					slog.With(slog.Any("error", err)).Error("failed to write ping message")
				}
				return
			}
		case <-a.ctx.Done():
			return
		case <-a.done:
			return
		}
	}
}

func IsErrorOfInterest(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, websocket.ErrCloseSent) {
		return false
	}

	if errors.Is(err, websocket.ErrBadHandshake) {
		return false
	}

	if errors.Is(err, net.ErrClosed) {
		return false
	}

	asCloseErr, ok := err.(*websocket.CloseError)
	if ok {
		switch asCloseErr.Code {
		case websocket.CloseNormalClosure, websocket.CloseGoingAway,
			websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure:
			return false
		}
	}

	return true
}
