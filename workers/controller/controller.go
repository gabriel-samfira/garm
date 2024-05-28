package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/cloudbase/garm/auth"
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/workers/common"
	"github.com/cloudbase/garm/workers/controller/entity"
)

func NewController(ctx context.Context, store dbCommon.Store) common.Worker {
	ctx = auth.GetAdminContext(ctx)
	return &Controller{
		store:    store,
		ctx:      ctx,
		entities: map[string]*entity.Worker{},
		errChan:  make(chan error, 2),
		stopped:  true,
	}
}

type Controller struct {
	ctx     context.Context
	quit    chan struct{}
	done    chan struct{}
	stopped bool
	mux     sync.Mutex
	errChan chan error

	entities map[string]*entity.Worker

	store    dbCommon.Store
	consumer dbCommon.Consumer
}

func (c *Controller) startWorkerForEntity(ghEntity params.GithubEntity) error {
	_, ok := c.entities[ghEntity.ID]
	if !ok {
		wrk := entity.NewWorker(c.ctx, ghEntity, c.store)
		c.entities[ghEntity.ID] = wrk
		if err := wrk.Start(); err != nil {
			return fmt.Errorf("cannot start worker: %w", err)
		}
	}
	return nil
}

func (c *Controller) loadEntities() error {
	repos, err := c.store.ListRepositories(c.ctx)
	if err != nil {
		return fmt.Errorf("cannot list repositories: %w", err)
	}
	orgs, err := c.store.ListOrganizations(c.ctx)
	if err != nil {
		return fmt.Errorf("cannot list organizations: %w", err)
	}
	enterprises, err := c.store.ListEnterprises(c.ctx)
	if err != nil {
		return fmt.Errorf("cannot list enterprises: %w", err)
	}

	for _, repo := range repos {
		ghEntity, err := repo.GetEntity()
		if err != nil {
			return fmt.Errorf("cannot get github entity: %w", err)
		}
		if err := c.startWorkerForEntity(ghEntity); err != nil {
			return fmt.Errorf("cannot start worker for repository: %w", err)
		}
	}
	for _, org := range orgs {
		ghEntity, err := org.GetEntity()
		if err != nil {
			return fmt.Errorf("cannot get github entity: %w", err)
		}
		if err := c.startWorkerForEntity(ghEntity); err != nil {
			return fmt.Errorf("cannot start worker for organization: %w", err)
		}
	}
	for _, ent := range enterprises {
		ghEntity, err := ent.GetEntity()
		if err != nil {
			return fmt.Errorf("cannot get github entity: %w", err)
		}
		if err := c.startWorkerForEntity(ghEntity); err != nil {
			return fmt.Errorf("cannot start worker for enterprise: %w", err)
		}
	}
	return nil
}

func (c *Controller) Start() (err error) {
	defer func() {
		if err != nil {
			c.Stop()
			select {
			case c.errChan <- err:
			default:
			}
		}
	}()

	c.mux.Lock()
	defer c.mux.Unlock()
	if !c.stopped {
		return nil
	}
	slog.InfoContext(c.ctx, "starting controller")

	// Context was canceled, which probably means that GARM
	// is shutting down. No point in allowing this worker to start.
	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}

	c.quit = make(chan struct{})
	c.done = make(chan struct{})
	c.stopped = false

	slog.InfoContext(c.ctx, "registering consumer", "name", "controller")
	consumer, err := watcher.RegisterConsumer(
		"controller",
		common.WithEntityTypeFilter(dbCommon.RepositoryEntityType),
		common.WithEntityTypeFilter(dbCommon.OrganizationEntityType),
		common.WithEntityTypeFilter(dbCommon.EnterpriseEntityType))
	if err != nil {
		return fmt.Errorf("cannot register consumer: %w", err)
	}
	c.consumer = consumer

	go c.loop()
	slog.InfoContext(c.ctx, "loading entities")
	if err := c.loadEntities(); err != nil {
		return fmt.Errorf("cannot load entities: %w", err)
	}
	return nil
}

func (c *Controller) Stop() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.stopped {
		return nil
	}
	c.stopped = true
	close(c.quit)
	if c.consumer != nil {
		c.consumer.Close()
	}

	for _, worker := range c.entities {
		if err := worker.Stop(); err != nil {
			slog.ErrorContext(c.ctx, "cannot stop worker", slog.Any("error", err))
		}
		delete(c.entities, worker.ID())
	}
	return nil
}

// Wait will block forever until the controller is stopped.
func (c *Controller) Wait() chan error {
	return c.errChan
}

func (c *Controller) handleEntityEvent(ghEntity params.GithubEntity, op dbCommon.OperationType) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.stopped {
		return
	}

	switch op {
	case dbCommon.CreateOperation:
		if err := c.startWorkerForEntity(ghEntity); err != nil {
			// Something really bad happened, and we couldn't create the worker.
			// We terminate the controller and allow the error to bubble up
			// to the main function, which in turn will stop the whole GARM
			// process.
			c.errChan <- err
			c.Stop()
			return
		}
	case dbCommon.DeleteOperation:
		worker, ok := c.entities[ghEntity.ID]
		if !ok {
			return
		}
		slog.InfoContext(c.ctx, "stopping worker", "entity_id", ghEntity.ID)
		if err := worker.Stop(); err != nil {
			slog.ErrorContext(c.ctx, "cannot stop worker", slog.Any("error", err))
		}
		slog.InfoContext(c.ctx, "deleting worker", "entity_id", ghEntity.ID)
		delete(c.entities, ghEntity.ID)
	default:
		// At this level we only care about the entity ID so we can keep track of
		// the workers. The workers will have a watcher on the specific entity
		// and will handle the updates.
		slog.DebugContext(c.ctx, "ignoring operation", slog.Any("operation", op))
	}
}

func (c *Controller) loop() {
	defer func() {
		c.mux.Lock()
		defer c.mux.Unlock()
		if c.stopped {
			return
		}
		c.stopped = true
		c.consumer.Close()
		c.errChan <- nil
		close(c.done)
	}()
	for {
		select {
		case payload := <-c.consumer.Watch():
			entity, err := githubEntityFromPayload(payload)
			if err != nil {
				slog.ErrorContext(c.ctx, "cannot get entity from payload", slog.Any("error", err), slog.Any("payload", payload))
				continue
			}
			go c.handleEntityEvent(entity, payload.Operation)
		case <-c.ctx.Done():
			return
		case <-c.quit:
			return
		}
	}
}
