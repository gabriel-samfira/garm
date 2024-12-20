package entity

import (
	"fmt"
	"sync"

	commonParams "github.com/cloudbase/garm-provider-common/params"

	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner/common"
	"github.com/cloudbase/garm/util/github/scalesets"
)

type githubResources struct {
	tools       []commonParams.RunnerApplicationDownload
	credentials params.GithubCredentials

	err error

	ghCli          common.GithubClient
	scaleSetClient *scalesets.ScaleSetClient

	mux sync.Mutex
}

func (g *githubResources) GetGithubClient() (common.GithubClient, error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	if g.err != nil {
		return nil, g.err
	}

	if g.ghCli == nil {
		return nil, fmt.Errorf("github client not initialized")
	}

	return g.ghCli, nil
}

func (g *githubResources) GetScaleSetClient() (*scalesets.ScaleSetClient, error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	if g.err != nil {
		return nil, g.err
	}

	if g.scaleSetClient == nil {
		return nil, fmt.Errorf("scaleset client not initialized")
	}
	return g.scaleSetClient, nil
}

func (g *githubResources) GetTools() ([]commonParams.RunnerApplicationDownload, error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	if g.err != nil {
		return nil, g.err
	}

	return g.tools, nil
}
