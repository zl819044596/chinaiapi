package deepseek

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/claude"
	"github.com/QuantumNous/new-api/relay/channel/openai"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/service/openaicompat"
	"github.com/QuantumNous/new-api/setting/reasoning"
	"github.com/QuantumNous/new-api/types"
	"github.com/gin-gonic/gin"
)

type Adaptor struct {
}

func (a *Adaptor) ConvertGeminiRequest(*gin.Context, *relaycommon.RelayInfo, *dto.GeminiChatRequest) (any, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, req *dto.ClaudeRequest) (any, error) {
	adaptor := claude.Adaptor{}
	convertedRequest, err := adaptor.ConvertClaudeRequest(c, info, req)
	if err != nil {
		return nil, err
	}
	claudeRequest, ok := convertedRequest.(*dto.ClaudeRequest)
	if !ok {
		return convertedRequest, nil
	}
	if err := applyDeepSeekV4ClaudeThinkingSuffix(info, claudeRequest); err != nil {
		return nil, err
	}
	return claudeRequest, nil
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	fimBaseUrl := info.ChannelBaseUrl
	switch info.RelayFormat {
	case types.RelayFormatClaude:
		return fmt.Sprintf("%s/anthropic/v1/messages", info.ChannelBaseUrl), nil
	default:
		if !strings.HasSuffix(info.ChannelBaseUrl, "/beta") {
			fimBaseUrl += "/beta"
		}
		switch info.RelayMode {
		case constant.RelayModeCompletions:
			return fmt.Sprintf("%s/completions", fimBaseUrl), nil
		default:
			return fmt.Sprintf("%s/v1/chat/completions", info.ChannelBaseUrl), nil
		}
	}
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)
	req.Set("Authorization", "Bearer "+info.ApiKey)
	return nil
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	if err := applyDeepSeekV4OpenAIThinkingSuffix(info, request); err != nil {
		return nil, err
	}

	return request, nil
}

func applyDeepSeekV4OpenAIThinkingSuffix(info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) error {
	modelName := request.Model
	if info != nil && info.ChannelMeta != nil && info.UpstreamModelName != "" {
		modelName = info.UpstreamModelName
	}
	baseModel, thinkingType, effort, ok := reasoning.ParseDeepSeekV4ThinkingSuffix(modelName)
	if !ok {
		return nil
	}
	thinking, err := common.Marshal(map[string]string{
		"type": thinkingType,
	})
	if err != nil {
		return fmt.Errorf("error marshalling thinking: %w", err)
	}
	request.Model = baseModel
	request.THINKING = thinking
	request.ReasoningEffort = effort
	if info != nil {
		if info.ChannelMeta != nil {
			info.UpstreamModelName = baseModel
		}
		info.ReasoningEffort = effort
	}
	return nil
}

func applyDeepSeekV4ClaudeThinkingSuffix(info *relaycommon.RelayInfo, request *dto.ClaudeRequest) error {
	modelName := request.Model
	if info != nil && info.ChannelMeta != nil && info.UpstreamModelName != "" {
		modelName = info.UpstreamModelName
	}
	baseModel, thinkingType, effort, ok := reasoning.ParseDeepSeekV4ThinkingSuffix(modelName)
	if !ok {
		return nil
	}
	request.Model = baseModel
	request.Thinking = &dto.Thinking{Type: thinkingType}
	if effort == "" {
		request.OutputConfig = nil
	} else {
		outputConfig, err := common.Marshal(map[string]string{
			"effort": effort,
		})
		if err != nil {
			return fmt.Errorf("error marshalling output_config: %w", err)
		}
		request.OutputConfig = outputConfig
	}
	if info != nil {
		if info.ChannelMeta != nil {
			info.UpstreamModelName = baseModel
		}
		info.ReasoningEffort = effort
	}
	return nil
}

func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return nil, nil
}

func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	// Convert Responses API request to Chat Completions format
	// since DeepSeek doesn't natively support the Responses API
	chatReq, err := openaicompat.ResponsesRequestToChatCompletionsRequest(&request)
	if err != nil {
		return nil, fmt.Errorf("deepseek: failed to convert responses request to chat completions: %w", err)
	}
	// Now convert from Chat Completions to DeepSeek format
	return a.ConvertOpenAIRequest(c, info, chatReq)
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	switch info.RelayFormat {
	case types.RelayFormatClaude:
		adaptor := claude.Adaptor{}
		return adaptor.DoResponse(c, resp, info)
	default:
		// For Responses API relay mode, DeepSeek returns Chat Completions format
		// but the client expects Responses format. Convert the response.
		if info.RelayMode == constant.RelayModeResponses || info.RelayMode == constant.RelayModeResponsesCompact {
			return handleDeepSeekResponsesResponse(c, resp, info)
		}
		adaptor := openai.Adaptor{}
		return adaptor.DoResponse(c, resp, info)
	}
}

// handleDeepSeekResponsesResponse converts a Chat Completions response from DeepSeek
// to the Responses API format expected by the client (e.g., Codex CLI).
func handleDeepSeekResponsesResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.Usage, *types.NewAPIError) {
	defer service.CloseResponseBodyGracefully(resp)

	// Read the Chat Completions response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeReadResponseBodyFailed, http.StatusInternalServerError)
	}

	// Parse as Chat Completions response
	var textResp dto.OpenAITextResponse
	if err := common.Unmarshal(responseBody, &textResp); err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeBadResponseBody, http.StatusInternalServerError)
	}

	// Build a Responses API response from the Chat Completions response
	outputItems := make([]map[string]any, 0)

	// First, check for reasoning content
	reasoningContent := ""
	for _, choice := range textResp.Choices {
		msg := choice.Message
		if msg.Role == "assistant" {
			// Add reasoning content if present
			if rc := msg.GetReasoningContent(); rc != "" {
				reasoningContent = rc
			}
		}
	}

	// Build output items from choices
	for i, choice := range textResp.Choices {
		msg := choice.Message
		if msg.Role == "assistant" {
			// Add tool calls
			for _, tc := range msg.ParseToolCalls() {
				if tc.ID != "" && tc.Function.Name != "" {
					outputItems = append(outputItems, map[string]any{
						"type":      "function_call",
						"call_id":   tc.ID,
						"name":      tc.Function.Name,
						"arguments": tc.Function.Arguments,
					})
				}
			}

			// Add the text content
			if msg.IsStringContent() {
				text := msg.StringContent()
				if text != "" {
					item := map[string]any{
						"type": "output_text",
						"text": text,
					}
					// Include annotations if reasoning content exists
					if reasoningContent != "" {
						item["annotations"] = []map[string]any{
							{
								"type":        "reasoning",
								"description": reasoningContent,
							},
						}
					}
					outputItems = append(outputItems, item)
				}
			}
		}
		_ = i // suppress unused
	}

	// If no output items were generated, create a default one
	if len(outputItems) == 0 {
		outputItems = append(outputItems, map[string]any{
			"type": "output_text",
			"text": "",
		})
	}

	// Build usage
	usage := &dto.Usage{
		PromptTokens:     textResp.Usage.PromptTokens,
		CompletionTokens: textResp.Usage.CompletionTokens,
		TotalTokens:      textResp.Usage.TotalTokens,
	}
	if textResp.Usage.PromptTokensDetails.CachedTokens > 0 {
		usage.PromptTokensDetails.CachedTokens = textResp.Usage.PromptTokensDetails.CachedTokens
	}

	// Build the Responses API response
	responsesResponse := map[string]any{
		"id":      textResp.Id,
		"object":  "response",
		"created": textResp.Created,
		"model":   textResp.Model,
		"output":  outputItems,
		"usage": map[string]any{
			"input_tokens":   textResp.Usage.PromptTokens,
			"output_tokens":  textResp.Usage.CompletionTokens,
			"total_tokens":   textResp.Usage.TotalTokens,
			"input_cached_tokens": textResp.Usage.PromptTokensDetails.CachedTokens,
		},
		"status": "completed",
	}

	if info.RelayMode == constant.RelayModeResponsesCompact {
		// For compact mode, add additional fields
		responsesResponse["conversation_id"] = ""
	}

	jsonData, err := common.Marshal(responsesResponse)
	if err != nil {
		return nil, types.NewOpenAIError(err, types.ErrorCodeInvalidRequest, http.StatusInternalServerError)
	}

	logger.LogDebug(c, "responses response: %s", jsonData)
	service.IOCopyBytesGracefully(c, resp, jsonData)

	return usage, nil
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}
