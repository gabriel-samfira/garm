// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	garmWs "github.com/cloudbase/garm-provider-common/util/websocket"
	"github.com/cloudbase/garm/client/agent"
	"github.com/cloudbase/garm/workers/websocket/agent/messaging"
)

// agentTokenCmd represents the agent token command
var agentTokenCmd = &cobra.Command{
	Use:          "agent",
	SilenceUsage: false,
	Short:        "Handle agent operations",
	Long:         `This command exposes a number of agent operations.`,
	Run:          nil,
}

var agentTokenCreateCmd = &cobra.Command{
	Use:          "token-create",
	Short:        "Create an agent token",
	Long:         `Create a metrics token.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		if len(args) != 1 {
			return fmt.Errorf("requires a runner name")
		}

		getAgentTokenReq := agent.NewGetAgentJWTTokenParams()
		getAgentTokenReq.AgentName = args[0]
		response, err := apiCli.Agent.GetAgentJWTToken(getAgentTokenReq, authToken)
		if err != nil {
			return err
		}
		fmt.Println(response.Payload.Token)

		return nil
	},
}

var agentShellCmd = &cobra.Command{
	Use:          "shell",
	Short:        "Execute an interactive shell",
	Long:         `Execute an interactive shell on the runner.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		if needsInit {
			return errNeedsInitError
		}

		if len(args) != 1 {
			return fmt.Errorf("requires a runner name")
		}

		var sessionID uuid.UUID

		handlerErr := make(chan struct{})
		resizeCh := make(chan [2]int, 1)
		handler := func(msgType int, msg []byte) error {
			switch msgType {
			case websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseMessage:
				os.Stderr.Write([]byte("remote server closed the connection"))
				close(handlerErr)
			case websocket.BinaryMessage, websocket.TextMessage:
				agentMsg, err := messaging.UnmarshalAgentMessage(msg)
				if err != nil {
					os.Stderr.Write([]byte("failed to unmarshal message"))
					close(handlerErr)
				}
				switch agentMsg.Type {
				case messaging.MessageTypeShellReady:
					shellReady, err := messaging.Unmarshal[messaging.ShellReadyMessage](agentMsg)
					if err != nil {
						os.Stderr.Write(fmt.Appendf(nil, "failed to unmarshal shell ready: %q", err))
						close(handlerErr)
					}
					sessionID = shellReady.SessionID
					if w, h, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
						resizeCh <- [2]int{w, h}
					}
				case messaging.MessageTypeShellExit:
					close(handlerErr)
				case messaging.MessageTypeShellData:
					shellData, err := messaging.Unmarshal[messaging.ShellDataMessage](agentMsg)
					if err != nil {
						os.Stderr.Write([]byte("failed to unmarshal shell data message"))
						close(handlerErr)
					}
					os.Stdout.Write(shellData.Data)
				}
			default:
				os.Stdout.Write(fmt.Appendf(nil, "invalid message type: %v", msgType))
			}
			return nil
		}

		// Put terminal in raw mode
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
		// Channel to stop on Ctrl+C
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)

		ctx, stop := signal.NotifyContext(context.Background(), signals...)
		defer stop()

		reader, err := garmWs.NewReader(ctx, mgr.BaseURL, fmt.Sprintf("/api/v1/ws/agent/%s/shell", args[0]), mgr.Token, handler)
		if err != nil {
			return err
		}

		if err := reader.Start(); err != nil {
			return err
		}

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := os.Stdin.Read(buf)
				if err != nil {
					os.Stderr.Write(fmt.Appendf(nil, "failed to write message: %q", err))
					close(handlerErr)
					return
				}

				if n > 0 && sessionID != uuid.Nil {
					msg := messaging.ShellDataMessage{
						SessionID: sessionID,
						Data:      buf[:n],
					}
					if err := reader.WriteMessage(websocket.BinaryMessage, msg.Marshal()); err != nil {
						os.Stderr.Write(fmt.Appendf(nil, "failed to write message: %q", err))
						close(handlerErr)
						return
					}
				}
			}
		}()

		// ---- Watch terminal resize ----
		go watchTermResize(resizeCh, sessionID)

		// ---- Send resize messages ----
		go func() {
			for size := range resizeCh {
				if sessionID == uuid.Nil {
					continue
				}
				msg := messaging.ShellResizeMessage{
					SessionID: sessionID,
					Cols:      uint16(size[0]),
					Rows:      uint16(size[1]),
				}
				reader.WriteMessage(websocket.BinaryMessage, msg.Marshal())
			}
		}()

		select {
		case <-ctx.Done():
		case <-reader.Done():
		case <-handlerErr:
		}

		return nil
	},
}

func init() {
	agentTokenCmd.AddCommand(
		agentTokenCreateCmd,
		// agentShellCmd,
	)

	rootCmd.AddCommand(agentTokenCmd)
}
