package internal

import (
	"fmt"

	"github.com/intility/cwc/pkg/chat"
	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/prompting"
	"github.com/intility/cwc/pkg/systemcontext"
	"github.com/intility/cwc/pkg/ui"
)

type NonInteractiveCmd struct {
	ui             ui.UI
	clientProvider config.ClientProvider
	promptResolver prompting.PromptResolver
	smGenerator    systemcontext.SystemMessageGenerator
}

func NewNonInteractiveCmd(
	clientProvider config.ClientProvider,
	promptResolver prompting.PromptResolver,
	smGenerator systemcontext.SystemMessageGenerator,
) *NonInteractiveCmd {
	return &NonInteractiveCmd{
		ui:             ui.NewUI(),
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

	generateSystemMessage, err := c.smGenerator.GenerateSystemMessage()
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

func (c *NonInteractiveCmd) printChunk(chunk *chat.ConversationChunk) {
	c.ui.PrintMessage(chunk.Content, ui.MessageTypeInfo)
}
