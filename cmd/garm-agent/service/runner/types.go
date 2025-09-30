package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/cloudbase/garm/params"
)

var (
	giteaJobStartedRegex  = regexp.MustCompile(" task [0-9]+ repo is .*")
	githubJobStartedRegex = regexp.MustCompile(" Job started .*")

	githubListenForJobs = regexp.MustCompile("Listening for Jobs")
	giteaListenForJobs  = regexp.MustCompile("runner: .*, declare successfully")
)

type Worker interface {
	Start() error
	Stop() error
	Wait() <-chan error
}

type RunnerStateManager interface {
	SetRunnerStarted(st bool)
	SetJobStarted()
	SetJobFinished()
}

type RunnerConfig interface {
	GetAgentID() uint
	GetAgentName() string
	IsEphemeral() bool
	GetServerURL() string
}

func (r *runnerCmd) isJobStartedLine(msg []byte) bool {
	switch params.EndpointType(r.forgeType) {
	case params.GiteaEndpointType:
		return giteaJobStartedRegex.Match(msg)
	case params.GithubEndpointType:
		return githubJobStartedRegex.Match(msg)
	default:
		return false
	}
}

func (r *runnerCmd) isRunnerStartedLine(msg []byte) bool {
	switch params.EndpointType(r.forgeType) {
	case params.GiteaEndpointType:
		return giteaListenForJobs.Match(msg)
	case params.GithubEndpointType:
		return githubListenForJobs.Match(msg)
	default:
		return false
	}
}

func NewRunnerConfig(cfg string, forgeType params.EndpointType) (RunnerConfig, error) {
	data, err := os.ReadFile(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read runner config: %w", err)
	}
	var runCfg RunnerConfig
	switch forgeType {
	case params.GiteaEndpointType:
		var giteaCfg GiteaRunnerConfig
		if err = json.Unmarshal(data, &giteaCfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal gitea runner config: %w", err)
		}
		runCfg = giteaCfg
	case params.GithubEndpointType:
		var githubCfg GitHubRunnerConfig
		if err = json.Unmarshal(data, &githubCfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal gitea runner config: %w", err)
		}
		runCfg = githubCfg
	default:
		return nil, fmt.Errorf("unknown forge type %s", forgeType)
	}
	return runCfg, nil
}

type GitHubRunnerConfig struct {
	AgentID   uint   `json:"agentId"`
	AgentName string `json:"agentName"`
	Ephemeral bool   `json:"ephemeral"`
	ServerURL string `json:"serverUrl"`
}

func (r GitHubRunnerConfig) GetAgentID() uint {
	return r.AgentID
}

func (r GitHubRunnerConfig) GetAgentName() string {
	return r.AgentName
}
func (r GitHubRunnerConfig) IsEphemeral() bool {
	return r.Ephemeral
}
func (r GitHubRunnerConfig) GetServerURL() string {
	return r.ServerURL
}

type GiteaRunnerConfig struct {
	AgentID   uint   `json:"id"`
	AgentName string `json:"name"`
	Ephemeral bool   `json:"ephemeral"`
	ServerURL string `json:"address"`
}

func (r GiteaRunnerConfig) GetAgentID() uint {
	return r.AgentID
}

func (r GiteaRunnerConfig) GetAgentName() string {
	return r.AgentName
}
func (r GiteaRunnerConfig) IsEphemeral() bool {
	return r.Ephemeral
}
func (r GiteaRunnerConfig) GetServerURL() string {
	return r.ServerURL
}
