package systemcontext_test

import (
	"github.com/intility/cwc/pkg/systemcontext"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/intility/cwc/mocks"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/templates"
)

func TestTemplatedSystemMessageGenerator_GenerateSystemMessage(t *testing.T) {
	type testConfig struct {
		locator      *mocks.TemplateLocator
		testTemplate *templates.Template
		templateVars map[string]string
		ctxRetriever *mocks.ContextRetriever
	}

	tests := []struct {
		name       string
		setupMocks func(testConfig)
		wantResult func(t *testing.T, result string)
		wantErr    func(t *testing.T, err error)
	}{
		{
			name: "use builtin system message if no template found",
			setupMocks: func(m testConfig) {
				m.ctxRetriever.On("RetrieveContext").Return("test_context", nil)
				m.locator.On("GetTemplate", "test").
					Return(nil, errors.TemplateNotFoundError{})
			},
			wantResult: func(t *testing.T, result string) {
				builtInMessage := systemcontext.CreateBuiltinSystemMessageFromContext("test_context")
				assert.Equal(t, builtInMessage, result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "returns error if template provider fails",
			setupMocks: func(m testConfig) {
				m.ctxRetriever.On("RetrieveContext").Return("test_context", nil)
				m.locator.On("GetTemplate", "test").
					Return(nil, assert.AnError)
			},
			wantResult: func(t *testing.T, result string) {
				assert.Empty(t, result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "render template without vars",
			setupMocks: func(m testConfig) {
				m.ctxRetriever.On("RetrieveContext").Return("test_context", nil)
				m.testTemplate = &templates.Template{SystemMessage: "test_message"}
				m.locator.On("GetTemplate", "test").
					Return(m.testTemplate, nil)
			},
			wantResult: func(t *testing.T, result string) {
				assert.Equal(t, "test_message", result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "render template with default var values",
			setupMocks: func(m testConfig) {
				m.ctxRetriever.On("RetrieveContext").Return("test_context", nil)
				m.testTemplate = &templates.Template{
					SystemMessage: "test_message {{.Variables.foo}}",
					Variables: []templates.TemplateVariable{
						{Name: "foo", DefaultValue: "bar"},
					},
				}
				m.locator.On("GetTemplate", "test").
					Return(m.testTemplate, nil)
			},
			wantResult: func(t *testing.T, result string) {
				assert.Equal(t, "test_message bar", result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "render template with replaced var values",
			setupMocks: func(m testConfig) {
				m.ctxRetriever.On("RetrieveContext").Return("test_context", nil)
				m.testTemplate = &templates.Template{
					SystemMessage: "test_message {{.Variables.foo}}",
					Variables: []templates.TemplateVariable{
						{Name: "foo", DefaultValue: "bar"},
					},
				}
				m.locator.On("GetTemplate", "test").
					Return(m.testTemplate, nil)
				m.templateVars["foo"] = "baz"
			},
			wantResult: func(t *testing.T, result string) {
				assert.Equal(t, "test_message baz", result)
			},
			wantErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locator := &mocks.TemplateLocator{}
			ctxRetriever := &mocks.ContextRetriever{}
			cfg := testConfig{
				locator:      locator,
				testTemplate: &templates.Template{},
				templateVars: map[string]string{},
				ctxRetriever: ctxRetriever,
			}

			tt.setupMocks(cfg)

			smg := systemcontext.NewTemplatedSystemMessageGenerator(
				locator,
				"test",
				cfg.templateVars,
				ctxRetriever,
			)

			res, err := smg.GenerateSystemMessage()

			locator.AssertExpectations(t)
			tt.wantResult(t, res)
			tt.wantErr(t, err)
		})
	}
}
