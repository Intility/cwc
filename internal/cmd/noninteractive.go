package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/intility/cwc/internal/client"
	"github.com/intility/cwc/internal/prompt"
	"github.com/intility/cwc/internal/systemmessage"
	"github.com/intility/cwc/pkg/chat"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/ui"
)

type NonInteractiveCmd struct {
	clientProvider client.Provider
	promptResolver prompt.Resolver
	smGenerator    systemmessage.Generator
}

func NewNonInteractiveCmd(
	clientProvider client.Provider,
	promptResolver prompt.Resolver,
	smGenerator systemmessage.Generator,
) *NonInteractiveCmd {
	return &NonInteractiveCmd{
		clientProvider: clientProvider,
		promptResolver: promptResolver,
		smGenerator:    smGenerator,
	}
}

func (c *NonInteractiveCmd) Run() error {
	openaiClient, err := c.clientProvider.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("error creating openaiClient: %w", err)
	}

	systemCtx, err := c.readContextFromStdIn()
	if err != nil {
		return fmt.Errorf("error reading context from stdin: %w", err)
	}

	generateSystemMessage, err := c.smGenerator.GenerateSystemMessage(systemCtx)
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	userPrompt := c.promptResolver.ResolvePrompt()

	if userPrompt == "" {
		return errors.NoPromptProvidedError{Message: "non-interactive mode requires a prompt"}
	}

	chatInstance := chat.NewChat(openaiClient, generateSystemMessage, c.printChunk)
	conversation := chatInstance.BeginConversation(userPrompt)

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
	ui.PrintMessage(chunk.Content, ui.MessageTypeInfo)
}
