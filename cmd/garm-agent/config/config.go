package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/cloudbase/garm/auth"
	"github.com/cloudbase/garm/params"
)

// NewConfig returns a new Config
func NewConfig(cfgFile string) (*Agent, error) {
	var config Agent
	if _, err := toml.DecodeFile(cfgFile, &config); err != nil {
		return nil, fmt.Errorf("error decoding toml: %w", err)
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}
	return &config, nil
}

type Agent struct {
	ServerURL string `toml:"server_url"`
	Token     string `toml:"token"`
	// WorkDir is the folder in which the runner will execute workflows.
	// In the case of github, it is expected that the WorkDir is also the installation
	// dir of the runner.
	WorkDir string `toml:"work_dir"`
	LogFile string `toml:"log_file"`
	Shell   string `toml:"shell"`
	// RunnerExecutable is the absolute path on disk to the runner executable. This can be
	// any executable (bash, binary, etc). For github it will most likely be the papth to
	// run.sh, which is usually in the WorkDir folder. For gitea the binary can be anywhere.
	RunnerExecutable string `toml:"runner_executable"`
}

func (a *Agent) Validate() error {
	if _, err := url.ParseRequestURI(a.ServerURL); err != nil {
		return fmt.Errorf("invalid server_url: %w", err)
	}

	if a.Token == "" {
		return fmt.Errorf("missing token")
	}

	if _, err := a.TokenClaims(); err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if a.RunnerExecutable == "" {
		return fmt.Errorf("runner executable path is not set")
	}
	mode, err := os.Stat(a.RunnerExecutable)
	if err != nil {
		return fmt.Errorf("failed to access runner executable %s: %w", a.RunnerExecutable, err)
	}

	if mode.IsDir() {
		return fmt.Errorf("runner executable seems to be a directory")
	}

	return nil
}

func (a *Agent) TokenClaims() (auth.InstanceJWTClaims, error) {
	claims := auth.InstanceJWTClaims{}
	_, err := jwt.ParseWithClaims(a.Token, &claims, nil)
	if err != nil && !errors.Is(err, jwt.ErrTokenUnverifiable) {
		return auth.InstanceJWTClaims{}, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	if !claims.IsAgent {
		return auth.InstanceJWTClaims{}, fmt.Errorf("token is not agent scoped")
	}
	return claims, nil
}

func (a *Agent) ForgeType() (params.EndpointType, error) {
	claims, err := a.TokenClaims()
	if err != nil {
		return "", fmt.Errorf("failed to get token claims: %w", err)
	}
	forgeType := params.EndpointType(claims.ForgeType)
	switch forgeType {
	case params.GithubEndpointType, params.GiteaEndpointType:
	default:
		return "", fmt.Errorf("invalid forge type: %s", forgeType)
	}
	return forgeType, nil
}
