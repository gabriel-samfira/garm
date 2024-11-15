package worker

import (
	dbCommon "github.com/cloudbase/garm/database/common"
	"github.com/cloudbase/garm/database/watcher"
	"github.com/cloudbase/garm/params"
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
	)
}
