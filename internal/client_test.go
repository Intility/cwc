package internal_test

import (
	"errors"
	"github.com/intility/cwc/internal"
	"github.com/intility/cwc/internal/mocks"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientFromConfig(t *testing.T) {
	mockConfigProvider := &mocks.ConfigProvider{}
	testClientConfig := openai.ClientConfig{BaseURL: "http://test"}
	mockConfigProvider.On("NewFromConfigFile").Return(testClientConfig, nil)

	clientProvider := internal.NewOpenAIClientProvider(mockConfigProvider)
	client, err := clientProvider.NewClientFromConfig()

	mockConfigProvider.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClientFromConfigError(t *testing.T) {
	mockConfigProvider := &mocks.ConfigProvider{}
	mockConfigProvider.On("NewFromConfigFile").Return(openai.ClientConfig{}, errors.New("error reading config"))

	clientProvider := internal.NewOpenAIClientProvider(mockConfigProvider)
	_, err := clientProvider.NewClientFromConfig()

	mockConfigProvider.AssertExpectations(t)
	assert.Error(t, err)
}
