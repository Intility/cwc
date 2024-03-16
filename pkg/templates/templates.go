package templates

type Template struct {
	// Name is the name of the template
	Name string `yaml:"name"`

	// Description is a short description of the template
	Description string `yaml:"description"`

	// DefaultPrompt is the prompt that is used if no prompt is provided
	DefaultPrompt string `yaml:"defaultPrompt,omitempty"`

	// SystemMessage is the message that primes the conversation
	SystemMessage string `yaml:"systemMessage"`

	// Variables is a list of input variables for the template
	Variables []TemplateVariable `yaml:"variables"`
}

type TemplateVariable struct {
	// Name is the name of the input variable
	Name string `yaml:"name"`

	// Description is a short description of the input variable
	Description string `yaml:"description"`

	// DefaultValue is the value used if no override is provided
	DefaultValue string `yaml:"defaultValue,omitempty"`
}

type TemplateLocator interface {
	// ListTemplates returns a list of available templates
	ListTemplates() ([]Template, error)

	// GetTemplate returns a template by name
	GetTemplate(name string) (*Template, error)
}
