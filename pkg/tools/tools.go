package tools

import (
	"path/filepath"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/config"
)

type ToolLocator interface {
	LocateTool(id string) *Tool
}

type ToolExecutor interface {
	Execute(tool Tool, args any) (string, error)
}

type Tool struct {
	definition       openai.FunctionDefinition
	shellExecutables []string
	webExecutables   []string
}

func (t *Tool) ShellExecutables() []string {
	return t.shellExecutables
}

func (t *Tool) HasShellExecutables() bool {
	return len(t.shellExecutables) > 0
}

func (t *Tool) WebExecutables() []string {
	return t.webExecutables
}

func (t *Tool) HasWebExecutables() bool {
	return len(t.webExecutables) > 0
}

func (t *Tool) Definition() openai.FunctionDefinition {
	return t.definition
}

type Toolkit struct {
	locator      ToolLocator
	enabledTools map[string]*Tool
}

func NewToolkitFromConfigFile(toolIDs ...string) *Toolkit {
	var paths []string

	// add global dir first and local dir second to allow local to override global
	cfgDir, err := config.GetConfigDir()
	if err == nil {
		paths = append(paths, filepath.Join(cfgDir, "tools.yaml"))
	}

	paths = append(paths, filepath.Join(".cwc", "tools.yaml"))

	toolkit := &Toolkit{
		locator:      NewConfigFileToolLocator(paths...),
		enabledTools: make(map[string]*Tool, len(toolIDs)),
	}

	toolkit.initTools(toolIDs...)

	return toolkit
}

func (t *Toolkit) initTools(toolsIDs ...string) {
	for _, toolID := range toolsIDs {
		tool := t.locator.LocateTool(toolID)
		if tool == nil {
			continue
		}

		t.enabledTools[toolID] = tool
	}
}

func (t *Toolkit) ListTools() []*Tool {
	tools := make([]*Tool, 0, len(t.enabledTools))

	for _, tool := range t.enabledTools {
		tools = append(tools, tool)
	}

	return tools
}

func (t *Toolkit) GetTool(id string) (*Tool, bool) {
	tool, ok := t.enabledTools[id]
	return tool, ok
}
