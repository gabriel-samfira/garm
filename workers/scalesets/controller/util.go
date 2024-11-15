package controller

import (
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
)

func composeWatcherFilters(entity params.GithubEntity) dbCommon.PayloadFilterFunc {
	// We want to watch for changes in either the controller or the
	// entity itself.
	return watcher.WithAny(
		watcher.WithAll(
			// Updates to the controller
			watcher.WithEntityTypeFilter(dbCommon.ControllerEntityType),
			watcher.WithOperationTypeFilter(dbCommon.UpdateOperation),
		),
		// Any operation on the entity we're managing the pool for.
		watcher.WithEntityFilter(entity),
		// Watch for changes to the github credentials
		watcher.WithGithubCredentialsFilter(entity.Credentials),
	)
}
