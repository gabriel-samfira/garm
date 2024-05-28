package entity

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/cloudbase/garm/auth"
	"github.com/cloudbase/garm/database/common"
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	garmUtil "github.com/cloudbase/garm/util"
	workerCommon "github.com/cloudbase/garm/workers/common"
)

func NewWorker(ctx context.Context, ghEntity params.GithubEntity, store common.Store) *Worker {
	consumerID := fmt.Sprintf("entity-%s", ghEntity.ID)
	ctx = garmUtil.WithContext(
		ctx,
		slog.Any("entity_name", ghEntity.String()),
		slog.Any("entity_type", ghEntity.EntityType),
		slog.Any("worker", consumerID))
	ctx = auth.GetAdminContext(ctx)

	return &Worker{
		ghEntity:   ghEntity,
		store:      store,
		ctx:        ctx,
		stopped:    true,
		consumerID: consumerID,
	}
}

type Worker struct {
	ctx     context.Context
	quit    chan struct{}
	stopped bool

	ghEntity   params.GithubEntity
	consumerID string

	store    common.Store
	consumer dbCommon.Consumer

	mux sync.Mutex
}

func (e *Worker) ID() string {
	return e.ghEntity.ID
}

func (e *Worker) Start() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	if !e.stopped {
		return nil
	}

	slog.InfoContext(
		e.ctx,
		"starting entity worker")

	// Context was canceled, which probably means that GARM
	// is shutting down. No point in allowing this worker to start.
	if e.ctx.Err() != nil {
		return e.ctx.Err()
	}

	if e.consumer != nil {
		e.consumer.Close()
	}

	consumer, err := watcher.RegisterConsumer(
		e.consumerID,
		// Watch for changes to the entity itself (github credentials, pool balancer type, etc.)
		workerCommon.WithEntityFilter(e.ghEntity),
		// Watch for changes to pools belonging to this entity.
		workerCommon.WithEntityPoolFilter(e.ghEntity),
		// Watch for jobs meant for this entity.
		workerCommon.WithEntityJobFilter(e.ghEntity))
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}
	e.consumer = consumer
	e.quit = make(chan struct{})
	e.stopped = false
	go e.loop()
	// nolint:golangci-lint,godox
	// TODO: load pools, jobs, etc. from the store
	return nil
}

func (e *Worker) Stop() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.stopped {
		return nil
	}
	slog.InfoContext(e.ctx, "stopping entity worker")
	e.stopped = true
	close(e.quit)
	e.consumer.Close()
	return nil
}

func (e *Worker) loop() {
	defer func() {
		e.mux.Lock()
		slog.InfoContext(e.ctx, "entity worker loop exited")
		defer e.mux.Unlock()
		if e.stopped {
			return
		}
		e.stopped = true
		close(e.quit)
		e.consumer.Close()
	}()
	for {
		select {
		case payload := <-e.consumer.Watch():
			slog.InfoContext(e.ctx, "received payload", slog.Any("payload", payload))
		case <-e.ctx.Done():
			return
		case <-e.quit:
			return
		}
	}
}
