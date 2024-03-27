package prompt

import (
	stdErrors "errors"

	"github.com/intility/cwc/internal/template"
	"github.com/intility/cwc/pkg/errors"
)

type Resolver interface {
	ResolvePrompt() string
}

type ArgsOrTemplatePromptResolver struct {
	args             []string
	templateName     string
	templateProvider template.Provider
}

func NewArgsOrTemplatePromptResolver(
	tmplProvider template.Provider,
	args []string,
	tmplName string,
) *ArgsOrTemplatePromptResolver {
	return &ArgsOrTemplatePromptResolver{
		args:             args,
		templateName:     tmplName,
		templateProvider: tmplProvider,
	}
}

func (r *ArgsOrTemplatePromptResolver) ResolvePrompt() string {
	var prompt string

	tmpl, err := r.templateProvider.GetTemplate(r.templateName)
	if err != nil {
		var templateNotFoundError errors.TemplateNotFoundError
		if stdErrors.As(err, &templateNotFoundError) {
			if len(r.args) == 0 {
				return ""
			}
		}
	} else {
		prompt = tmpl.DefaultPrompt
	}

	if len(r.args) > 0 {
		prompt = r.args[0]
	}

	return prompt
}
