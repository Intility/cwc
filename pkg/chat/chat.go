package chat

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"

	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/tools"
	"github.com/intility/cwc/pkg/ui"
)

type Chat struct {
	client        *openai.Client
	systemMessage string
	chunkHandler  MessageChunkHandler
	kit           *tools.Toolkit
}

type MessageChunkHandler func(chunk *ConversationChunk)

func NewChat(client *openai.Client, systemMessage string, onChunk MessageChunkHandler) *Chat {
	return &Chat{
		client:        client,
		systemMessage: systemMessage,
		chunkHandler:  onChunk,
		kit:           nil,
	}
}

func (c *Chat) UseToolkit(kit *tools.Toolkit) {
	c.kit = kit
}

func (c *Chat) BeginConversation(initialMessage string) *Conversation {
	conversation := &Conversation{
		client:  c.client,
		wg:      sync.WaitGroup{},
		onChunk: c.chunkHandler,
		tools:   make([]tools.Tool, 0),
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: c.systemMessage,
			},
		},
	}

	if c.kit != nil {
		for _, tool := range c.kit.ListTools() {
			conversation.tools = append(conversation.tools, *tool)
		}
	}

	conversation.Reply(initialMessage)

	return conversation
}

type Conversation struct {
	client   *openai.Client
	messages []openai.ChatCompletionMessage
	wg       sync.WaitGroup
	tools    []tools.Tool
	onChunk  func(chunk *ConversationChunk)
}

func (c *Conversation) addMessage(role string, message string) {
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    role,
		Content: message,
	})
}

type ConversationChunk struct {
	Role           string
	Content        string
	IsInitialChunk bool
	IsFinalChunk   bool
	IsErrorChunk   bool
}

func (c *Conversation) OnMessageChunk(onChunk func(chunk *ConversationChunk)) {
	c.onChunk = onChunk
}

func (c *Conversation) WaitMyTurn() {
	c.wg.Wait()
}

func (c *Conversation) Reply(message string) {
	c.wg.Add(1)

	c.addMessage(openai.ChatMessageRoleUser, message)

	ctx := context.Background()

	go func() {
		c.onChunk(&ConversationChunk{
			Role:           openai.ChatMessageRoleAssistant,
			Content:        "",
			IsInitialChunk: true,
			IsFinalChunk:   false,
			IsErrorChunk:   false,
		})

		err := c.processMessages(ctx)
		if err != nil {
			c.onChunk(&ConversationChunk{
				Role:           openai.ChatMessageRoleAssistant,
				Content:        "Sorry, I'm having trouble processing your request: " + err.Error(),
				IsInitialChunk: false,
				IsFinalChunk:   true,
				IsErrorChunk:   true,
			})
		}

		c.wg.Done()
	}()
}

func (c *Conversation) processMessages(ctx context.Context) error {
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4TurboPreview,
		Messages: c.messages,
		Stream:   true,
	}

	if len(c.tools) > 0 {
		var openaiTools []openai.Tool

		for _, tool := range c.tools {
			def := tool.Definition()
			openaiTools = append(openaiTools, openai.Tool{
				Type:     openai.ToolTypeFunction,
				Function: &def,
			})
		}

		req.Tools = openaiTools
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("error creating chat completion stream: %w", err)
	}

	defer stream.Close()

	return c.handleStream(ctx, stream)
}

func (c *Conversation) handleStream(ctx context.Context, stream *openai.ChatCompletionStream) error { //nolint:funlen
	var reply strings.Builder

	callDetector := tools.NewToolCallDetector()
	toolWasCalled := false
answer:
	for {
		response, err := stream.Recv()
		if stderrors.Is(err, io.EOF) {
			c.onChunk(&ConversationChunk{
				Role:           openai.ChatMessageRoleAssistant,
				Content:        "",
				IsInitialChunk: false,
				IsFinalChunk:   true,
				IsErrorChunk:   false,
			})

			break answer
		}

		if err != nil {
			return fmt.Errorf("error receiving chat completion response: %w", err)
		}

		if len(response.Choices) == 0 {
			continue answer
		}

		callDetector.Collect(response)

		// check for tool_calls response
		if callDetector.IsToolCallReady() {
			err = c.handleToolCalls(callDetector)

			if err != nil {
				ui.PrintMessage("Error handling tool calls: "+err.Error(), ui.MessageTypeError)

				continue
			}

			toolWasCalled = true

			break
		}

		reply.WriteString(response.Choices[0].Delta.Content)

		c.onChunk(&ConversationChunk{
			Role:           response.Choices[0].Delta.Role,
			Content:        response.Choices[0].Delta.Content,
			IsInitialChunk: false,
			IsFinalChunk:   false,
			IsErrorChunk:   false,
		})
	}

	if toolWasCalled {
		err := c.processMessages(ctx)
		if err != nil {
			return fmt.Errorf("error processing messages: %w", err)
		}

		return nil
	}

	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: reply.String(),
	})

	return nil
}

func (c *Conversation) handleToolCalls(detector *tools.ToolCallDetector) error {
	toolCall := detector.DetectedToolCall()
	if toolCall == nil {
		return errors.NoToolCallsDetectedError{}
	}

	defer detector.Flush()

	// reconstruct the assistant message from the streamed responses
	// gathered by the detector
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: "",
		ToolCalls: []openai.ToolCall{
			{
				Index: toolCall.Index,
				ID:    toolCall.ID,
				Type:  openai.ToolTypeFunction,
				Function: openai.FunctionCall{
					Name:      toolCall.Name,
					Arguments: toolCall.Args,
				},
			},
		},
	})

	for _, tool := range c.tools {
		if toolCall.Name == tool.Definition().Name {
			ui.PrintMessage("[executing tool: "+tool.Definition().Name+"] ", ui.MessageTypeSuccess)

			if tool.HasShellExecutables() {
				toolExecutor := tools.NewShellExecutor()

				toolResponse, err := toolExecutor.Execute(tool, toolCall.Args)
				if err != nil {
					return fmt.Errorf("error executing tool: %w", err)
				}

				c.messages = append(c.messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    toolResponse,
					ToolCallID: toolCall.ID,
				})
			}

			if tool.HasWebExecutables() {
				c.messages = append(c.messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    "Web tool execution not yet supported",
					ToolCallID: toolCall.ID,
				})
			}
		}
	}

	return nil
}
