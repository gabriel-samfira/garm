package worker

import (
	dbCommon "github.com/cloudbase/garm/database/common"
)

// handleInstanceWatcherEvent watches for events that involve instances which belong to
// a particular scale set.
func (s *ScaleSetWorker) handleInstanceWatcherEvent(event dbCommon.ChangePayload) {
	// switch event.EntityType {
	// case dbCommon.GithubCredentialsEntityType:
	// 	s.handleGithubCredentialsWatcherEvent(event)
	// case dbCommon.RepositoryEntityType, dbCommon.OrganizationEntityType, dbCommon.EnterpriseEntityType:
	// 	s.handleGithubEntityWatcherEvent(event)
	// case dbCommon.ScaleSetEntityType:
	// 	s.handleScaleSetWatcherEvent(event)
	// default:
	// 	slog.ErrorContext(s.ctx, "invalid entity type", "entity_type", event.EntityType)
	// }
}

func (s *ScaleSetWorker) runScaleSetInstanceWatcher() {
	defer s.scalesetConsumer.Close()
	for {
		select {
		case <-s.quit:
			return
		case <-s.ctx.Done():
			return
		case event, ok := <-s.scalesetInstanceConsumer.Watch():
			if !ok {
				return
			}
			go s.handleInstanceWatcherEvent(event)
		}
	}
}
