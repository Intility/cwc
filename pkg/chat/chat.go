package chat

import (
	"context"
	stderrors "errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"strings"
	"sync"
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
	go func() {
		err := c.processMessages(context.Background())
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
		Model: openai.GPT4TurboPreview,
		//MaxTokens: 4096,
		Messages: c.messages,
		Stream:   true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	defer stream.Close()

	var reply strings.Builder

	c.onChunk(&ConversationChunk{
		Role:           openai.ChatMessageRoleAssistant,
		Content:        "",
		IsInitialChunk: true,
		IsFinalChunk:   false,
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
			})
			break answer
		}

		if err != nil {
			return err
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
		})
	}

	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: reply.String(),
	})

	return nil
}
