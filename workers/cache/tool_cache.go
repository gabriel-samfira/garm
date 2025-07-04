package cache

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	commonParams "github.com/cloudbase/garm-provider-common/params"
	"github.com/cloudbase/garm/cache"
	"github.com/cloudbase/garm/params"
	garmUtil "github.com/cloudbase/garm/util"
	"github.com/cloudbase/garm/util/github"
)

func newToolsUpdater(ctx context.Context, entity params.GithubEntity) *toolsUpdater {
	return &toolsUpdater{
		ctx:    ctx,
		entity: entity,
		quit:   make(chan struct{}),
	}
}

type toolsUpdater struct {
	ctx context.Context

	entity     params.GithubEntity
	tools      []commonParams.RunnerApplicationDownload
	lastUpdate time.Time

	mux     sync.Mutex
	running bool
	quit    chan struct{}

	reset chan struct{}
}

func (t *toolsUpdater) Start() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.running {
		return nil
	}

	t.running = true
	t.quit = make(chan struct{})

	go t.loop()
	return nil
}

func (t *toolsUpdater) Stop() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.running {
		return nil
	}

	t.running = false
	close(t.quit)

	return nil
}

func (t *toolsUpdater) updateTools() error {
	slog.DebugContext(t.ctx, "updating tools", "entity", t.entity.String())
	entity, ok := cache.GetEntity(t.entity.ID)
	if !ok {
		return fmt.Errorf("getting entity from cache: %s", t.entity.ID)
	}
	ghCli, err := github.Client(t.ctx, entity)
	if err != nil {
		return fmt.Errorf("getting github client: %w", err)
	}

	tools, err := garmUtil.FetchTools(t.ctx, ghCli)
	if err != nil {
		return fmt.Errorf("fetching tools: %w", err)
	}
	t.lastUpdate = time.Now().UTC()
	t.tools = tools

	slog.DebugContext(t.ctx, "updating tools cache", "entity", t.entity.String())
	cache.SetGithubToolsCache(entity, tools)
	return nil
}

func (t *toolsUpdater) Reset() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if !t.running {
		return
	}

	if t.reset != nil {
		close(t.reset)
		t.reset = nil
	}
}

func (t *toolsUpdater) loop() {
	defer t.Stop()

	// add some jitter. When spinning up multiple entities, we add
	// jitter to prevent stampeeding herd.
	randInt, err := rand.Int(rand.Reader, big.NewInt(3000))
	if err != nil {
		randInt = big.NewInt(0)
	}
	time.Sleep(time.Duration(randInt.Int64()) * time.Millisecond)

	var resetTime time.Time
	now := time.Now().UTC()
	if now.After(t.lastUpdate.Add(40 * time.Minute)) {
		if err := t.updateTools(); err != nil {
			slog.ErrorContext(t.ctx, "initial tools update error", "error", err)
			resetTime = now.Add(5 * time.Minute)
			slog.ErrorContext(t.ctx, "initial tools update error", "error", err)
		} else {
			// Tools are usually valid for 1 hour.
			resetTime = t.lastUpdate.Add(40 * time.Minute)
		}
	}

	for {
		if t.reset == nil {
			t.reset = make(chan struct{})
		}
		// add some jitter
		randInt, err := rand.Int(rand.Reader, big.NewInt(300))
		if err != nil {
			randInt = big.NewInt(0)
		}
		timer := time.NewTimer(resetTime.Sub(now) + time.Duration(randInt.Int64())*time.Second)
		select {
		case <-t.quit:
			slog.DebugContext(t.ctx, "stopping tools updater")
			timer.Stop()
			return
		case <-timer.C:
			slog.DebugContext(t.ctx, "updating tools")
			now = time.Now().UTC()
			if err := t.updateTools(); err == nil {
				slog.ErrorContext(t.ctx, "updating tools", "error", err)
				resetTime = now.Add(5 * time.Minute)
			} else {
				// Tools are usually valid for 1 hour.
				resetTime = t.lastUpdate.Add(40 * time.Minute)
			}
		case <-t.reset:
			slog.DebugContext(t.ctx, "resetting tools updater")
			timer.Stop()
			now = time.Now().UTC()
			if err := t.updateTools(); err != nil {
				slog.ErrorContext(t.ctx, "updating tools", "error", err)
				resetTime = now.Add(5 * time.Minute)
			} else {
				// Tools are usually valid for 1 hour.
				resetTime = t.lastUpdate.Add(40 * time.Minute)
			}
		}
		timer.Stop()
	}
}
