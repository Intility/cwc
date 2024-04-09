package internal

import (
	"fmt"
	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/chat"
	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/prompting"
	"github.com/intility/cwc/pkg/systemcontext"
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
	ui             ui.UI
	clientProvider config.ClientProvider
	promptResolver prompting.PromptResolver
	smGenerator    systemcontext.SystemMessageGenerator
	chatOptions    InteractiveChatOptions
}

func NewInteractiveCmd(
	promptResolver prompting.PromptResolver,
	clientProvider config.ClientProvider,
	smGenerator systemcontext.SystemMessageGenerator,
	chatOptions InteractiveChatOptions,
) *InteractiveCmd {
	return &InteractiveCmd{
		ui:             ui.NewUI(),
		promptResolver: promptResolver,
		clientProvider: clientProvider,
		chatOptions:    chatOptions,
		smGenerator:    smGenerator,
	}
}

func (c *InteractiveCmd) Run() error {
	openaiClient, err := c.clientProvider.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("error creating openaiClient: %w", err)
	}

	generatedSystemMessage, err := c.smGenerator.GenerateSystemMessage()
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	c.ui.PrintMessage("Type '/exit' to end the chat.\n", ui.MessageTypeNotice)

	userPrompt := c.promptResolver.ResolvePrompt()

	if userPrompt == "" {
		c.ui.PrintMessage("👤: ", ui.MessageTypeInfo)
		userPrompt = c.ui.ReadUserInput()
	} else {
		c.ui.PrintMessage(fmt.Sprintf("👤: %s\n", userPrompt), ui.MessageTypeInfo)
	}

	if userPrompt == "/exit" {
		return nil
	}

	c.handleChat(openaiClient, generatedSystemMessage, userPrompt)

	return nil
}

func (c *InteractiveCmd) handleChat(client *openai.Client, systemMessage string, prompt string) {
	chatInstance := chat.NewChat(client, systemMessage, c.printMessageChunk)
	conversation := chatInstance.BeginConversation(prompt)

	for {
		conversation.WaitMyTurn()
		c.ui.PrintMessage("👤: ", ui.MessageTypeInfo)

		userMessage := c.ui.ReadUserInput()

		if userMessage == "/exit" {
			break
		}

		conversation.Reply(userMessage)
	}
}

func (c *InteractiveCmd) printMessageChunk(chunk *chat.ConversationChunk) {
	if chunk.IsInitialChunk {
		c.ui.PrintMessage("🤖: ", ui.MessageTypeInfo)
		return
	}

	if chunk.IsErrorChunk {
		c.ui.PrintMessage(chunk.Content, ui.MessageTypeError)
	}

	if chunk.IsFinalChunk {
		c.ui.PrintMessage("\n", ui.MessageTypeInfo)
	}

	c.ui.PrintMessage(chunk.Content, ui.MessageTypeInfo)
}
