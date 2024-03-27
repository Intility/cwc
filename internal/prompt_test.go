package internal_test

import (
	stdErrors "errors"
	"github.com/intility/cwc/internal"
	"github.com/intility/cwc/internal/mocks"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/templates"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArgsOrTemplatePromptResolver_ResolvePromptUsesTemplatePrompt(t *testing.T) {
	mockTemplateProvider := &mocks.TemplateProvider{}
	testTemplate := &templates.Template{DefaultPrompt: "test"}
	mockTemplateProvider.On("GetTemplate", "test").Return(testTemplate, nil)

	resolver := internal.NewArgsOrTemplatePromptResolver(mockTemplateProvider, []string{}, "test")

	prompt := resolver.ResolvePrompt()

	mockTemplateProvider.AssertExpectations(t)
	assert.Equal(t, "test", prompt)
}

func TestArgsOrTemplatePromptResolver_ResolvePromptUsesArgsPrompt(t *testing.T) {
	args := []string{"arg prompt"}
	mockTemplateProvider := &mocks.TemplateProvider{}
	testTemplate := &templates.Template{DefaultPrompt: "template prompt"}
	mockTemplateProvider.On("GetTemplate", "test").Return(testTemplate, nil)
	resolver := internal.NewArgsOrTemplatePromptResolver(mockTemplateProvider, args, "test")

	prompt := resolver.ResolvePrompt()

	mockTemplateProvider.AssertExpectations(t)
	assert.Equal(t, "arg prompt", prompt)
}

func TestArgsOrTemplatePromptResolver_ResolvePromptReturnsEmptyStringWhenTemplateNotFound(t *testing.T) {
	mockTemplateProvider := &mocks.TemplateProvider{}
	err := errors.TemplateNotFoundError{TemplateName: "test"}
	mockTemplateProvider.On("GetTemplate", "test").Return(nil, err)
	resolver := internal.NewArgsOrTemplatePromptResolver(mockTemplateProvider, []string{}, "test")

	prompt := resolver.ResolvePrompt()

	mockTemplateProvider.AssertExpectations(t)
	assert.Equal(t, "", prompt)
}

func TestArgsOrTemplatePromptResolver_ResolvePromptReturnsEmptyStringWhenTemplateProviderFails(t *testing.T) {
	mockTemplateProvider := &mocks.TemplateProvider{}
	err := stdErrors.New("error")
	mockTemplateProvider.On("GetTemplate", "test").Return(nil, err)
	resolver := internal.NewArgsOrTemplatePromptResolver(mockTemplateProvider, []string{}, "test")

	prompt := resolver.ResolvePrompt()

	mockTemplateProvider.AssertExpectations(t)
	assert.Equal(t, "", prompt)
}
