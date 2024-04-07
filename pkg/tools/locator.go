package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"gopkg.in/yaml.v3"

	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/ui"
)

type ConfigFileToolLocator struct {
	// Paths to search for config files
	Paths []string
}

func NewConfigFileToolLocator(paths ...string) *ConfigFileToolLocator {
	return &ConfigFileToolLocator{
		Paths: paths,
	}
}

type ConfiguredTool struct {
	Name        string                    `yaml:"name"`
	Description string                    `yaml:"description"`
	Shell       []string                  `yaml:"shell"`
	Web         []string                  `yaml:"web"`
	Parameters  []ConfiguredToolParameter `yaml:"parameters"`
}

type ConfiguredToolParameter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
}

type ToolConfig struct {
	Tools []ConfiguredTool `yaml:"tools"`
}

type ToolParser struct {
	dataTypes map[string]jsonschema.DataType
}

func NewToolParser() *ToolParser {
	return &ToolParser{
		dataTypes: map[string]jsonschema.DataType{
			"string":  jsonschema.String,
			"integer": jsonschema.Integer,
			"number":  jsonschema.Number,
			"boolean": jsonschema.Boolean,
			"object":  jsonschema.Object,
			"array":   jsonschema.Array,
		},
	}
}

func (t *ToolParser) ParseArguments(args []ConfiguredToolParameter) (map[string]jsonschema.Definition, error) {
	propertiesMap := make(map[string]jsonschema.Definition)

	for _, param := range args {
		dataType, ok := t.dataTypes[param.Type]
		if !ok {
			supportedTypes := make([]string, 0, len(t.dataTypes))
			for k := range t.dataTypes {
				supportedTypes = append(supportedTypes, k)
			}

			return nil, &errors.InvalidToolSpecError{
				Message: fmt.Sprintf("tool '%s' has an invalid data type: %s. supported types: %s",
					param.Name, param.Type, strings.Join(supportedTypes, ", ")),
			}
		}

		propertiesMap[param.Name] = jsonschema.Definition{
			Type:        dataType,
			Description: param.Description,
		}
	}

	return propertiesMap, nil
}

func (c *ConfigFileToolLocator) LocateTool(toolID string) *Tool {
	// loop the paths in descending order, so that the last path takes precedence
	argsParser := NewToolParser()

	for i := len(c.Paths) - 1; i >= 0; i-- {
		file, err := os.ReadFile(c.Paths[i])
		if err != nil {
			continue
		}

		var config ToolConfig

		err = yaml.Unmarshal(file, &config)
		if err != nil {
			continue
		}

		for _, configuredTool := range config.Tools {
			if configuredTool.Name != toolID {
				continue
			}

			propertiesMap, err := argsParser.ParseArguments(configuredTool.Parameters)
			if err != nil {
				ui.PrintMessage(err.Error(), ui.MessageTypeError)
				continue
			}

			// convert the ConfiguredTool to a Tool
			var toolDef openai.FunctionDefinition
			toolDef.Name = configuredTool.Name
			toolDef.Description = configuredTool.Description
			toolDef.Parameters = jsonschema.Definition{
				Type:       jsonschema.Object,
				Properties: propertiesMap,
			}

			tool := &Tool{definition: toolDef, executor: c.getExecutor(configuredTool)}

			return tool
		}
	}

	return nil
}

func (c *ConfigFileToolLocator) getExecutor(tool ConfiguredTool) ToolExecutor { //nolint:ireturn
	if len(tool.Shell) > 0 {
		return NewShellExecutor(tool.Shell)
	}

	if len(tool.Web) > 0 {
		return NewWebExecutor(tool.Web)
	}

	return nil
}

type MockLocator struct{}

func (m *MockLocator) LocateTool(id string) *Tool {
	toolDef := openai.FunctionDefinition{
		Name:        "diff",
		Description: "Get the git diff between two refs.",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"source": {
					Type:        jsonschema.String,
					Description: "The source ref to compare. E.g. 'main', 'HEAD~1', etc.",
				},
				"target": {
					Type:        jsonschema.String,
					Description: "The target ref to compare. E.g. 'main', 'HEAD~1', etc.",
				},
			},
			Required: []string{"target"},
		},
	}

	return &Tool{
		definition: toolDef,
		executor:   NewShellExecutor([]string{`git diff {{or .source ""}} {{.target}}`}),
	}
}
