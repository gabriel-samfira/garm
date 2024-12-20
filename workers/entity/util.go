package entity

import (
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
)

func composeWatcherFilters(entity params.GithubEntity) dbCommon.PayloadFilterFunc {
	return watcher.WithAny(
		watcher.WithEntityFilter(entity),
		// Watch for changes to pools belonging to this entity.
		watcher.WithEntityPoolFilter(entity),
		// Watch for changes to scale sets belonging to this entity.
		watcher.WithEntityScaleSetFilter(entity),
		// Watch for credentials updates.
		watcher.WithGithubCredentialsFilter(entity.Credentials),
		// Watch for jobs meant for this entity.
		watcher.WithEntityJobFilter(entity),
		watcher.WithAll(
			// Updates to the controller
			watcher.WithEntityTypeFilter(dbCommon.ControllerEntityType),
			watcher.WithOperationTypeFilter(dbCommon.UpdateOperation),
		),
	)
}
