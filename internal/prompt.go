package internal

import "github.com/intility/cwc/pkg/templates"

type PromptResolver interface {
	ResolvePrompt() string
}

type ArgsOrTemplatePromptResolver struct {
	args            []string
	templateName    string
	templateLocator templates.TemplateLocator
}

func NewArgsOrTemplatePromptResolver(
	templateLocator templates.TemplateLocator,
	args []string,
	tmplName string,
) *ArgsOrTemplatePromptResolver {
	return &ArgsOrTemplatePromptResolver{
		args:            args,
		templateName:    tmplName,
		templateLocator: templateLocator,
	}
}

func (r *ArgsOrTemplatePromptResolver) ResolvePrompt() string {
	var prompt string

	tmpl, err := r.templateLocator.GetTemplate(r.templateName)
	if err == nil && tmpl.DefaultPrompt != "" {
		prompt = tmpl.DefaultPrompt
	}

	if len(r.args) > 0 {
		prompt = r.args[0]
	}

	return prompt
}
