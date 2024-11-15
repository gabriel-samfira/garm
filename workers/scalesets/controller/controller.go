package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/workers/scalesets/worker"
)

func NewScaleSetController(ctx context.Context, store dbCommon.Store) (*ScaleSetController, error) {
	consumerID := "scaleset-controller"
	consumer, err := watcher.RegisterConsumer(
		ctx, consumerID,
		watcher.WithAll(
			watcher.WithEntityTypeFilter(dbCommon.ScaleSetEntityType),
			watcher.WithAny(
				watcher.WithOperationTypeFilter(dbCommon.CreateOperation),
				watcher.WithOperationTypeFilter(dbCommon.DeleteOperation),
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer for scaleset controller: %w", err)
	}
	controller := &ScaleSetController{
		ctx:       ctx,
		consumer:  consumer,
		scalesets: make(map[uint]*worker.ScaleSetWorker),
	}
	return controller, nil
}

type ScaleSetController struct {
	store    dbCommon.Store
	consumer dbCommon.Consumer

	scalesets map[uint]*worker.ScaleSetWorker

	mux     sync.Mutex
	ctx     context.Context
	quit    chan struct{}
	running bool
}

func (s *ScaleSetController) Start() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	scaleSets, err := s.store.ListAllScaleSets(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to start scale set controller: %w", err)
	}

	for _, set := range scaleSets {
		work, err := worker.NewScaleSetWorker(s.ctx, set, s.store)
		if err != nil {
			return fmt.Errorf("failed to create scaleset worker for %d (%s): %w", set.ID, set.Name, err)
		}

		if err := work.Start(); err != nil {
			return fmt.Errorf("failed to start scaleset worker for %d (%s): %w", set.ID, set.Name, err)
		}
		s.scalesets[set.ID] = work
	}

	s.running = true
	go s.runWatcher()

	return nil
}

func (s *ScaleSetController) Stop() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.running {
		return nil
	}

	close(s.quit)
	s.running = false

	return nil
}

func (s *ScaleSetController) handleScaleSetDelete(scaleset params.ScaleSet) {
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

func (s *ScaleSetController) handleScaleSetUpdateOrCreate(scaleset params.ScaleSet) {
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

	worker, err := worker.NewScaleSetWorker(s.ctx, scaleset, s.store)
	if err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(s.ctx, "failed to get entity", "scaleset_id", scaleset.ID)
		return
	}

	if err := worker.Start(); err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(s.ctx, "failed to start scale set worker", "scaleset_id", scaleset.ID)
		return
	}
	s.scalesets[scaleset.ID] = worker
}

func (s *ScaleSetController) handleWatcherEvent(event dbCommon.ChangePayload) {
	if event.EntityType != dbCommon.ScaleSetEntityType {
		slog.InfoContext(s.ctx, "got invalid entity type from watcher", "entity_type", event.EntityType)
		return
	}

	switch event.Operation {
	case dbCommon.CreateOperation:
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
	case dbCommon.UpdateOperation:
		scaleSet, ok := event.Payload.(params.ScaleSet)
		if !ok {
			slog.InfoContext(s.ctx, "got invalid scaleset payload from watcher", "payload", event.Payload)
		}
		s.handleScaleSetUpdateOrCreate(scaleSet)
	default:
		slog.InfoContext(s.ctx, "invalid operation for scaleset controller", "operation", event.Operation)
	}
}

func (s *ScaleSetController) runWatcher() {
	defer s.consumer.Close()
	for {
		select {
		case <-s.quit:
			return
		case <-s.ctx.Done():
			return
		case event, ok := <-s.consumer.Watch():
			if !ok {
				return
			}
			go s.handleWatcherEvent(event)
		}
	}
}
