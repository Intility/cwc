package client

import (
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/internal/config"
)

type Provider interface {
	NewClientFromConfig() (*openai.Client, error)
}

type OpenAIClientProvider struct {
	cfg config.Provider
}

func NewOpenAIClientProvider(provider config.Provider) *OpenAIClientProvider {
	return &OpenAIClientProvider{cfg: provider}
}

func (c *OpenAIClientProvider) NewClientFromConfig() (*openai.Client, error) {
	cfg, err := c.cfg.NewFromConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	client := openai.NewClientWithConfig(cfg)

	return client, nil
}
