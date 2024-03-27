package internal

import (
	"fmt"
	"strings"
	tt "text/template"

	"github.com/intility/cwc/pkg/errors"
)

type SystemMessageGenerator interface {
	GenerateSystemMessage(ctx string) (string, error)
}

type TemplatedSystemMessageGenerator struct {
	templateProvider TemplateProvider
	templateName     string
	templateVars     map[string]string
}

func NewTemplatedSystemMessageGenerator(
	templateProvider TemplateProvider,
	templateName string,
	templateVars map[string]string,
) *TemplatedSystemMessageGenerator {
	return &TemplatedSystemMessageGenerator{
		templateProvider: templateProvider,
		templateName:     templateName,
		templateVars:     templateVars,
	}
}

func (smg *TemplatedSystemMessageGenerator) GenerateSystemMessage(ctx string) (string, error) {
	tmpl, err := smg.templateProvider.GetTemplate(smg.templateName)

	if smg.templateVars == nil {
		smg.templateVars = make(map[string]string)
	}

	// if no template found, create a basic template as fallback
	if errors.IsTemplateNotFoundError(err) {
		return smg.createBuiltinSystemMessageFromContext(ctx), nil
	}

	// compile the template.SystemMessage as a go template
	compiledTemplate, err := tt.New("systemMessage").Parse(tmpl.SystemMessage)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	type valueBag struct {
		Context   string
		Variables map[string]string
	}

	// populate the variables map with default values if not provided
	for _, v := range tmpl.Variables {
		if _, ok := smg.templateVars[v.Name]; !ok {
			smg.templateVars[v.Name] = v.DefaultValue
		}
	}

	values := valueBag{
		Context:   ctx,
		Variables: smg.templateVars,
	}

	writer := &strings.Builder{}
	err = compiledTemplate.Execute(writer, values)

	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return writer.String(), nil
}

func (smg *TemplatedSystemMessageGenerator) createBuiltinSystemMessageFromContext(ctx string) string {
	var systemMessage strings.Builder

	systemMessage.WriteString("You are a helpful coding assistant. ")
	systemMessage.WriteString("Below you will find relevant context to answer the user's question.\n\n")
	systemMessage.WriteString("Context:\n")
	systemMessage.WriteString(ctx)
	systemMessage.WriteString("\n\n")
	systemMessage.WriteString("Please follow the users instructions, you can do this!")

	return systemMessage.String()
}
