package chat

import (
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

type Chat struct {
	client        *openai.Client
	systemMessage string
	chunkHandler  MessageChunkHandler
}

type MessageChunkHandler func(chunk *ConversationChunk)

func NewChat(client *openai.Client, systemMessage string, onChunk MessageChunkHandler) *Chat {
	return &Chat{
		client:        client,
		systemMessage: systemMessage,
		chunkHandler:  onChunk,
	}
}

func (c *Chat) BeginConversation(initialMessage string) *Conversation {
	conversation := &Conversation{
		client:  c.client,
		wg:      sync.WaitGroup{},
		onChunk: c.chunkHandler,
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: c.systemMessage,
			},
		},
	}

	conversation.Reply(initialMessage)

	return conversation
}

type Conversation struct {
	client   *openai.Client
	messages []openai.ChatCompletionMessage
	wg       sync.WaitGroup
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

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("error creating chat completion stream: %w", err)
	}

	defer stream.Close()

	return c.handleStream(stream)
}

func (c *Conversation) handleStream(stream *openai.ChatCompletionStream) error {
	var reply strings.Builder

	c.onChunk(&ConversationChunk{
		Role:           openai.ChatMessageRoleAssistant,
		Content:        "",
		IsInitialChunk: true,
		IsFinalChunk:   false,
		IsErrorChunk:   false,
	})

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

		reply.WriteString(response.Choices[0].Delta.Content)

		c.onChunk(&ConversationChunk{
			Role:           response.Choices[0].Delta.Role,
			Content:        response.Choices[0].Delta.Content,
			IsInitialChunk: false,
			IsFinalChunk:   false,
			IsErrorChunk:   false,
		})
	}

	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: reply.String(),
	})

	return nil
}
