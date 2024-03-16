package templates

import (
	stdErrors "errors"
	"fmt"

	"github.com/intility/cwc/pkg/errors"
)

// MergedTemplateLocator is a TemplateLocator that merges templates from multiple locators
// making the last applied locator the one that takes precedence in case of name conflicts.
type MergedTemplateLocator struct {
	locators []TemplateLocator
}

// NewMergedTemplateLocator creates a new MergedTemplateLocator.
func NewMergedTemplateLocator(locators ...TemplateLocator) *MergedTemplateLocator {
	return &MergedTemplateLocator{
		locators: locators,
	}
}

// ListTemplates returns a list of available templates.
func (c *MergedTemplateLocator) ListTemplates() ([]Template, error) {
	// Merge templates from all locators
	templates := make(map[string]Template)

	for _, l := range c.locators {
		t, err := l.ListTemplates()
		if err != nil {
			return nil, fmt.Errorf("error listing templates: %w", err)
		}

		for _, template := range t {
			templates[template.Name] = template
		}
	}

	mergedTemplates := make([]Template, 0, len(templates))
	for _, t := range templates {
		mergedTemplates = append(mergedTemplates, t)
	}

	return mergedTemplates, nil
}

// GetTemplate returns a template by name.
func (c *MergedTemplateLocator) GetTemplate(name string) (*Template, error) {
	// Get template from the last locator that has it
	for i := len(c.locators) - 1; i >= 0; i-- {
		tmpl, err := c.locators[i].GetTemplate(name)

		// if template not found, continue to the next locator
		var templateNotFoundError errors.TemplateNotFoundError
		if stdErrors.As(err, &templateNotFoundError) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("error getting template: %w", err)
		}

		return tmpl, nil
	}

	return nil, errors.TemplateNotFoundError{TemplateName: name}
}
