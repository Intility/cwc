package config_test

import (
	"errors"
	"github.com/intility/cwc/pkg/config"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"

	"github.com/intility/cwc/mocks"
)

func TestNewClientFromConfig(t *testing.T) {

	// Define the test cases
	type testConfig struct {
		cfgProvider  *mocks.ConfigProvider
		clientConfig openai.ClientConfig
	}

	tests := []struct {
		name       string
		setupMocks func(testConfig)
		wantResult func(t *testing.T, result *openai.Client)
		wantErr    func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupMocks: func(m testConfig) {
				m.cfgProvider.On("NewFromConfigFile").Return(m.clientConfig, nil)
			},
			wantResult: func(t *testing.T, result *openai.Client) {
				assert.NotNil(t, result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "error loading config",
			setupMocks: func(m testConfig) {
				m.cfgProvider.On("NewFromConfigFile").
					Return(openai.ClientConfig{}, errors.New("error reading config"))
			},
			wantResult: func(t *testing.T, result *openai.Client) {
				assert.Nil(t, result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfigProvider := &mocks.ConfigProvider{}
			cfg := testConfig{cfgProvider: mockConfigProvider, clientConfig: openai.ClientConfig{}}
			tt.setupMocks(cfg)

			clientProvider := config.NewOpenAIClientProvider(mockConfigProvider)
			res, err := clientProvider.NewClientFromConfig()

			mockConfigProvider.AssertExpectations(t)
			tt.wantResult(t, res)
			tt.wantErr(t, err)
		})

	}
}
