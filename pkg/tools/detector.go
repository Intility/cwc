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

func (t *ToolCallDetector) DetectedToolCalls() []ToolCall {
	if !t.isReady {
		return nil
	}

	// construct the name and args from the responses
	toolCalls := make([]ToolCall, 0)

	var currentToolCall ToolCall

	// split the responses based on the presence of ToolCalls.ID.
	// This separates the tool calls from each other as multiple may be present in the responses.

	for _, res := range t.responses {
		for _, calls := range res.Choices[0].Delta.ToolCalls {
			if calls.ID != "" {
				if currentToolCall.ID != "" {
					toolCalls = append(toolCalls, currentToolCall)
					currentToolCall = ToolCall{Index: nil, ID: "", Name: "", Args: ""}
				}

				currentToolCall.Index = calls.Index
				currentToolCall.ID = calls.ID
				currentToolCall.Name = calls.Function.Name
			} else {
				currentToolCall.Name += calls.Function.Name // in case the name is split across multiple responses
				currentToolCall.Args += calls.Function.Arguments
			}
		}
	}

	if currentToolCall.ID != "" {
		toolCalls = append(toolCalls, currentToolCall)
	}

	return toolCalls
}

func (t *ToolCallDetector) Flush() {
	t.responses = make([]openai.ChatCompletionStreamResponse, 0)
	t.isReady = false
}
