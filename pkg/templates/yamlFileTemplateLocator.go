package templates

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/intility/cwc/pkg/errors"
)

type YamlFileTemplateLocator struct {
	// Path is the path to the directory containing the templates
	Path string
}

// configFile is a struct that represents the yaml file containing the templates.
type configFile struct {
	Templates []Template `yaml:"templates"`
}

func NewYamlFileTemplateLocator(path string) *YamlFileTemplateLocator {
	return &YamlFileTemplateLocator{
		Path: path,
	}
}

func (y *YamlFileTemplateLocator) ListTemplates() ([]Template, error) {
	// no configured templates file is a valid state
	// and should not return an error
	_, err := os.Stat(y.Path)
	if os.IsNotExist(err) {
		return []Template{}, nil
	}

	file, err := os.Open(y.Path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	decoder := yaml.NewDecoder(file)

	var cfg configFile
	err = decoder.Decode(&cfg)

	if err != nil {
		return nil, fmt.Errorf("error decoding file: %w", err)
	}

	return cfg.Templates, nil
}

func (y *YamlFileTemplateLocator) GetTemplate(name string) (*Template, error) {
	templates, err := y.ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("error getting template: %w", err)
	}

	for _, tmpl := range templates {
		if tmpl.Name == name {
			return &tmpl, nil
		}
	}

	return nil, errors.TemplateNotFoundError{TemplateName: name}
}
