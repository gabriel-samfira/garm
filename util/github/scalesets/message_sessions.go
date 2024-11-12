// Copyright 2024 Cloudbase Solutions SRL
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

package scalesets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	runnerErrors "github.com/cloudbase/garm-provider-common/errors"
	"github.com/cloudbase/garm/params"
)

const maxCapacityHeader = "X-ScaleSetMaxCapacity"

func NewMessageSession(ctx context.Context, cli *ScaleSetClient, session *params.RunnerScaleSetSession) (*MessageSession, error) {
	sess := &MessageSession{
		ssCli:   cli,
		session: session,
		ctx:     ctx,
		done:    make(chan struct{}),
		closed:  false,
	}
	go sess.loop()
	return sess, nil
}

type MessageSession struct {
	ssCli   *ScaleSetClient
	session *params.RunnerScaleSetSession
	ctx     context.Context

	done    chan struct{}
	closed  bool
	lastErr error

	mux sync.Mutex
}

func (m *MessageSession) Close() error {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.closed {
		return nil
	}
	close(m.done)
	m.closed = true
	return nil
}

func (m *MessageSession) LastError() error {
	return m.lastErr
}

func (m *MessageSession) loop() {
	timer := time.NewTimer(1 * time.Minute)
	defer timer.Stop()
	if m.closed {
		return
	}
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.done:
			return
		case <-timer.C:
			if err := m.maybeRefreshToken(m.ctx); err != nil {
				// We endlessly retry. If it's a transient error, it should eventually
				// work, if it's credentials issues, users can update them.
				slog.With(slog.Any("error", err)).ErrorContext(m.ctx, "failed to refresh message queue token")
				m.lastErr = err
			}
		}
	}
}

func (m *MessageSession) SessionsRelativeURL() (string, error) {
	if m.session == nil {
		return "", fmt.Errorf("session is nil")
	}
	if m.session.RunnerScaleSet == nil {
		return "", fmt.Errorf("runner scale set is nil")
	}
	relativePath := fmt.Sprintf("%s/%d/sessions/%s", scaleSetEndpoint, m.session.RunnerScaleSet.Id, m.session.SessionId.String())
	return relativePath, nil
}

func (m *MessageSession) Refresh(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	relPath, err := m.SessionsRelativeURL()
	if err != nil {
		return fmt.Errorf("failed to get session URL: %w", err)
	}
	req, err := m.ssCli.newActionsRequest(ctx, http.MethodPatch, relPath, nil)
	if err != nil {
		return fmt.Errorf("failed to create message delete request: %w", err)
	}
	resp, err := m.ssCli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete message session: %w", err)
	}

	var refreshedSession params.RunnerScaleSetSession
	if err := json.NewDecoder(resp.Body).Decode(&refreshedSession); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	m.session = &refreshedSession
	return nil
}

func (m *MessageSession) maybeRefreshToken(ctx context.Context) error {
	if m.session == nil {
		return fmt.Errorf("session is nil")
	}

	if m.session.ExpiresIn(2 * time.Minute) {
		if err := m.Refresh(ctx); err != nil {
			return fmt.Errorf("failed to refresh message queue token: %w", err)
		}
	}
	return nil
}

func (m *MessageSession) GetMessage(ctx context.Context, lastMessageId int64, maxCapacity uint) (params.RunnerScaleSetMessage, error) {
	u, err := url.Parse(m.session.MessageQueueUrl)
	if err != nil {
		return params.RunnerScaleSetMessage{}, err
	}

	if lastMessageId > 0 {
		q := u.Query()
		q.Set("lastMessageId", strconv.FormatInt(lastMessageId, 10))
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return params.RunnerScaleSetMessage{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json; api-version=6.0-preview")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.session.MessageQueueAccessToken))
	req.Header.Set(maxCapacityHeader, fmt.Sprintf("%d", maxCapacity))

	resp, err := m.ssCli.Do(req)
	if err != nil {
		return params.RunnerScaleSetMessage{}, fmt.Errorf("request to %s failed: %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		return params.RunnerScaleSetMessage{}, nil
	}

	var message params.RunnerScaleSetMessage
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return params.RunnerScaleSetMessage{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return message, nil
}

func (m *MessageSession) DeleteMessage(ctx context.Context, messageId int64) error {
	u, err := url.Parse(m.session.MessageQueueUrl)
	if err != nil {
		return err
	}

	u.Path = fmt.Sprintf("%s/%d", u.Path, messageId)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.session.MessageQueueAccessToken))

	resp, err := m.ssCli.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (s *ScaleSetClient) CreateMessageSession(ctx context.Context, runnerScaleSetId int, owner string) (*MessageSession, error) {
	path := fmt.Sprintf("%s/%d/sessions", scaleSetEndpoint, runnerScaleSetId)

	newSession := params.RunnerScaleSetSession{
		OwnerName: owner,
	}

	requestData, err := json.Marshal(newSession)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session data: %w", err)
	}

	req, err := s.newActionsRequest(ctx, http.MethodPost, path, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	var createdSession params.RunnerScaleSetSession
	if err := json.NewDecoder(resp.Body).Decode(&createdSession); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &MessageSession{
		ssCli:   s,
		session: &createdSession,
	}, nil
}

func (s *ScaleSetClient) DeleteMessageSession(ctx context.Context, session *MessageSession) error {
	path, err := session.SessionsRelativeURL()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	req, err := s.newActionsRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to create message delete request: %w", err)
	}

	_, err = s.Do(req)
	if err != nil {
		if !errors.Is(err, runnerErrors.ErrNotFound) {
			return fmt.Errorf("failed to delete message session: %w", err)
		}
	}
	return nil
}
