package scaleset

import (
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/util/github/scalesets"
)

type scaleSetHelper interface {
	ScaleSetCLI() *scalesets.ScaleSetClient
	GetScaleSet() params.ScaleSet
	SetLastMessageID(id int64) error
	SetDesiredRunnerCount(count int) error
	Owner() string
	HandleJobsCompleted(jobs []params.ScaleSetJobMessage) error
	HandleJobsStarted(jobs []params.ScaleSetJobMessage) error
	HandleJobsAvailable(jobs []params.ScaleSetJobMessage) error
}
