package entity

import (
	"fmt"
	"log/slog"

	garmUtil "github.com/cloudbase/garm/util"
)

func (s *EntityWorker) updateTools() error {
	// Update tools cache.
	tools, err := garmUtil.FetchTools(r.ctx, r.ghcli)
	if err != nil {
		slog.With(slog.Any("error", err)).ErrorContext(
			s.ctx, "failed to update tools for entity", "entity", s.ghEntity.String())
		s.setPoolRunningState(false, err.Error())
		return fmt.Errorf("failed to update tools for entity %s: %w", s.ghEntity.String(), err)
	}
	s.mux.Lock()
	s.tools = tools
	s.mux.Unlock()

	slog.DebugContext(s.ctx, "successfully updated tools")
	s.setPoolRunningState(true, "")
	return err
}
