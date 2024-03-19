package cmd

import (
	stdErrors "errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	tt "text/template"

	"github.com/intility/cwc/pkg/chat"
	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/pathmatcher"
	"github.com/intility/cwc/pkg/templates"
	"github.com/intility/cwc/pkg/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

const (
	warnFileSizeThreshold = 100000
	longDescription       = `The 'cwc' command initiates a new chat session, 
providing granular control over the inclusion and exclusion of files via regular expression patterns. 
It allows for specification of paths to include or exclude files from the chat context.

Features at a glance:

- Regex-based file inclusion and exclusion patterns
- .gitignore integration for ignoring files
- Option to specify directories for inclusion scope
- Interactive file selection and confirmation
- Reading from standard input for a non-interactive session
- Use of templates for system messages and default prompts

The command can also receive context from standard input, useful for piping the output from another command as input.

Examples:

Including all '.go' files while excluding the 'vendor/' directory:
> cwc --include='.*.go$' --exclude='vendor/'

Including 'main.go' files from a specific path:
> cwc --include='main.go' --paths='./cmd'

Using the output of another command:
> git diff | cwc "Short commit message for these changes"

Using a specific template:
> cwc --template=tech_writer --template-variables rizz=max
`
)

func CreateRootCommand() *cobra.Command {
	var (
		includeFlag              string
		excludeFlag              string
		pathsFlag                []string
		excludeFromGitignoreFlag bool
		excludeGitDirFlag        bool
		templateFlag             string
		templateVariablesFlag    map[string]string
	)

	loginCmd := createLoginCmd()
	logoutCmd := createLogoutCmd()

	cmd := &cobra.Command{
		Use:   "cwc [prompt]",
		Short: "starts a new chat session",
		Long:  longDescription,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if isPiped(os.Stdin) {
				return nonInteractive(args, templateFlag, templateVariablesFlag)
			}

			chatOpts := &chatOptions{
				includeFlag:              includeFlag,
				excludeFlag:              excludeFlag,
				pathsFlag:                pathsFlag,
				excludeFromGitignoreFlag: excludeFromGitignoreFlag,
				excludeGitDirFlag:        excludeGitDirFlag,
				templateFlag:             templateFlag,
				templateVariablesFlag:    templateVariablesFlag,
			}

			var prompt string
			if len(args) > 0 {
				prompt = args[0]
			}

			return interactiveChat(prompt, chatOpts)
		},
	}

	initFlags(cmd, &flags{
		includeFlag:              &includeFlag,
		excludeFlag:              &excludeFlag,
		pathsFlag:                &pathsFlag,
		excludeFromGitignoreFlag: &excludeFromGitignoreFlag,
		excludeGitDirFlag:        &excludeGitDirFlag,
		templateFlag:             &templateFlag,
		templateVariablesFlag:    &templateVariablesFlag,
	})

	cmd.AddCommand(loginCmd)
	cmd.AddCommand(logoutCmd)
	cmd.AddCommand(createTemplatesCmd())

	return cmd
}

type flags struct {
	includeFlag              *string
	excludeFlag              *string
	pathsFlag                *[]string
	excludeFromGitignoreFlag *bool
	excludeGitDirFlag        *bool
	templateFlag             *string
	templateVariablesFlag    *map[string]string
}

func initFlags(cmd *cobra.Command, flags *flags) {
	cmd.Flags().StringVarP(flags.includeFlag, "include", "i", ".*", "a regular expression to match files to include")
	cmd.Flags().StringVarP(flags.excludeFlag, "exclude", "x", "", "a regular expression to match files to exclude")
	cmd.Flags().StringSliceVarP(flags.pathsFlag, "paths", "p", []string{"."}, "a list of paths to search for files")
	cmd.Flags().BoolVarP(flags.excludeFromGitignoreFlag,
		"exclude-from-gitignore", "e", true, "exclude files from .gitignore")
	cmd.Flags().BoolVarP(flags.excludeGitDirFlag, "exclude-git-dir", "g", true, "exclude the .git directory")
	cmd.Flags().StringVarP(flags.templateFlag, "template", "t", "default", "the name of the template to use")
	cmd.Flags().StringToStringVarP(flags.templateVariablesFlag,
		"template-variables", "v", nil, "variables to use in the template")

	cmd.Flag("include").
		Usage = "Specify a regex pattern to include files. " +
		"For example, to include only Markdown files, use --include '\\.md$'"
	cmd.Flag("exclude").
		Usage = "Specify a regex pattern to exclude files. For example, to exclude test files, use --exclude '_test\\\\.go$'"
	cmd.Flag("paths").
		Usage = "Specify a list of paths to search for files. For example, " +
		"to search in the 'cmd' and 'pkg' directories, use --paths cmd,pkg"
	cmd.Flag("exclude-from-gitignore").
		Usage = "Exclude files from .gitignore. If set to false, files mentioned in .gitignore will not be excluded"
	cmd.Flag("exclude-git-dir").
		Usage = "Exclude the .git directory. If set to false, the .git directory will not be excluded"
	cmd.Flag("template").
		Usage = "Specify the name of the template to use. For example, " +
		"to use a template named 'tech_writer', use --template tech_writer"
	cmd.Flag("template-variables").
		Usage = "Specify variables to use in the template. For example, to use the variable 'name' " +
		"with the value 'John', use --template-variables name=John"
}

func isPiped(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

func interactiveChat(prompt string, chatOpts *chatOptions) error {
	// Load configuration
	cfg, err := config.NewFromConfigFile()
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	client := openai.NewClientWithConfig(cfg)

	// Gather context from files
	files, fileTree, err := gatherAndPrintContext(chatOpts)
	if err != nil {
		return err
	} else if len(files) == 0 { // No files found, terminating or confirming to proceed
		if !askConfirmation("No files found matching the given criteria.\n", ui.MessageTypeWarning) {
			return nil
		}
	}

	contextStr := createContext(fileTree, files)

	systemMessage, err := createSystemMessage(contextStr, chatOpts.templateFlag, chatOpts.templateVariablesFlag)
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	ui.PrintMessage("Type '/exit' to end the chat.\n", ui.MessageTypeNotice)

	if prompt == "" {
		prompt = getPromptFromUserOrTemplate(chatOpts.templateFlag)
	} else {
		ui.PrintMessage(fmt.Sprintf("👤: %s\n", prompt), ui.MessageTypeInfo)
	}

	if prompt == "/exit" {
		return nil
	}

	handleChat(client, systemMessage, prompt)

	return nil
}

func getPromptFromUserOrTemplate(templateName string) string {
	// get default prompt from template
	var prompt string

	if templateName == "" {
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)
		return ui.ReadUserInput()
	}

	template, err := getTemplate(templateName)
	if err != nil {
		ui.PrintMessage(err.Error()+"\n", ui.MessageTypeWarning)
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)

		return ui.ReadUserInput()
	}

	if template.DefaultPrompt == "" {
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)
		prompt = ui.ReadUserInput()
	} else {
		prompt = template.DefaultPrompt
		ui.PrintMessage(fmt.Sprintf("👤: %s\n", prompt), ui.MessageTypeInfo)
	}

	return prompt
}

func handleChat(client *openai.Client, systemMessage string, prompt string) {
	chatInstance := chat.NewChat(client, systemMessage, printMessageChunk)
	conversation := chatInstance.BeginConversation(prompt)

	for {
		conversation.WaitMyTurn()
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)

		userMessage := ui.ReadUserInput()

		if userMessage == "/exit" {
			break
		}

		conversation.Reply(userMessage)
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

// gatherAndPrintContext gathers file context based on provided options and prints it out.
func gatherAndPrintContext(chatOptions *chatOptions) ([]filetree.File, string, error) {
	files, rootNode, err := gatherContext(chatOptions)
	if err != nil {
		return nil, "", err
	}

	for _, file := range files {
		printLargeFileWarning(file)
	}

	fileTree := filetree.GenerateFileTree(rootNode, "", true)

	ui.PrintMessage("The following files will be used as context:\n", ui.MessageTypeInfo)
	ui.PrintMessage(fileTree, ui.MessageTypeInfo)

	return files, fileTree, nil
}

// askConfirmation prompts the user if they want to proceed with no files.
func askConfirmation(prompt string, messageType ui.MessageType) bool {
	ui.PrintMessage(prompt, messageType)

	if !ui.AskYesNo("Do you wish to proceed?", false) {
		ui.PrintMessage("See ya later!", ui.MessageTypeInfo)
		return false
	}

	return true
}

func printLargeFileWarning(file filetree.File) {
	if len(file.Data) > warnFileSizeThreshold {
		largeFileMsg := fmt.Sprintf(
			"warning: %s is very large (%d bytes) and will degrade performance.\n",
			file.Path, len(file.Data))

		ui.PrintMessage(largeFileMsg, ui.MessageTypeWarning)
	}
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
	var templateNotFoundError errors.TemplateNotFoundError
	if err != nil && stdErrors.As(err, &templateNotFoundError) {
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

func printMessageChunk(chunk *chat.ConversationChunk) {
	if chunk.IsInitialChunk {
		ui.PrintMessage("🤖: ", ui.MessageTypeInfo)
		return
	}

	if chunk.IsErrorChunk {
		ui.PrintMessage(chunk.Content, ui.MessageTypeError)
	}

	if chunk.IsFinalChunk {
		ui.PrintMessage("\n", ui.MessageTypeInfo)
	}

	ui.PrintMessage(chunk.Content, ui.MessageTypeInfo)
}

func nonInteractive(args []string, templateName string, templateVars map[string]string) error {
	var prompt string

	template, err := getTemplate(templateName)
	if err != nil {
		// if no template found, create a basic template as fallback
		var templateNotFoundError errors.TemplateNotFoundError
		if stdErrors.As(err, &templateNotFoundError) {
			if len(args) == 0 {
				return &errors.NoPromptProvidedError{Message: "no prompt provided"}
			}
		}
	}

	prompt = template.DefaultPrompt

	// args takes precedence over template.DefaultPrompt
	if len(args) > 0 {
		prompt = args[0]
	}

	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}

	systemContext := string(inputBytes)

	systemMessage, err := createSystemMessage(systemContext, templateName, templateVars)
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	cfg, err := config.NewFromConfigFile()
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	client := openai.NewClientWithConfig(cfg)

	onChunk := func(chunk *chat.ConversationChunk) {
		ui.PrintMessage(chunk.Content, ui.MessageTypeInfo)
	}
	chatInstance := chat.NewChat(client, systemMessage, onChunk)
	conversation := chatInstance.BeginConversation(prompt)

	conversation.WaitMyTurn()

	return nil
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

type chatOptions struct {
	includeFlag              string
	excludeFlag              string
	pathsFlag                []string
	excludeFromGitignoreFlag bool
	excludeGitDirFlag        bool
	templateFlag             string
	templateVariablesFlag    map[string]string
}

func gatherContext(opts *chatOptions) ([]filetree.File, *filetree.FileNode, error) {
	includeFlag := opts.includeFlag
	excludeFlag := opts.excludeFlag
	pathsFlag := opts.pathsFlag
	excludeFromGitignoreFlag := opts.excludeFromGitignoreFlag
	excludeGitDirFlag := opts.excludeGitDirFlag

	var excludeMatchers []pathmatcher.PathMatcher

	// add exclude flag to excludeMatchers
	if excludeFlag != "" {
		excludeMatcher, err := pathmatcher.NewRegexPathMatcher(excludeFlag)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating exclude matcher: %w", err)
		}

		excludeMatchers = append(excludeMatchers, excludeMatcher)
	}

	if excludeFromGitignoreFlag {
		gitignoreMatcher, err := pathmatcher.NewGitignorePathMatcher()
		if err != nil {
			if errors.IsGitNotInstalledError(err) {
				ui.PrintMessage("warning: git not found in PATH, skipping .gitignore\n", ui.MessageTypeWarning)
			} else {
				return nil, nil, fmt.Errorf("error creating gitignore matcher: %w", err)
			}
		}

		excludeMatchers = append(excludeMatchers, gitignoreMatcher)
	}

	if excludeGitDirFlag {
		gitDirMatcher, err := pathmatcher.NewRegexPathMatcher(`^\.git(/|\\)`)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating git directory matcher: %w", err)
		}

		excludeMatchers = append(excludeMatchers, gitDirMatcher)
	}

	excludeMatcher := pathmatcher.NewCompoundPathMatcher(excludeMatchers...)

	// includeMatcher
	includeMatcher, err := pathmatcher.NewRegexPathMatcher(includeFlag)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating include matcher: %w", err)
	}

	files, rootNode, err := filetree.GatherFiles(&filetree.FileGatherOptions{
		IncludeMatcher: includeMatcher,
		ExcludeMatcher: excludeMatcher,
		PathScopes:     pathsFlag,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error gathering files: %w", err)
	}

	return files, rootNode, nil
}
