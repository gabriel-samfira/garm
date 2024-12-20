package entity

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/cloudbase/garm/auth"
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner/common"
	garmUtil "github.com/cloudbase/garm/util"
	"github.com/cloudbase/garm/util/github"
	"github.com/cloudbase/garm/util/github/scalesets"
	"github.com/cloudbase/garm/workers/scalesets/worker"
	"github.com/pkg/errors"
)

func NewEntityController(ctx context.Context, ghEntity params.GithubEntity, store dbCommon.Store, providers map[string]common.Provider, controllerInfo params.ControllerInfo) (*EntityWorker, error) {
	consumerID := fmt.Sprintf("entity-%s", ghEntity.ID)
	ctx = garmUtil.WithSlogContext(
		ctx,
		slog.Any("entity_name", ghEntity.String()),
		slog.Any("entity_type", ghEntity.EntityType),
		slog.Any("worker", consumerID))
	ctx = auth.GetAdminContext(ctx)

	controller := &EntityWorker{
		ctx:        ctx,
		consumerID: consumerID,
		store:      store,
		scalesets:  make(map[uint]*worker.ScaleSetWorker),
		providers:  providers,
		githubResources: &githubResources{
			credentials: ghEntity.Credentials,
		},
		ghEntity:       ghEntity,
		controllerInfo: controllerInfo,
	}
	return controller, nil
}

type EntityWorker struct {
	consumerID     string
	store          dbCommon.Store
	consumer       dbCommon.Consumer
	ghEntity       params.GithubEntity
	controllerInfo params.ControllerInfo

	githubResources *githubResources

	providers map[string]common.Provider
	scalesets map[uint]*worker.ScaleSetWorker

	mux     sync.Mutex
	ctx     context.Context
	quit    chan struct{}
	stopped bool
}

func (e *EntityWorker) Start() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	if !e.stopped {
		return nil
	}

	slog.InfoContext(
		e.ctx,
		"starting entity worker")

	ghCli, err := github.GithubClient(e.ctx, e.ghEntity)
	if err != nil {
		e.githubResources.err = err
		return errors.Wrap(err, "getting github client")
	}
	e.githubResources.ghCli = ghCli

	scalesetCli, err := scalesets.NewClient(ghCli)
	if err != nil {
		e.githubResources.err = err
		return errors.Wrap(err, "getting scalesets client")
	}
	e.githubResources.scaleSetClient = scalesetCli

	consumer, err := watcher.RegisterConsumer(
		e.ctx, e.consumerID,
		composeWatcherFilters(e.ghEntity))

	if err != nil {
		return fmt.Errorf("failed to create consumer for scaleset controller: %w", err)
	}
	e.consumer = consumer

	e.quit = make(chan struct{})
	e.stopped = false
	go e.loop()
	// nolint:golangci-lint,godox
	// TODO: load pools, scale sets jobs, etc, from the store
	return nil

	// scaleSets, err := s.store.ListAllScaleSets(s.ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to start scale set controller: %w", err)
	// }

	// for _, set := range scaleSets {
	// 	work, err := worker.NewScaleSetWorker(s.ctx, set, s.store, s.providers)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create scaleset worker for %d (%s): %w", set.ID, set.Name, err)
	// 	}

	// 	if err := work.Start(); err != nil {
	// 		return fmt.Errorf("failed to start scaleset worker for %d (%s): %w", set.ID, set.Name, err)
	// 	}
	// 	s.scalesets[set.ID] = work
	// }

	// return nil
}

func (e *EntityWorker) Stop() error {
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

func (e *EntityWorker) loop() {
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
			go e.handleWatcherEvent(payload)
		case <-e.ctx.Done():
			return
		case <-e.quit:
			return
		}
	}
}
