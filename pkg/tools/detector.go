package tools

import (
	"github.com/sashabaranov/go-openai"
)

type ToolCallDetector struct {
	isReady   bool
	responses []openai.ChatCompletionStreamResponse
}

func NewToolCallDetector() *ToolCallDetector {
	return &ToolCallDetector{
		isReady:   false,
		responses: make([]openai.ChatCompletionStreamResponse, 0),
	}
}

func (t *ToolCallDetector) IsToolCallReady() bool {
	// last response has finish reason tool_calls
	if len(t.responses) == 0 {
		return false
	}

	lastResponse := t.responses[len(t.responses)-1]
	t.isReady = lastResponse.Choices[0].FinishReason == openai.FinishReasonToolCalls

	return t.isReady
}

func (t *ToolCallDetector) Collect(res openai.ChatCompletionStreamResponse) {
	t.responses = append(t.responses, res)
}

type ToolCall struct {
	Index *int
	ID    string
	Name  string
	Args  string
}

func (t *ToolCallDetector) DetectedToolCall() *ToolCall {
	if !t.isReady {
		return nil
	}

	// construct the name and args from the responses
	toolCall := &ToolCall{
		Index: nil,
		ID:    "",
		Name:  "",
		Args:  "",
	}

	for _, res := range t.responses {
		for _, calls := range res.Choices[0].Delta.ToolCalls {
			if toolCall.Index == nil && calls.Index != nil {
				toolCall.Index = calls.Index
			}

			toolCall.ID += calls.ID
			toolCall.Name += calls.Function.Name
			toolCall.Args += calls.Function.Arguments
		}
	}

	return toolCall
}

func (t *ToolCallDetector) Flush() {
	t.responses = make([]openai.ChatCompletionStreamResponse, 0)
	t.isReady = false
}
