package worker

import (
	"log/slog"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/params"
)

// handleWatcherEvent dispatches the event to the appropriate handler. The handlers must be
// idempotent to avoid event feedback loops that will endlessly send events which cause changes
// in the db which in turn generate new events and so on.
func (s *ScaleSetWorker) handleWatcherEvent(event dbCommon.ChangePayload) {
	dbEntityType := dbCommon.DatabaseEntityType(s.entity.EntityType)
	switch event.EntityType {
	case dbCommon.GithubCredentialsEntityType:
		s.handleGithubCredentialsWatcherEvent(event)
	case dbEntityType:
		s.handleGithubEntityWatcherEvent(event)
	case dbCommon.ScaleSetEntityType:
		s.handleScaleSetWatcherEvent(event)
	case dbCommon.ControllerEntityType:
		s.handleControllerWatcherEvent(event)
	default:
		slog.ErrorContext(s.ctx, "invalid entity type", "entity_type", event.EntityType)
	}
}

func (s *ScaleSetWorker) handleScaleSetWatcherEvent(event dbCommon.ChangePayload) {
	if event.EntityType != dbCommon.ScaleSetEntityType {
		slog.InfoContext(s.ctx, "invalid event received by scaleset event handler", "event_type", event.EntityType)
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()
	scaleSet, ok := event.Payload.(params.ScaleSet)
	if !ok {
		slog.ErrorContext(s.ctx, "invalid payload type for event", "expected", dbCommon.ScaleSetEntityType)
		return
	}
	s.scaleset = scaleSet
}

func (s *ScaleSetWorker) handleGithubEntityWatcherEvent(event dbCommon.ChangePayload)      {}
func (s *ScaleSetWorker) handleGithubCredentialsWatcherEvent(event dbCommon.ChangePayload) {}
func (s *ScaleSetWorker) handleControllerWatcherEvent(event dbCommon.ChangePayload)        {}

func (s *ScaleSetWorker) runScaleSetWatcher() {
	defer s.scalesetConsumer.Close()
	for {
		select {
		case <-s.quit:
			return
		case <-s.ctx.Done():
			return
		case event, ok := <-s.scalesetConsumer.Watch():
			if !ok {
				return
			}
			go s.handleWatcherEvent(event)
		}
	}
}
