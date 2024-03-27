package config

import (
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/config"
)

type Provider interface {
	LoadConfig() (*config.Config, error)
	NewFromConfigFile() (openai.ClientConfig, error)
	GetConfigDir() (string, error)
}

type DefaultProvider struct{}

func NewDefaultProvider() *DefaultProvider {
	return &DefaultProvider{}
}

func (c *DefaultProvider) LoadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return cfg, nil
}

func (c *DefaultProvider) NewFromConfigFile() (openai.ClientConfig, error) {
	cfg, err := config.NewFromConfigFile()
	if err != nil {
		return openai.ClientConfig{}, fmt.Errorf("error reading config: %w", err)
	}

	return cfg, nil
}

func (c *DefaultProvider) GetConfigDir() (string, error) {
	cfgDir, err := config.GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("error getting config dir: %w", err)
	}

	return cfgDir, nil
}
