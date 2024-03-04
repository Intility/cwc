package cmd

import (
	"context"
	stderrors "errors"
	"fmt"
	"github.com/emilkje/cwc/pkg/pathmatcher"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"io"

	"github.com/emilkje/cwc/pkg/config"
	"github.com/emilkje/cwc/pkg/errors"
	"github.com/emilkje/cwc/pkg/filetree"
	"github.com/emilkje/cwc/pkg/ui"
)

var (
	includeFlag              string
	excludeFlag              string
	pathsFlag                []string
	excludeFromGitignoreFlag bool
	excludeGitDirFlag        bool
)

var CwcCmd = &cobra.Command{
	Use:   "cwc [prompt]",
	Short: "starts a new chat session",
	Long: `Starts a new chat session with detailed control over what files are included or excluded using regular expressions.

The --include flag allows you to specify a regular expression to match the files that should be included in the chat context. Only files that match the given pattern will be considered.

The --exclude flag allows you to specify a regular expression to match the files that should be excluded from the chat context. If a file matches the given pattern, it will be ignored even if it matches the include pattern.

If the --exclude-from-gitignore is provided and set to true (which is the default), the files mentioned in the .gitignore file will be excluded automatically.

For example, if you want to include all '.go' files but exclude files in 'vendor/' directory, you would start the chat with:
> cwc --include='\.go$' --exclude='vendor/'

If you have multiple files called 'main.go' for example, you can use the --paths qualifier to specify which files to include. For example:
> cwc --include='\.go$' --paths='./cmd'`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(c *cobra.Command, args []string) error {

		cfg, err := config.NewFromConfigFile()
		if err != nil {

			// check of validation error
			if validationErr, ok := errors.AsConfigValidationError(err); ok {
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
					return err
				}
			} else {
				return err
			}
		}

		client := openai.NewClientWithConfig(cfg)

		var excludeMatchers []pathmatcher.PathMatcher

		// add exclude flag to excludeMatchers
		if excludeFlag != "" {
			excludeMatcher, err := pathmatcher.NewRegexPathMatcher(excludeFlag)
			if err != nil {
				return err
			}
			excludeMatchers = append(excludeMatchers, excludeMatcher)
		}

		if excludeFromGitignoreFlag {
			gitignoreMatcher, err := pathmatcher.NewGitignorePathMatcher(".gitignore")

			if err != nil {
				if !errors.IsFileNotExistError(err) {
					return err
				}

				ui.PrintMessage(".gitignore does not exist, skipping.\n", ui.MessageTypeInfo)
			} else {
				excludeMatchers = append(excludeMatchers, gitignoreMatcher)
			}
		}

		if excludeGitDirFlag {
			// TODO: fix this so that .github is not excluded
			gitDirMatcher, err := pathmatcher.NewRegexPathMatcher(`^.*\.git`)
			if err != nil {
				return err
			}
			excludeMatchers = append(excludeMatchers, gitDirMatcher)
		}

		excludeMatcher := pathmatcher.NewCompoundPathMatcher(excludeMatchers...)

		// includeMatcher
		includeMatcher, err := pathmatcher.NewRegexPathMatcher(includeFlag)
		if err != nil {
			return err
		}

		files, rootNode, err := filetree.GatherFiles(includeMatcher, excludeMatcher, pathsFlag)

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

		// confirm with the user that the files are correct
		ui.PrintMessage("The following files will be used as context:\n", ui.MessageTypeInfo)
		fileTree := filetree.GenerateFileTree(rootNode, "", true)
		ui.PrintMessage(fileTree, ui.MessageTypeInfo)

		// warn the user of files larger than 100kb
		for _, file := range files {
			if len(file.Data) > 100000 {
				ui.PrintMessage(fmt.Sprintf("warning: %s is very large (%d bytes) and will degrade performance.\n", file.Path, len(file.Data)), ui.MessageTypeWarning)
			}
		}

		if !ui.AskYesNo("Do you wish to proceed?", true) {
			ui.PrintMessage("See ya later!", ui.MessageTypeInfo)
			return nil
		}

		contextStr := "Context:\n\n"
		contextStr += "## File tree\n\n"
		contextStr += "```\n" + fileTree + "```\n\n"
		contextStr += "## File contents\n\n"
		for _, file := range files {
			// find extension by splitting on ".". if no extension, use
			contextStr += fmt.Sprintf("./%s\n```%s\n%s\n```\n\n", file.Path, file.Type, file.Data)
		}

		systemMessage := "You are a helpful coding assistant. Below you will find relevant context to answer the user's question.\n\n" +
			contextStr + "\n\n" +
			"Please follow the users instructions, you can do this!"

		ui.PrintMessage("Type '/exit' to end the chat.\n", ui.MessageTypeNotice)

		initialUserMessage := ""
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

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: initialUserMessage,
			},
		}

		for {
			req := openai.ChatCompletionRequest{
				Model: openai.GPT4TurboPreview,
				//MaxTokens: 4096,
				Messages: messages,
				Stream:   true,
			}

			ctx := context.Background()
			stream, err := client.CreateChatCompletionStream(ctx, req)
			if err != nil {
				return err
			}

			messageStr := ""

			ui.PrintMessage("🤖: ", ui.MessageTypeInfo)
		answer:
			for {
				response, err := stream.Recv()
				if stderrors.Is(err, io.EOF) {
					break answer
				}

				if err != nil {
					return err
				}

				if len(response.Choices) == 0 {
					continue answer
				}

				messageStr = response.Choices[0].Delta.Content
				ui.PrintMessage(response.Choices[0].Delta.Content, ui.MessageTypeInfo)
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: messageStr,
			})

			// read user input until newline
			ui.PrintMessage("\n👤: ", ui.MessageTypeInfo)
			userInput := ui.ReadUserInput()

			// check for slash commands
			if userInput == "/exit" {
				break
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			})

			// close the stream for the current request
			stream.Close()
		}
		return nil
	},
}

func init() {

	CwcCmd.Flags().StringVarP(&includeFlag, "include", "i", ".*", "a regular expression to match files to include")
	CwcCmd.Flags().StringVarP(&excludeFlag, "exclude", "x", "", "a regular expression to match files to exclude")
	CwcCmd.Flags().StringSliceVarP(&pathsFlag, "paths", "p", []string{"."}, "a list of paths to search for files")
	CwcCmd.Flags().BoolVarP(&excludeFromGitignoreFlag, "exclude-from-gitignore", "e", true, "exclude files from .gitignore")
	CwcCmd.Flags().BoolVarP(&excludeGitDirFlag, "exclude-git-dir", "g", true, "exclude the .git directory")

	CwcCmd.Flag("include").Usage = "Specify a regex pattern to include files. For example, to include only Markdown files, use --include '\\.md$'"
	CwcCmd.Flag("exclude").Usage = "Specify a regex pattern to exclude files. For example, to exclude test files, use --exclude '_test\\\\.go$'"
	CwcCmd.Flag("paths").Usage = "Specify a list of paths to search for files. For example, to search in the 'cmd' and 'pkg' directories, use --paths cmd,pkg"
	CwcCmd.Flag("exclude-from-gitignore").Usage = "Exclude files from .gitignore. If set to false, files mentioned in .gitignore will not be excluded"
	CwcCmd.Flag("exclude-git-dir").Usage = "Exclude the .git directory. If set to false, the .git directory will not be excluded"

	// Add the login command to the root command so that it is available to the CLI
	CwcCmd.AddCommand(loginCmd)
	CwcCmd.AddCommand(logoutCmd)
}
