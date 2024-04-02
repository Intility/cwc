package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/intility/cwc/pkg/chat"
	cwcui "github.com/intility/cwc/pkg/ui"
)

type NonInteractiveCmd struct {
	prompt       string
	templateName string
	templateVars map[string]string
	ui           cwcui.UI
}

func NewNonInteractiveCmd(args []string, templateName string, templateVars map[string]string) *NonInteractiveCmd {
	prompt := determinePrompt(args, templateName)
	ui := cwcui.NewUI()

	return &NonInteractiveCmd{prompt: prompt, templateName: templateName, templateVars: templateVars, ui: ui}
}

func (c *NonInteractiveCmd) Run() error {
	client, err := newClientFromConfig()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	systemCtx, err := c.readContextFromStdIn()
	if err != nil {
		return fmt.Errorf("error reading context from stdin: %w", err)
	}

	systemMessage, err := createSystemMessage(systemCtx, c.templateName, c.templateVars)
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	chatInstance := chat.NewChat(client, systemMessage, c.printChunk)
	conversation := chatInstance.BeginConversation(c.prompt)

	conversation.WaitMyTurn()

	return nil
}

func (c *NonInteractiveCmd) readContextFromStdIn() (string, error) {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}

	return string(inputBytes), nil
}

func (c *NonInteractiveCmd) printChunk(chunk *chat.ConversationChunk) {
	c.ui.PrintMessage(chunk.Content, cwcui.MessageTypeInfo)
}
