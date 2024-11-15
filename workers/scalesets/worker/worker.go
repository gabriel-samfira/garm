package worker

import (
	"context"
	"fmt"
	"sync"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner/common"
	"github.com/cloudbase/garm/util/github"
	"github.com/cloudbase/garm/util/github/scalesets"
)

func NewScaleSetWorker(ctx context.Context, scaleset params.ScaleSet, store dbCommon.Store) (worker *ScaleSetWorker, err error) {
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
		quit:                     make(chan struct{}),
	}, nil
}

type ScaleSetWorker struct {
	controllerInfo params.ControllerInfo
	scalesetClient *scalesets.ScaleSetClient
	ghCli          common.GithubClient

	store                    dbCommon.Store
	scaleset                 params.ScaleSet
	scalesetConsumer         dbCommon.Consumer
	scalesetInstanceConsumer dbCommon.Consumer
	ctx                      context.Context
	quit                     chan struct{}
	running                  bool

	active      bool
	faultReason error

	entity params.GithubEntity

	mux sync.Mutex
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

	go s.runScaleSetWatcher()
	go s.runScaleSetInstanceWatcher()
	go s.consolidate()
	s.running = true
	return nil
}

func (s *ScaleSetWorker) Stop() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.running {
		return nil
	}

	close(s.quit)
	s.running = false

	return nil
}
