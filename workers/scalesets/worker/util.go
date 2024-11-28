package worker

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner/common"
	garmUtil "github.com/cloudbase/garm/util"
)

func composeWatcherFilters(scaleSet params.ScaleSet, ghEntity params.GithubEntity) dbCommon.PayloadFilterFunc {
	return watcher.WithAny(
		// Watch for changes on the scale set. We use this to update the
		// internal state of the entity.
		watcher.WithScaleSetFilter(scaleSet),
		// Watch for changes on credentials. If credentials get updated,
		// we need to refresh the clients.
		watcher.WithGithubCredentialsFilter(ghEntity.Credentials),
		// Watch for changes on the github entity (repo, org, enterprise). Credentials
		// for the entity may be swapped out completely, in which case, we need to refresh
		// the clients.
		watcher.WithEntityFilter(ghEntity),
		// Watch for controller changes.
		watcher.WithAll(
			// Updates to the controller
			watcher.WithEntityTypeFilter(dbCommon.ControllerEntityType),
			watcher.WithOperationTypeFilter(dbCommon.UpdateOperation),
		),
	)
}

func (s *ScaleSetWorker) startLoopForFunction(f func() error, interval time.Duration, name string, alwaysRun bool) {
	slog.InfoContext(
		s.ctx, "starting loop for entity",
		"loop_name", name)
	jitter := rand.IntN(10000)
	ticker := time.NewTicker(interval + time.Duration(jitter)*time.Millisecond)
	s.wg.Add(1)

	defer func() {
		slog.InfoContext(
			s.ctx, "pool loop exited",
			"loop_name", name)
		ticker.Stop()
		s.wg.Done()
	}()

	for {
		shouldRun := s.faultReason == nil
		if alwaysRun {
			shouldRun = true
		}
		switch shouldRun {
		case true:
			select {
			case <-ticker.C:
				if err := f(); err != nil {
					slog.With(slog.Any("error", err)).ErrorContext(
						s.ctx, "error in loop",
						"loop_name", name)
				}
			case <-s.ctx.Done():
				// daemon is shutting down.
				return
			case <-s.quit:
				// this worker was stopped.
				return
			}
		default:
			select {
			case <-s.ctx.Done():
				// daemon is shutting down.
				return
			case <-s.quit:
				// this worker was stopped.
				return
			default:
				s.waitForTimeoutOrCancelled(common.BackoffTimer)
			}
		}
	}
}

func (s *ScaleSetWorker) waitForTimeoutOrCancelled(timeout time.Duration) {
	slog.DebugContext(
		s.ctx, fmt.Sprintf("sleeping for %.2f minutes", timeout.Minutes()))
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-s.ctx.Done():
	case <-s.quit:
	}
}

func (s *ScaleSetWorker) updateTools() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	// Update tools cache.
	tools, err := garmUtil.FetchTools(s.ctx, s.ghCli)
	if err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(
			s.ctx, "failed to update tools for entity", "entity", s.entity.String())
		s.faultReason = err
		return fmt.Errorf("failed to update tools for entity %s: %w", s.entity.String(), err)
	}
	s.mux.Lock()
	s.tools = tools
	s.mux.Unlock()

	slog.DebugContext(s.ctx, "successfully updated tools")
	s.faultReason = nil
	return err
}
