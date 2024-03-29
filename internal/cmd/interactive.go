package cmd

import (
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/chat"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/pathmatcher"
	"github.com/intility/cwc/pkg/ui"
)

type InteractiveChatOptions struct {
	IncludePattern    string
	ExcludePattern    string
	Paths             []string
	TemplateName      string
	TemplateVariables map[string]string
}

type InteractiveCmd struct {
	prompt      string
	chatOptions InteractiveChatOptions
}

func NewInteractiveCmd(args []string, chatOptions InteractiveChatOptions) *InteractiveCmd {
	prompt := determinePrompt(args, chatOptions.TemplateName)
	return &InteractiveCmd{prompt: prompt, chatOptions: chatOptions}
}

func (c *InteractiveCmd) Run() error {
	client, err := newClientFromConfig()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	files, fileTree, err := c.gatherAndPrintContext()
	if err != nil {
		return err
	} else if len(files) == 0 { // No files found, terminating or confirming to proceed
		if !askConfirmation("No files found matching the given criteria.\n", ui.MessageTypeWarning) {
			return nil
		}
	}

	contextStr := createContext(fileTree, files)

	systemMessage, err := createSystemMessage(contextStr, c.chatOptions.TemplateName, c.chatOptions.TemplateVariables)
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	ui.PrintMessage("Type '/exit' to end the chat.\n", ui.MessageTypeNotice)

	if c.prompt == "" {
		ui.PrintMessage("👤: ", ui.MessageTypeInfo)
		c.prompt = ui.ReadUserInput()
	} else {
		ui.PrintMessage(fmt.Sprintf("👤: %s\n", c.prompt), ui.MessageTypeInfo)
	}

	if c.prompt == "/exit" {
		return nil
	}

	c.handleChat(client, systemMessage, c.prompt)

	return nil
}

func (c *InteractiveCmd) handleChat(client *openai.Client, systemMessage string, prompt string) {
	chatInstance := chat.NewChat(client, systemMessage, c.printMessageChunk)
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

func (c *InteractiveCmd) printMessageChunk(chunk *chat.ConversationChunk) {
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

func (c *InteractiveCmd) gatherAndPrintContext() ([]filetree.File, string, error) {
	// gatherAndPrintContext gathers file context based on provided options and prints it out.
	files, rootNode, err := c.gatherContext()
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

func (c *InteractiveCmd) gatherContext() ([]filetree.File, *filetree.FileNode, error) {
	var excludeMatchers []pathmatcher.PathMatcher

	// add exclude flag to excludeMatchers
	if c.chatOptions.ExcludePattern != "" {
		excludeMatcher, err := pathmatcher.NewRegexPathMatcher(c.chatOptions.ExcludePattern)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating exclude matcher: %w", err)
		}

		excludeMatchers = append(excludeMatchers, excludeMatcher)
	}

	excludeMatchersFromConfig, err := excludeMatchersFromConfig()
	if err != nil {
		return nil, nil, err
	}

	excludeMatchers = append(excludeMatchers, excludeMatchersFromConfig...)

	excludeMatcher := pathmatcher.NewCompoundPathMatcher(excludeMatchers...)

	includeMatcher, err := pathmatcher.NewRegexPathMatcher(c.chatOptions.IncludePattern)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating include matcher: %w", err)
	}

	files, rootNode, err := filetree.GatherFiles(&filetree.FileGatherOptions{
		IncludeMatcher: includeMatcher,
		ExcludeMatcher: excludeMatcher,
		PathScopes:     c.chatOptions.Paths,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error gathering files: %w", err)
	}

	return files, rootNode, nil
}
