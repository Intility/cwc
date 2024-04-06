package prompting_test

import (
	stdErrors "errors"
	"github.com/intility/cwc/pkg/prompting"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/intility/cwc/mocks"
	"github.com/intility/cwc/pkg/templates"
)

func TestArgsOrTemplatePromptResolver_ResolvePrompt(t *testing.T) {
	// Define the test cases
	type testConfig struct {
		locator      *mocks.TemplateLocator
		testTemplate *templates.Template
	}

	tests := []struct {
		name         string
		args         []string
		templateName string
		setupMocks   func(testConfig)
		wantResult   func(t *testing.T, result string)
	}{
		{
			name:         "template default prompt",
			args:         []string{},
			templateName: "test",
			setupMocks: func(m testConfig) {
				m.testTemplate.DefaultPrompt = "foo"
				m.locator.On("GetTemplate", "test").Return(m.testTemplate, nil)
			},
			wantResult: func(t *testing.T, prompt string) {
				assert.Equal(t, "foo", prompt)
			},
		},
		{
			name:         "args prompt when error getting template",
			args:         []string{"bar"},
			templateName: "test",
			setupMocks: func(m testConfig) {
				m.locator.On("GetTemplate", "test").Return(nil, stdErrors.New("error"))
			},
			wantResult: func(t *testing.T, prompt string) {
				assert.Equal(t, "bar", prompt)
			},
		},
		{
			name:         "args prompt overrides template",
			args:         []string{"bar"},
			templateName: "test",
			setupMocks: func(m testConfig) {
				m.testTemplate.DefaultPrompt = "foo"
				m.locator.On("GetTemplate", "test").Return(m.testTemplate, nil)
			},
			wantResult: func(t *testing.T, prompt string) {
				assert.Equal(t, "bar", prompt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locator := &mocks.TemplateLocator{}
			cfg := testConfig{locator: locator, testTemplate: &templates.Template{}}
			tt.setupMocks(cfg)

			resolver := prompting.NewArgsOrTemplatePromptResolver(locator, tt.args, tt.templateName)
			prompt := resolver.ResolvePrompt()

			locator.AssertExpectations(t)
			tt.wantResult(t, prompt)
		})
	}
}
