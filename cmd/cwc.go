package cmd

import (
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"

	"github.com/emilkje/cwc/pkg/chat"
	"github.com/emilkje/cwc/pkg/config"
	"github.com/emilkje/cwc/pkg/errors"
	"github.com/emilkje/cwc/pkg/filetree"
	"github.com/emilkje/cwc/pkg/pathmatcher"
	"github.com/emilkje/cwc/pkg/ui"
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

The command can also receive context from standard input, useful for piping the output from another command as input.

Examples:

Including all '.go' files while excluding the 'vendor/' directory:
> cwc --include='.*.go$' --exclude='vendor/'

Including 'main.go' files from a specific path:
> cwc --include='main.go' --paths='./cmd'

Using the output of another command:
> git diff | cwc "Short commit message for these changes"`
)

func CreateRootCommand() *cobra.Command {
	var (
		includeFlag              string
		excludeFlag              string
		pathsFlag                []string
		excludeFromGitignoreFlag bool
		excludeGitDirFlag        bool
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
				// stdin is not a terminal, typically piped from another command
				if len(args) == 0 {
					return &errors.NoPromptProvidedError{Message: "no prompt provided"}
				}

				var systemContext string

				inputBytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("error reading from stdin: %w", err)
				}
				systemContext = string(inputBytes)

				return nonInteractive(systemContext, args[0])
			}

			gatherOpts := &chatOptions{
				includeFlag:              includeFlag,
				excludeFlag:              excludeFlag,
				pathsFlag:                pathsFlag,
				excludeFromGitignoreFlag: excludeFromGitignoreFlag,
				excludeGitDirFlag:        excludeGitDirFlag,
			}

			return interactiveChat(cmd, args, gatherOpts, loginCmd)
		},
	}

	initFlags(cmd, &flags{
		includeFlag:              &includeFlag,
		excludeFlag:              &excludeFlag,
		pathsFlag:                &pathsFlag,
		excludeFromGitignoreFlag: &excludeFromGitignoreFlag,
		excludeGitDirFlag:        &excludeGitDirFlag,
	})

	cmd.AddCommand(loginCmd)
	cmd.AddCommand(logoutCmd)

	return cmd
}

type flags struct {
	includeFlag              *string
	excludeFlag              *string
	pathsFlag                *[]string
	excludeFromGitignoreFlag *bool
	excludeGitDirFlag        *bool
}

func initFlags(cmd *cobra.Command, flags *flags) {
	cmd.Flags().StringVarP(flags.includeFlag, "include", "i", ".*", "a regular expression to match files to include")
	cmd.Flags().StringVarP(flags.excludeFlag, "exclude", "x", "", "a regular expression to match files to exclude")
	cmd.Flags().StringSliceVarP(flags.pathsFlag, "paths", "p", []string{"."}, "a list of paths to search for files")
	cmd.Flags().BoolVarP(flags.excludeFromGitignoreFlag,
		"exclude-from-gitignore", "e", true, "exclude files from .gitignore")
	cmd.Flags().BoolVarP(flags.excludeGitDirFlag, "exclude-git-dir", "g", true, "exclude the .git directory")

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
}

func isPiped(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

func interactiveChat(c *cobra.Command, args []string, gatherOpts *chatOptions, loginCmd *cobra.Command) error {
	cfg, err := config.NewFromConfigFile()
	if err != nil {
		var validationErr errors.ConfigValidationError
		if !stderrors.As(err, &validationErr) {
			return fmt.Errorf("error reading config: %w", err)
		}

		for _, e := range validationErr.Errors {
			ui.PrintMessage(fmt.Sprintf("error: %s\n", e), ui.MessageTypeError)
		}

		// prompt the user to sign in to refresh the config
		if !ui.AskYesNo("do you want to login now?", true) {
			ui.PrintMessage("see ya later!", ui.MessageTypeInfo)
			return nil
		}

		// login
		err = loginCmd.RunE(c, args)
		if err != nil {
			return fmt.Errorf("error logging in: %w", err)
		}
	}

	client := openai.NewClientWithConfig(cfg)

	files, rootNode, err := gatherContext(gatherOpts)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		ui.PrintMessage("No files found matching the given criteria.\n", ui.MessageTypeWarning)

		if !ui.AskYesNo("Do you wish to proceed?", false) {
			ui.PrintMessage("See ya later!", ui.MessageTypeInfo)
			return nil
		}

		return nil
	}

	fileTree := filetree.GenerateFileTree(rootNode, "", true)
	// confirm with the user that the files are correct
	ui.PrintMessage("The following files will be used as context:\n", ui.MessageTypeInfo)
	ui.PrintMessage(fileTree, ui.MessageTypeInfo)

	// warn the user of files larger than 100kb
	for _, file := range files {
		if len(file.Data) > warnFileSizeThreshold {
			largeFileMsg := fmt.Sprintf(
				"warning: %s is very large (%d bytes) and will degrade performance.\n",
				file.Path, len(file.Data))

			ui.PrintMessage(largeFileMsg, ui.MessageTypeWarning)
		}
	}

	// confirm with the user that the files are correct
	if !ui.AskYesNo("Do you wish to proceed?", true) {
		ui.PrintMessage("See ya later!", ui.MessageTypeInfo)
		return nil
	}

	contextStr := "File tree:\n\n"
	contextStr += "```\n" + fileTree + "```\n\n"
	contextStr += "File contents:\n\n"

	for _, file := range files {
		// find extension by splitting on ".". if no extension, use
		contextStr += fmt.Sprintf("./%s\n```%s\n%s\n```\n\n", file.Path, file.Type, file.Data)
	}

	systemMessage := createSystemMessageFromContext(contextStr)

	ui.PrintMessage("Type '/exit' to end the chat.\n", ui.MessageTypeNotice)

	var initialUserMessage string
	if len(args) > 0 {
		initialUserMessage = args[0]
		ui.PrintMessage(fmt.Sprintf("👤: %s\n", initialUserMessage), ui.MessageTypeInfo)
	} else {
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)
		initialUserMessage = ui.ReadUserInput()
	}

	if initialUserMessage == "/exit" {
		return nil
	}

	chatInstance := chat.NewChat(client, systemMessage, printMessageChunk)
	conversation := chatInstance.BeginConversation(initialUserMessage)

	for {
		conversation.WaitMyTurn()
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)

		userMessage := ui.ReadUserInput()

		if userMessage == "/exit" {
			break
		}

		conversation.Reply(userMessage)
	}

	return nil
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

func nonInteractive(systemMessage string, prompt string) error {
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

func createSystemMessageFromContext(context string) string {
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
		gitDirMatcher, err := pathmatcher.NewRegexPathMatcher(`^.*\.git$`)
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
