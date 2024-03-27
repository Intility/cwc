package internal

import (
	"fmt"
	"path/filepath"

	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/templates"
)

type TemplateProvider interface {
	GetTemplate(templateName string) (*templates.Template, error)
}

type DefaultTemplateProvider struct {
	configProvider ConfigProvider
}

func NewTemplateProvider(cfgProvider ConfigProvider) *DefaultTemplateProvider {
	return &DefaultTemplateProvider{configProvider: cfgProvider}
}

func (tp *DefaultTemplateProvider) GetTemplate(templateName string) (*templates.Template, error) {
	if templateName == "" {
		templateName = "default"
	}

	var locators []templates.TemplateLocator

	configDir, err := tp.configProvider.GetConfigDir()
	if err == nil {
		locators = append(locators, templates.NewYamlFileTemplateLocator(filepath.Join(configDir, "templates.yaml")))
	}

	locators = append(locators, templates.NewYamlFileTemplateLocator(filepath.Join(".cwc", "templates.yaml")))
	mergedLocator := templates.NewMergedTemplateLocator(locators...)

	tmpl, err := mergedLocator.GetTemplate(templateName)
	if err != nil {
		if errors.IsTemplateNotFoundError(err) {
			return nil, errors.TemplateNotFoundError{TemplateName: templateName}
		}

		return nil, fmt.Errorf("error getting template: %w", err)
	}

	return tmpl, nil
}
