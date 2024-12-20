package entity

import (
	"log/slog"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/workers/scalesets/worker"
)

func (s *EntityWorker) handleScaleSetUpdateOrCreate(scaleset params.ScaleSet) {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.scalesets[scaleset.ID]
	if ok {
		// An error might have occured when creating a scaleset worker due to invalid credentials
		// or some other reason. In which case, we won't have a scaleset worker set up. This function
		// also handles cases where scale sets were updated. If the issue that caused an error is
		// fixed, a second call to this function should start it successfully. Otherwise, an update
		// will be handled by the worker.
		slog.InfoContext(s.ctx, "scale set worker already exists", "scaleset_id", scaleset.ID)
		return
	}

	var err error
	worker, err := worker.NewScaleSetWorker(s.ctx, scaleset, s.store, s.providers)
	if err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(s.ctx, "failed to get entity", "scaleset_id", scaleset.ID)
		return
	}

	defer func(err error) {
		// cleanup watchers.
		if err != nil {
			worker.Stop()
		}
	}(err)

	err = worker.Start()
	if err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(s.ctx, "failed to start scale set worker", "scaleset_id", scaleset.ID)
		return
	}
	s.scalesets[scaleset.ID] = worker
}

func (s *EntityWorker) handleScaleSetDelete(scaleset params.ScaleSet) {
	set, ok := s.scalesets[scaleset.ID]
	if !ok {
		slog.InfoContext(s.ctx, "could not find scaleset worker", "scaleset_id", scaleset.ID)
		return
	}

	if err := set.Stop(); err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(s.ctx, "failed to stop scaleset worker", "scaleset_id", scaleset.ID)
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.scalesets, scaleset.ID)
}

func (s *EntityWorker) handleScaleSetEvent(event dbCommon.ChangePayload) {
	if event.EntityType != dbCommon.ScaleSetEntityType {
		slog.InfoContext(s.ctx, "got invalid entity type from watcher", "entity_type", event.EntityType)
		return
	}

	switch event.Operation {
	case dbCommon.CreateOperation, dbCommon.UpdateOperation:
		scaleSet, ok := event.Payload.(params.ScaleSet)
		if !ok {
			slog.InfoContext(s.ctx, "got invalid scaleset payload from watcher", "payload", event.Payload)
		}
		s.handleScaleSetUpdateOrCreate(scaleSet)
	case dbCommon.DeleteOperation:
		scaleSet, ok := event.Payload.(params.ScaleSet)
		if !ok {
			slog.InfoContext(s.ctx, "got invalid scaleset payload from watcher", "payload", event.Payload)
		}
		s.handleScaleSetDelete(scaleSet)
	default:
		slog.InfoContext(s.ctx, "invalid operation for scaleset controller", "operation", event.Operation)
	}
}
