package entity

import (
	"log/slog"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/params"
)

func (s *EntityWorker) handleWatcherEvent(event dbCommon.ChangePayload) {
	entityType := dbCommon.DatabaseEntityType(s.ghEntity.EntityType)
	switch event.EntityType {
	case dbCommon.ScaleSetEntityType:
		s.handleScaleSetEvent(event)
	case dbCommon.PoolEntityType:
		slog.InfoContext(s.ctx, "got pool payload event")
	case entityType:
		slog.InfoContext(s.ctx, "got entity payload event")
	case dbCommon.JobEntityType:
		slog.InfoContext(s.ctx, "got job payload event")
	case dbCommon.GithubCredentialsEntityType:
		slog.InfoContext(s.ctx, "got credentials update event")
	case dbCommon.ControllerEntityType:
		slog.InfoContext(s.ctx, "got controller update event")
		s.handleControllerUpdateEvent(event)
	default:
		slog.ErrorContext(s.ctx, "invalid entity type", "entity_type", event.EntityType)
	}
}

func (s *EntityWorker) handleControllerUpdateEvent(event dbCommon.ChangePayload) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if event.EntityType != dbCommon.ControllerEntityType {
		slog.InfoContext(s.ctx, "invalid payload for controller update event", "entity_type", event.EntityType)
		return
	}

	controllerInfo, ok := event.Payload.(params.ControllerInfo)
	if !ok {
		slog.InfoContext(s.ctx, "invalid payload for controller update event")
		return
	}

	slog.DebugContext(s.ctx, "updating controller info", "controller_info", controllerInfo)
	s.controllerInfo = controllerInfo
}
