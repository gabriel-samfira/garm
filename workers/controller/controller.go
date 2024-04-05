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
	"github.com/cloudbase/garm/runner"
	"github.com/cloudbase/garm/workers/common"
	"github.com/cloudbase/garm/workers/controller/entity"
)

func NewController(ctx context.Context, store dbCommon.Store, r *runner.Runner) common.Worker {
	ctx = auth.GetAdminContext(ctx)
	return &Controller{
		store:    store,
		ctx:      ctx,
		entities: map[string]*entity.Worker{},
		errChan:  make(chan error, 2),
		stopped:  true,
		runner:   r,
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
	runner   *runner.Runner
	consumer dbCommon.Consumer
}

func (c *Controller) startWorkerForEntity(ent entityGetter) error {
	ghEntity, err := ent.GetEntity()
	if err != nil {
		return fmt.Errorf("cannot get entity: %w", err)
	}
	_, ok := c.entities[ghEntity.ID]
	if !ok {
		wrk := entity.NewWorker(c.ctx, ghEntity, c.runner)
		c.entities[ghEntity.ID] = wrk
		if err := wrk.Start(); err != nil {
			return fmt.Errorf("cannot start worker: %w", err)
		}
	}
	return nil
}

func (c *Controller) loadEntities() error {
	repos, err := c.runner.ListRepositories(c.ctx)
	if err != nil {
		return fmt.Errorf("cannot list repositories: %w", err)
	}
	orgs, err := c.runner.ListOrganizations(c.ctx)
	if err != nil {
		return fmt.Errorf("cannot list organizations: %w", err)
	}
	enterprises, err := c.runner.ListEnterprises(c.ctx)
	if err != nil {
		return fmt.Errorf("cannot list enterprises: %w", err)
	}

	for _, repo := range repos {
		if err := c.startWorkerForEntity(repo); err != nil {
			return fmt.Errorf("cannot start worker for repository: %w", err)
		}
	}
	for _, org := range orgs {
		if err := c.startWorkerForEntity(org); err != nil {
			return fmt.Errorf("cannot start worker for organization: %w", err)
		}
	}
	for _, ent := range enterprises {
		if err := c.startWorkerForEntity(ent); err != nil {
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
		worker, ok := c.entities[ghEntity.ID]
		if !ok {
			worker = entity.NewWorker(c.ctx, ghEntity, c.runner)
			c.entities[ghEntity.ID] = worker
		}
		if err := worker.Start(); err != nil {
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
