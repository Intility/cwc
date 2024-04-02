package cmd

import (
	stdErrors "errors"
	"fmt"
	"path/filepath"
	"strings"
	tt "text/template"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/pathmatcher"
	"github.com/intility/cwc/pkg/templates"
	cwcui "github.com/intility/cwc/pkg/ui"
)

const (
	warnFileSizeThreshold = 100000
)

func newClientFromConfig() (*openai.Client, error) {
	cfg, err := config.NewFromConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	client := openai.NewClientWithConfig(cfg)

	return client, nil
}

func determinePrompt(args []string, templateName string) string {
	var prompt string

	template, err := getTemplate(templateName)
	if err != nil {
		var templateNotFoundError errors.TemplateNotFoundError
		if stdErrors.As(err, &templateNotFoundError) {
			if len(args) == 0 {
				return ""
			}
		}
	} else {
		prompt = template.DefaultPrompt
	}

	// args takes precedence over template.DefaultPrompt
	if len(args) > 0 {
		prompt = args[0]
	}

	return prompt
}

func getTemplate(templateName string) (*templates.Template, error) {
	if templateName == "" {
		templateName = "default"
	}

	var locators []templates.TemplateLocator

	configDir, err := config.GetConfigDir()
	if err == nil {
		locators = append(locators, templates.NewYamlFileTemplateLocator(filepath.Join(configDir, "templates.yaml")))
	}

	locators = append(locators, templates.NewYamlFileTemplateLocator(filepath.Join(".cwc", "templates.yaml")))
	mergedLocator := templates.NewMergedTemplateLocator(locators...)

	tmpl, err := mergedLocator.GetTemplate(templateName)
	if err != nil {
		if errors.IsTemplateNotFoundError(err) {
			return nil, errors.TemplateNotFoundError{TemplateName: templateName}
		}

		return nil, fmt.Errorf("error getting template: %w", err)
	}

	return tmpl, nil
}

func createSystemMessage(ctx string, templateName string, templateVariables map[string]string) (string, error) {
	template, err := getTemplate(templateName)

	if templateVariables == nil {
		templateVariables = make(map[string]string)
	}

	// if no template found, create a basic template as fallback
	if errors.IsTemplateNotFoundError(err) {
		return createBuiltinSystemMessageFromContext(ctx), nil
	}

	// compile the template.SystemMessage as a go template
	tmpl, err := tt.New("systemMessage").Parse(template.SystemMessage)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	type valueBag struct {
		Context   string
		Variables map[string]string
	}

	// populate the variables map with default values if not provided
	for _, v := range template.Variables {
		if _, ok := templateVariables[v.Name]; !ok {
			templateVariables[v.Name] = v.DefaultValue
		}
	}

	values := valueBag{
		Context:   ctx,
		Variables: templateVariables,
	}

	writer := &strings.Builder{}
	err = tmpl.Execute(writer, values)

	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return writer.String(), nil
}

func createBuiltinSystemMessageFromContext(context string) string {
	var systemMessage strings.Builder

	systemMessage.WriteString("You are a helpful coding assistant. ")
	systemMessage.WriteString("Below you will find relevant context to answer the user's question.\n\n")
	systemMessage.WriteString("Context:\n")
	systemMessage.WriteString(context)
	systemMessage.WriteString("\n\n")
	systemMessage.WriteString("Please follow the users instructions, you can do this!")

	return systemMessage.String()
}

// askConfirmation prompts the user if they want to proceed with no files.
func askConfirmation(prompt string, messageType cwcui.MessageType) bool {
	ui := cwcui.NewUI()
	ui.PrintMessage(prompt, messageType)

	if !ui.AskYesNo("Do you wish to proceed?", false) {
		ui.PrintMessage("See ya later!", cwcui.MessageTypeInfo)
		return false
	}

	return true
}

func excludeMatchersFromConfig() ([]pathmatcher.PathMatcher, error) {
	var excludeMatchers []pathmatcher.PathMatcher

	ui := cwcui.NewUI() //nolint:varnamelen

	cfg, err := config.LoadConfig()
	if err != nil {
		return excludeMatchers, fmt.Errorf("error loading config: %w", err)
	}

	if cfg.UseGitignore {
		gitignoreMatcher, err := pathmatcher.NewGitignorePathMatcher()
		if err != nil {
			if errors.IsGitNotInstalledError(err) {
				ui.PrintMessage("warning: git not found in PATH, skipping .gitignore\n", cwcui.MessageTypeWarning)
			} else {
				return nil, fmt.Errorf("error creating gitignore matcher: %w", err)
			}
		}

		excludeMatchers = append(excludeMatchers, gitignoreMatcher)
	}

	if cfg.ExcludeGitDir {
		gitDirMatcher, err := pathmatcher.NewRegexPathMatcher(`^\.git(/|\\)`)
		if err != nil {
			return nil, fmt.Errorf("error creating git directory matcher: %w", err)
		}

		excludeMatchers = append(excludeMatchers, gitDirMatcher)
	}

	return excludeMatchers, nil
}

func printLargeFileWarning(file filetree.File) {
	ui := cwcui.NewUI() //nolint:varnamelen
	if len(file.Data) > warnFileSizeThreshold {
		largeFileMsg := fmt.Sprintf(
			"warning: %s is very large (%d bytes) and will degrade performance.\n",
			file.Path, len(file.Data))

		ui.PrintMessage(largeFileMsg, cwcui.MessageTypeWarning)
	}
}

func createContext(fileTree string, files []filetree.File) string {
	contextStr := "File tree:\n\n"
	contextStr += "```\n" + fileTree + "```\n\n"
	contextStr += "File contents:\n\n"

	for _, file := range files {
		// find extension by splitting on ".". if no extension, use
		contextStr += fmt.Sprintf("./%s\n```%s\n%s\n```\n\n", file.Path, file.Type, file.Data)
	}

	return contextStr
}
