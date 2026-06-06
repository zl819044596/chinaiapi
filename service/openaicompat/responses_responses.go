package openaicompat

import (
	"github.com/QuantumNous/new-api/dto"
)

// ResponsesResponseToChatCompletionsResponse converts a Responses API response
// to a Chat Completions text response. This is needed by the chat_via_responses flow.
func ResponsesResponseToChatCompletionsResponse(resp *dto.OpenAIResponsesResponse, id string) (*dto.OpenAITextResponse, *dto.Usage, error) {
	usage := &dto.Usage{}
	if resp.Usage != nil {
		usage.PromptTokens = resp.Usage.InputTokens
		usage.CompletionTokens = resp.Usage.OutputTokens
		usage.TotalTokens = resp.Usage.TotalTokens
		if resp.Usage.InputTokensDetails != nil {
			usage.PromptTokensDetails.CachedTokens = resp.Usage.InputTokensDetails.CachedTokens
		}
	}

	// Extract the response text from output items
	outputText := ExtractOutputTextFromResponses(resp)

	chatResp := &dto.OpenAITextResponse{
		Id:      id,
		Model:   resp.Model,
		Object:  "chat.completion",
		Created: resp.CreatedAt,
		Choices: []dto.OpenAITextResponseChoice{
			{
				Index: 0,
				Message: dto.Message{
					Role:    "assistant",
					Content: outputText,
				},
				FinishReason: "stop",
			},
		},
		Usage: *usage,
	}

	return chatResp, usage, nil
}

// ExtractOutputTextFromResponses extracts the text content from a Responses API response.
func ExtractOutputTextFromResponses(resp *dto.OpenAIResponsesResponse) string {
	if resp == nil {
		return ""
	}

	var textParts []string
	for _, item := range resp.Output {
		if item.Type == "message" {
			for _, content := range item.Content {
				if content.Type == "output_text" {
					textParts = append(textParts, content.Text)
				}
			}
		}
	}

	var result string
	for _, t := range textParts {
		if result != "" {
			result += "\n"
		}
		result += t
	}
	return result
}
