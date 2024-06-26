package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/intility/cwc/internal"
	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/prompting"
	"github.com/intility/cwc/pkg/systemcontext"
	"github.com/intility/cwc/pkg/templates"
	cwcui "github.com/intility/cwc/pkg/ui"
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
	chatOpts := internal.InteractiveChatOptions{
		IncludePattern:    "",
		ExcludePattern:    "",
		Paths:             []string{},
		TemplateName:      "",
		TemplateVariables: nil,
	}

	loginCmd := createLoginCmd()
	logoutCmd := createLogoutCmd()

	rootCmd := &cobra.Command{
		Use:   "cwc [prompt]",
		Short: "starts a new chat session",
		Long:  longDescription,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cfgProvider, err := getPlatformSpecificConfigProvider()
			if err != nil {
				return fmt.Errorf("error getting config provider: %w", err)
			}

			if isPiped(os.Stdin) {
				nic := createNonInteractiveCommand(cfgProvider, args, chatOpts.TemplateName, chatOpts.TemplateVariables)

				err = nic.Run()
				if err != nil {
					return fmt.Errorf("error running non-interactive command: %w", err)
				}

				return nil
			}

			interactiveCmd := createInteractiveCommand(args, chatOpts, cfgProvider)

			err = interactiveCmd.Run()
			if err != nil {
				return fmt.Errorf("error running interactive command: %w", err)
			}

			return nil
		},
	}

	initFlags(rootCmd, &chatOpts)

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(createTemplatesCmd())
	rootCmd.AddCommand(createConfigCommand())

	return rootCmd
}

func createNonInteractiveCommand(
	cfgProvider config.Provider,
	args []string,
	templateName string,
	templateVars map[string]string,
) *internal.NonInteractiveCmd {
	clientProvider := config.NewOpenAIClientProvider(cfgProvider)
	templateLocator := getTemplateLocator(cfgProvider)
	promptResolver := prompting.NewArgsOrTemplatePromptResolver(templateLocator, args, templateName)

	contextRetriever := systemcontext.NewIOReaderContextRetriever(os.Stdin)
	smGenerator := systemcontext.NewTemplatedSystemMessageGenerator(
		templateLocator,
		templateName,
		templateVars,
		contextRetriever,
	)

	return internal.NewNonInteractiveCmd(
		clientProvider,
		promptResolver,
		smGenerator,
	)
}

func createInteractiveCommand(
	args []string,
	opts internal.InteractiveChatOptions,
	cfgProvider config.Provider,
) *internal.InteractiveCmd {
	clientProvider := config.NewOpenAIClientProvider(cfgProvider)
	templateLocator := getTemplateLocator(cfgProvider)
	promptResolver := prompting.NewArgsOrTemplatePromptResolver(templateLocator, args, opts.TemplateName)

	retrieverConfig := systemcontext.FileContextRetrieverOptions{
		CfgProvider:    cfgProvider,
		IncludePattern: opts.IncludePattern,
		ExcludePattern: opts.ExcludePattern,
		SearchScopes:   opts.Paths,
		ContextPrinter: printContext,
	}

	contextRetriever := systemcontext.NewFileContextRetriever(retrieverConfig)

	smGenerator := systemcontext.NewTemplatedSystemMessageGenerator(
		templateLocator,
		opts.TemplateName,
		opts.TemplateVariables,
		contextRetriever,
	)

	return internal.NewInteractiveCmd(
		promptResolver,
		clientProvider,
		smGenerator,
		opts,
	)
}

func getPlatformSpecificConfigProvider() (config.Provider, error) { //nolint: ireturn
	var cfgProvider config.Provider

	if config.IsWSL() {
		configDir, err := config.GetConfigDir()
		if err != nil {
			return nil, fmt.Errorf("error getting config directory: %w", err)
		}

		keyStore := config.NewAPIKeyFileStore(filepath.Join(configDir, "api.key"))
		cfgProvider = config.NewDefaultProvider(config.WithKeyStore(keyStore))
	} else {
		cfgProvider = config.NewDefaultProvider()
	}

	return cfgProvider, nil
}

func printContext(fileTree string, files []filetree.File) {
	ui := cwcui.NewUI()
	ui.PrintMessage(fileTree, cwcui.MessageTypeInfo)

	for _, file := range files {
		printLargeFileWarning(file)
	}
}

func printLargeFileWarning(file filetree.File) {
	if len(file.Data) > warnFileSizeThreshold {
		largeFileMsg := fmt.Sprintf(
			"warning: %s is very large (%d bytes) and will degrade performance.\n",
			file.Path, len(file.Data))

		ui := cwcui.NewUI()
		ui.PrintMessage(largeFileMsg, cwcui.MessageTypeWarning)
	}
}

func getTemplateLocator(cfgProvider config.Provider) *templates.MergedTemplateLocator {
	var locators []templates.TemplateLocator

	configDir, err := cfgProvider.GetConfigDir()
	if err == nil {
		locators = append(locators, templates.NewYamlFileTemplateLocator(filepath.Join(configDir, "templates.yaml")))
	}

	locators = append(locators, templates.NewYamlFileTemplateLocator(filepath.Join(".cwc", "templates.yaml")))
	mergedLocator := templates.NewMergedTemplateLocator(locators...)

	return mergedLocator
}

func initFlags(cmd *cobra.Command, opts *internal.InteractiveChatOptions) {
	cmd.Flags().StringVarP(&opts.IncludePattern, "include", "i", ".*", "a regular expression to match files to include")
	cmd.Flags().StringVarP(&opts.ExcludePattern, "exclude", "x", "", "a regular expression to match files to exclude")
	cmd.Flags().StringSliceVarP(&opts.Paths, "paths", "p", []string{"."}, "a list of paths to search for files")
	cmd.Flags().StringVarP(&opts.TemplateName, "template", "t", "default", "the name of the template to use")
	cmd.Flags().StringToStringVarP(&opts.TemplateVariables,
		"template-variables", "v", nil, "variables to use in the template")

	cmd.Flag("include").
		Usage = "Specify a regex pattern to include files. " +
		"For example, to include only Markdown files, use --include '\\.md$'"
	cmd.Flag("exclude").
		Usage = "Specify a regex pattern to exclude files. For example, to exclude test files, use --exclude '_test\\\\.go$'"
	cmd.Flag("paths").
		Usage = "Specify a list of paths to search for files. For example, " +
		"to search in the 'cmd' and 'pkg' directories, use --paths cmd,pkg"
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
