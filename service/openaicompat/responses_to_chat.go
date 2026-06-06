package openaicompat

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/samber/lo"
)

// ResponsesRequestToChatCompletionsRequest converts an OpenAI Responses API request
// to a Chat Completions API request. This allows Responses API consumers (like Codex CLI)
// to work with providers that only support the Chat Completions format (like DeepSeek).
func ResponsesRequestToChatCompletionsRequest(req *dto.OpenAIResponsesRequest) (*dto.GeneralOpenAIRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	messages := make([]dto.Message, 0)

	// Convert Instructions to system message
	if len(req.Instructions) > 0 {
		var instructionsStr string
		if err := common.Unmarshal(req.Instructions, &instructionsStr); err == nil {
			if strings.TrimSpace(instructionsStr) != "" {
				msg := dto.Message{Role: "system"}
				msg.SetStringContent(instructionsStr)
				messages = append(messages, msg)
			}
		}
	}

	// Parse Input items and convert to messages
	if len(req.Input) > 0 {
		inputMessages := parseResponsesInput(req.Input)
		messages = append(messages, inputMessages...)
	}

	// Build tools as ToolCallRequest
	var tools []dto.ToolCallRequest
	if len(req.Tools) > 0 {
		var rawTools []map[string]any
		if err := common.Unmarshal(req.Tools, &rawTools); err == nil {
			for _, t := range rawTools {
				toolType, _ := t["type"].(string)
				if toolType == "function" {
					tool := dto.ToolCallRequest{
						Type: "function",
					}
					if name, ok := t["name"].(string); ok {
						tool.Function.Name = name
					}
					if desc, ok := t["description"].(string); ok {
						tool.Function.Description = desc
					}
					if params, ok := t["parameters"]; ok {
						tool.Function.Parameters = params
					}
					tools = append(tools, tool)
				}
			}
		}
	}

	// Tool choice
	var toolChoice any
	if len(req.ToolChoice) > 0 {
		var rawChoice map[string]any
		if err := common.Unmarshal(req.ToolChoice, &rawChoice); err == nil {
			if typeVal, ok := rawChoice["type"].(string); ok && typeVal == "function" {
				if name, ok := rawChoice["name"].(string); ok {
					toolChoice = map[string]any{
						"type": "function",
						"function": map[string]string{
							"name": name,
						},
					}
				}
			} else {
				toolChoice = string(req.ToolChoice)
			}
		} else {
			var str string
			if err := common.Unmarshal(req.ToolChoice, &str); err == nil {
				toolChoice = str
			} else {
				toolChoice = string(req.ToolChoice)
			}
		}
	}

	// Max tokens
	var maxTokens *uint
	if req.MaxOutputTokens != nil {
		maxTokens = lo.ToPtr(lo.FromPtr(req.MaxOutputTokens))
	}

	// Reasoning effort
	var reasoningEffort string
	if req.Reasoning != nil && req.Reasoning.Effort != "" {
		reasoningEffort = req.Reasoning.Effort
	}

	// Build the chat request
	out := &dto.GeneralOpenAIRequest{
		Model:           req.Model,
		Messages:        messages,
		Stream:          lo.ToPtr(lo.FromPtrOr(req.Stream, false)),
		MaxTokens:       maxTokens,
		Temperature:     req.Temperature,
		TopP:            req.TopP,
		Tools:           tools,
		ToolChoice:      toolChoice,
		User:            req.User,
		Metadata:        req.Metadata,
		Store:           req.Store,
		ReasoningEffort: reasoningEffort,
	}

	// Text format (response_format)
	if len(req.Text) > 0 {
		var textConfig map[string]any
		if err := common.Unmarshal(req.Text, &textConfig); err == nil {
			if format, ok := textConfig["format"].(map[string]any); ok {
				if formatType, ok := format["type"].(string); ok {
					out.ResponseFormat = &dto.ResponseFormat{
						Type: formatType,
					}
				}
			}
		}
	}

	return out, nil
}

// parseResponsesInput converts Responses API input items into chat messages.
func parseResponsesInput(input []byte) []dto.Message {
	if len(input) == 0 {
		return nil
	}

	var messages []dto.Message

	// Try parsing as array of typed items
	var items []map[string]any
	if err := common.Unmarshal(input, &items); err != nil {
		// Try as plain string
		var str string
		if err := common.Unmarshal(input, &str); err == nil && strings.TrimSpace(str) != "" {
			msg := dto.Message{Role: "user"}
			msg.SetStringContent(str)
			messages = append(messages, msg)
		}
		return messages
	}

	var currentToolCalls []dto.ToolCallRequest

	for _, item := range items {
		itemType, _ := item["type"].(string)

		switch itemType {
		case "message":
			role, _ := item["role"].(string)
			content := extractContentString(item["content"])

			msg := dto.Message{Role: role}
			msg.SetStringContent(content)
			messages = append(messages, msg)

		case "function_call":
			callID, _ := item["call_id"].(string)
			name, _ := item["name"].(string)
			args, _ := item["arguments"].(string)

			if callID != "" && name != "" {
				currentToolCalls = append(currentToolCalls, dto.ToolCallRequest{
					ID:   callID,
					Type: "function",
					Function: dto.FunctionRequest{
						Name:      name,
						Arguments: args,
					},
				})
			}

		case "function_call_output":
			// Flush pending tool calls to last assistant message
			if len(currentToolCalls) > 0 && len(messages) > 0 {
				lastMsg := &messages[len(messages)-1]
				if lastMsg.Role == "assistant" {
					lastMsg.SetToolCalls(currentToolCalls)
				}
				currentToolCalls = nil
			}

			callID, _ := item["call_id"].(string)
			output, _ := item["output"].(string)
			msg := dto.Message{
				Role:       "tool",
				ToolCallId: callID,
			}
			msg.SetStringContent(output)
			messages = append(messages, msg)

		case "input_text":
			text, _ := item["text"].(string)
			msg := dto.Message{Role: "user"}
			msg.SetStringContent(text)
			messages = append(messages, msg)

		case "output_text":
			text, _ := item["text"].(string)
			if len(messages) > 0 && messages[len(messages)-1].Role == "assistant" {
				lastMsg := &messages[len(messages)-1]
				if lastMsg.IsStringContent() {
					existing := lastMsg.StringContent()
					if existing != "" {
						lastMsg.SetStringContent(existing + "\n" + text)
					} else {
						lastMsg.SetStringContent(text)
					}
				}
			} else {
				msg := dto.Message{Role: "assistant"}
				msg.SetStringContent(text)
				messages = append(messages, msg)
			}
		}
	}

	// Flush remaining tool calls
	if len(currentToolCalls) > 0 && len(messages) > 0 {
		lastMsg := &messages[len(messages)-1]
		if lastMsg.Role == "assistant" {
			lastMsg.SetToolCalls(currentToolCalls)
		}
	}

	return messages
}

// extractContentString extracts a string representation from content
// which can be a string, array of content parts, or nil.
func extractContentString(content any) string {
	if content == nil {
		return ""
	}

	switch v := content.(type) {
	case string:
		return v
	case []any:
		var textParts []string
		for _, part := range v {
			if m, ok := part.(map[string]any); ok {
				typeVal, _ := m["type"].(string)
				switch typeVal {
				case "input_text", "text":
					if text, ok := m["text"].(string); ok {
						textParts = append(textParts, text)
					}
				}
			}
		}
		return strings.Join(textParts, "\n")
	}

	return fmt.Sprintf("%v", content)
}
