package systemcontext

import (
	"fmt"
	"strings"
	tt "text/template"

	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/templates"
)

const (
	defaultTemplateName = "default"
)

type SystemMessageGenerator interface {
	GenerateSystemMessage() (string, error)
}

type TemplatedSystemMessageGenerator struct {
	templateLocator  templates.TemplateLocator
	templateName     string
	templateVars     map[string]string
	contextRetriever ContextRetriever
}

func NewTemplatedSystemMessageGenerator(
	templateLocator templates.TemplateLocator,
	templateName string,
	templateVars map[string]string,
	contextRetriever ContextRetriever,
) *TemplatedSystemMessageGenerator {
	return &TemplatedSystemMessageGenerator{
		templateLocator:  templateLocator,
		templateName:     templateName,
		templateVars:     templateVars,
		contextRetriever: contextRetriever,
	}
}

func (smg *TemplatedSystemMessageGenerator) GenerateSystemMessage() (string, error) {
	ctx, err := smg.contextRetriever.RetrieveContext()
	if err != nil {
		return "", fmt.Errorf("error retrieving context: %w", err)
	}

	tmpl, err := smg.templateLocator.GetTemplate(smg.templateName)

	if smg.templateVars == nil {
		smg.templateVars = make(map[string]string)
	}

	// if no template found, create a basic template as fallback
	if err != nil {
		if errors.IsTemplateNotFoundError(err) {
			// exit with error if the user has requested a custom template and it is not found
			if smg.templateName != defaultTemplateName {
				return "", fmt.Errorf("template not found: %w", err)
			}

			return CreateBuiltinSystemMessageFromContext(ctx), nil
		}

		return "", fmt.Errorf("error getting template: %w", err)
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

func CreateBuiltinSystemMessageFromContext(ctx string) string {
	var systemMessage strings.Builder

	systemMessage.WriteString("You are a helpful coding assistant. ")
	systemMessage.WriteString("Below you will find relevant context to answer the user's question.\n\n")
	systemMessage.WriteString("Context:\n")
	systemMessage.WriteString(ctx)
	systemMessage.WriteString("\n\n")
	systemMessage.WriteString("Please follow the users instructions, you can do this!")

	return systemMessage.String()
}
