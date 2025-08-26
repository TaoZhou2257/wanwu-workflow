package wanwu_util

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/cloudwego/eino-ext/libs/acl/openai"
	model "github.com/coze-dev/coze-studio/backend/api/model/crossdomain/modelmgr"
	"github.com/coze-dev/coze-studio/backend/infra/contract/chatmodel"
	"github.com/coze-dev/coze-studio/backend/infra/contract/modelmgr"
	chatmodelImpl "github.com/coze-dev/coze-studio/backend/infra/impl/chatmodel"
)

func CreateChatModel(ctx context.Context, llmParams *model.LLMParams) (chatmodel.ToolCallingChatModel, *modelmgr.Model, error) {
	if llmParams == nil {
		return nil, nil, errors.New("empty llmParams")
	}
	var topP *float32
	if llmParams.TopP != nil {
		topPF32 := float32(*llmParams.TopP)
		topP = &topPF32
	}
	var temperature *float32
	if llmParams.Temperature != nil {
		temperatureF32 := float32(*llmParams.Temperature)
		temperature = &temperatureF32
	}
	var maxTokens *int
	if llmParams.MaxTokens != 0 {
		maxTokens = &llmParams.MaxTokens
	}
	var responseFormatType openai.ChatCompletionResponseFormatType
	switch llmParams.ResponseFormat {
	case model.ResponseFormatText:
		responseFormatType = openai.ChatCompletionResponseFormatTypeText
	case model.ResponseFormatMarkdown:
		responseFormatType = openai.ChatCompletionResponseFormatTypeJSONObject
	case model.ResponseFormatJSON:
		responseFormatType = openai.ChatCompletionResponseFormatTypeJSONObject
	default:
		responseFormatType = openai.ChatCompletionResponseFormatTypeText
	}
	baseUrl, err := url.JoinPath(os.Getenv("WANWU_CALLBACK_LLM_BASE_URL"), strconv.Itoa(int(llmParams.ModelType)))
	if err != nil {
		return nil, nil, err
	}
	m, err := chatmodelImpl.NewDefaultFactory().CreateChatModel(ctx, chatmodel.ProtocolOpenAI, &chatmodel.Config{
		BaseURL:     baseUrl,
		Model:       llmParams.ModelName,
		TopP:        topP,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		OpenAI:      &chatmodel.OpenAIConfig{ResponseFormat: &openai.ChatCompletionResponseFormat{Type: responseFormatType}},
	})
	return m, nil, err
}
