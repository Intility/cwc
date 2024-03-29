package internal

import (
	"fmt"
	"github.com/intility/cwc/pkg/config"

	"github.com/sashabaranov/go-openai"
)

type ClientProvider interface {
	NewClientFromConfig() (*openai.Client, error)
}

type OpenAIClientProvider struct {
	cfg config.ConfigProvider
}

func NewOpenAIClientProvider(provider config.ConfigProvider) *OpenAIClientProvider {
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
