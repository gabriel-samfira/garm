package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/pkg/errors"

	runnerErrors "github.com/cloudbase/garm-provider-common/errors"
	commonParams "github.com/cloudbase/garm-provider-common/params"
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner/common"
	"github.com/cloudbase/garm/util/github"
	"github.com/cloudbase/garm/util/github/scalesets"
)

func NewScaleSetWorker(ctx context.Context, scaleset params.ScaleSet, store dbCommon.Store, providers map[string]common.Provider) (worker *ScaleSetWorker, err error) {
	if scaleset.ID == 0 {
		return nil, fmt.Errorf("invalid scale set")
	}

	var scalesetConsumer dbCommon.Consumer
	var scalesetInstanceConsumer dbCommon.Consumer
	defer func() {
		if err != nil {
			if scalesetConsumer != nil {
				scalesetConsumer.Close()
			}

			if scalesetInstanceConsumer != nil {
				scalesetInstanceConsumer.Close()
			}
		}
	}()

	entity, err := scaleset.GithubEntity()
	if err != nil {
		return nil, fmt.Errorf("failed to get entity from scale set: %w", err)
	}
	ghEntity, err := store.GetGithubEntity(ctx, entity.EntityType, entity.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get github entity from db: %w", err)
	}

	scaleSetConsumerID := fmt.Sprintf("scaleset-worker-%d", scaleset.ID)
	scalesetConsumer, err = watcher.RegisterConsumer(
		ctx, scaleSetConsumerID,
		composeWatcherFilters(scaleset, ghEntity),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	scaleSetInstanceConsumerID := fmt.Sprintf("scaleset-instance-worker-%d", scaleset.ID)
	scalesetInstanceConsumer, err = watcher.RegisterConsumer(
		ctx, scaleSetInstanceConsumerID,
		watcher.WithScaleSetInstanceFilter(scaleset),
	)

	ghClient, err := github.GithubClient(ctx, ghEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to get github client: %w", err)
	}

	scalesetClient, err := scalesets.NewClient(ghClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get scaleset client: %w", err)
	}

	controllerInfo, err := store.ControllerInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get controller info: %w", err)
	}

	return &ScaleSetWorker{
		controllerInfo:           controllerInfo,
		scalesetClient:           scalesetClient,
		ghCli:                    ghClient,
		store:                    store,
		entity:                   ghEntity,
		scaleset:                 scaleset,
		scalesetConsumer:         scalesetConsumer,
		scalesetInstanceConsumer: scalesetInstanceConsumer,
		ctx:                      ctx,
		wg:                       &sync.WaitGroup{},
		providers:                providers,
		quit:                     make(chan struct{}),
	}, nil
}

type ScaleSetWorker struct {
	controllerInfo params.ControllerInfo
	scalesetClient *scalesets.ScaleSetClient
	ghCli          common.GithubClient
	providers      map[string]common.Provider

	store                    dbCommon.Store
	scaleset                 params.ScaleSet
	scalesetConsumer         dbCommon.Consumer
	scalesetInstanceConsumer dbCommon.Consumer
	ctx                      context.Context
	quit                     chan struct{}
	running                  bool

	faultReason error
	tools       []commonParams.RunnerApplicationDownload

	entity params.GithubEntity

	mux sync.Mutex
	wg  *sync.WaitGroup
}

func (s *ScaleSetWorker) ID() uint {
	return s.scaleset.ID
}

// PoolID returns a unique ID for the scale set which will be used
// to tag instances in providers. This will allow us to reuse the current
// provider interface without making it aware of scale sets.
func (s *ScaleSetWorker) PoolID() string {
	return fmt.Sprintf("%s-%d", s.ghCli.GetEntity().ID, s.ID())
}

func (s *ScaleSetWorker) Start() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.running {
		return nil
	}

	initialToolUpdate := make(chan struct{}, 1)
	go func() {
		slog.Info("running initial tool update")
		if err := s.updateTools(); err != nil {
			slog.With(slog.Any("error", err)).Error("failed to update tools")
		}
		initialToolUpdate <- struct{}{}
	}()

	go s.runScaleSetWatcher()
	go s.runScaleSetInstanceWatcher()
	go func() {
		select {
		case <-s.quit:
			return
		case <-s.ctx.Done():
			return
		case <-initialToolUpdate:
		}
		defer close(initialToolUpdate)
		go s.startLoopForFunction(s.updateTools, 50*time.Minute, "update_tools", true)
	}()
	s.running = true
	return nil
}

func (s *ScaleSetWorker) Stop() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.scalesetConsumer != nil {
		s.scalesetConsumer.Close()
	}

	if s.scalesetInstanceConsumer != nil {
		s.scalesetInstanceConsumer.Close()
	}

	if !s.running {
		return nil
	}

	close(s.quit)
	s.running = false

	if err := s.Wait(); err != nil {
		return errors.Wrap(err, "stopping scale set")
	}
	return nil
}

func (s *ScaleSetWorker) Wait() error {
	done := make(chan struct{})
	timer := time.NewTimer(60 * time.Second)
	go func() {
		s.wg.Wait()
		timer.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-timer.C:
		return errors.Wrap(runnerErrors.ErrTimeout, "waiting for scaleset to stop")
	}
	return nil
}
