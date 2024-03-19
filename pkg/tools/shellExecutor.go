package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"text/template"
)

type ShellExecutor struct{}

func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

func (s *ShellExecutor) Execute(tool Tool, arguments string) (string, error) {
	results := make([]string, 0)

	args := make(map[string]string)

	err := json.Unmarshal([]byte(arguments), &args)
	if err != nil {
		return "", fmt.Errorf("error parsing arguments: %w", err)
	}

	// create templates from the tool's scripts
	for _, script := range tool.ShellExecutables() {
		tmpl := template.New("script")

		tmpl, err = tmpl.Parse(script)
		if err != nil {
			return "", fmt.Errorf("error parsing shell script: %w", err)
		}

		// execute the template with the arguments
		var renderedScript strings.Builder

		err = tmpl.Execute(&renderedScript, args)
		if err != nil {
			return "", fmt.Errorf("error rendering shell script: %w", err)
		}

		// execute the rendered script
		scriptToExecute := renderedScript.String()
		cmd := exec.Command("sh", "-c", scriptToExecute)

		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("error executing shell script: %w", err)
		}

		results = append(results, string(out))
	}

	// execute the tool with the arguments
	// return the result
	return strings.Join(results, "\n"), nil
}
