package cmd

import (
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/chat"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/pathmatcher"
	cwcui "github.com/intility/cwc/pkg/ui"
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
	ui          cwcui.UI
}

func NewInteractiveCmd(args []string, chatOptions InteractiveChatOptions) *InteractiveCmd {
	prompt := determinePrompt(args, chatOptions.TemplateName)
	ui := cwcui.NewUI()

	return &InteractiveCmd{prompt: prompt, chatOptions: chatOptions, ui: ui}
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
		if !askConfirmation("No files found matching the given criteria.\n", cwcui.MessageTypeWarning) {
			return nil
		}
	}

	contextStr := createContext(fileTree, files)

	systemMessage, err := createSystemMessage(contextStr, c.chatOptions.TemplateName, c.chatOptions.TemplateVariables)
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	c.ui.PrintMessage("Type '/exit' to end the chat.\n", cwcui.MessageTypeNotice)

	if c.prompt == "" {
		c.ui.PrintMessage("👤: ", cwcui.MessageTypeInfo)
		c.prompt = c.ui.ReadUserInput()
	} else {
		c.ui.PrintMessage(fmt.Sprintf("👤: %s\n", c.prompt), cwcui.MessageTypeInfo)
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
		c.ui.PrintMessage("👤: ", cwcui.MessageTypeInfo)

		userMessage := c.ui.ReadUserInput()

		if userMessage == "/exit" {
			break
		}

		conversation.Reply(userMessage)
	}
}

func (c *InteractiveCmd) printMessageChunk(chunk *chat.ConversationChunk) {
	if chunk.IsInitialChunk {
		c.ui.PrintMessage("🤖: ", cwcui.MessageTypeInfo)
		return
	}

	if chunk.IsErrorChunk {
		c.ui.PrintMessage(chunk.Content, cwcui.MessageTypeError)
	}

	if chunk.IsFinalChunk {
		c.ui.PrintMessage("\n", cwcui.MessageTypeInfo)
	}

	c.ui.PrintMessage(chunk.Content, cwcui.MessageTypeInfo)
}

func (c *InteractiveCmd) gatherAndPrintContext() ([]filetree.File, string, error) {
	// gatherAndPrintContext gathers file context based on provided options and prints it out.
	ui := cwcui.NewUI() //nolint:varnamelen

	files, rootNode, err := c.gatherContext()
	if err != nil {
		return nil, "", err
	}

	for _, file := range files {
		printLargeFileWarning(file)
	}

	fileTree := filetree.GenerateFileTree(rootNode, "", true)

	ui.PrintMessage("The following files will be used as context:\n", cwcui.MessageTypeInfo)
	ui.PrintMessage(fileTree, cwcui.MessageTypeInfo)

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
