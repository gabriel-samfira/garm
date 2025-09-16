package config

import (
	"fmt"
	"net/url"

	"github.com/BurntSushi/toml"
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
	WorkDir   string `toml:"work_dir"`
}

func (a *Agent) Validate() error {
	if _, err := url.ParseRequestURI(a.ServerURL); err != nil {
		return fmt.Errorf("invalid server_url: %w", err)
	}
	return nil
}
