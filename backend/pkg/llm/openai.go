package llm

import (
	"context"
	"encoding/json"
	"net/http"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
)

const (
	openAIResponsesPath = "/v1/responses"
	headerAuthorization = "Authorization"
	bearerPrefix        = "Bearer "

	openAIResponseCreated    = "response.created"
	openAITextDelta          = "response.output_text.delta"
	openAIResponseCompleted  = "response.completed"
	openAIResponseIncomplete = "response.incomplete"
	openAIResponseFailed     = "response.failed"
	openAIErrorEvent         = "error"
	openAIMaxOutputReason    = "max_output_tokens"
)

type OpenAI struct {
	baseURL string
	apiKey  string
	effort  string
	client  *http.Client
}

func NewOpenAI(cfg configs.ProviderConfig, client *http.Client) *OpenAI {
	return &OpenAI{baseURL: cfg.OpenAIBaseURL, apiKey: cfg.OpenAIAPIKey, effort: cfg.OpenAIReasoningEffort, client: client}
}

type openAIBody struct {
	Model           string          `json:"model"`
	Input           []Message       `json:"input"`
	Stream          bool            `json:"stream"`
	Store           bool            `json:"store"`
	MaxOutputTokens int             `json:"max_output_tokens"`
	Reasoning       openAIReasoning `json:"reasoning"`
}

type openAIReasoning struct {
	Effort string `json:"effort"`
}

type openAIError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type openAIEvent struct {
	Type     string      `json:"type"`
	Delta    string      `json:"delta"`
	Error    openAIError `json:"error"`
	Response struct {
		ID    string `json:"id"`
		Usage struct {
			InputTokens  int64 `json:"input_tokens"`
			OutputTokens int64 `json:"output_tokens"`
		} `json:"usage"`
		IncompleteDetails struct {
			Reason string `json:"reason"`
		} `json:"incomplete_details"`
		Error openAIError `json:"error"`
	} `json:"response"`
}

func (o *OpenAI) Stream(ctx context.Context, req Request, onDelta func(string) error) (Result, error) {
	var (
		result    Result
		completed bool
	)

	err := streamSSE(ctx, o.client, target{
		provider: constants.ProviderOpenAI,
		url:      o.baseURL + openAIResponsesPath,
		headers:  map[string]string{headerAuthorization: bearerPrefix + o.apiKey},
		body: openAIBody{
			Model:           req.Model,
			Input:           req.Messages,
			Stream:          true,
			MaxOutputTokens: req.MaxOutputTokens,
			Reasoning:       openAIReasoning{Effort: o.effort},
		},
	}, func(_, data string) error {
		var event openAIEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return failure.ErrWithTrace(err)
		}

		done, err := o.apply(event, onDelta, &result)
		completed = completed || done

		return err
	})
	if err != nil {
		return result, err
	}
	if !completed {
		return result, failure.ErrWithTrace(errStreamEnded)
	}

	return result, nil
}

func (o *OpenAI) apply(event openAIEvent, onDelta func(string) error, result *Result) (bool, error) {
	switch event.Type {
	case openAIResponseCreated:
		result.ProviderRequestID = event.Response.ID
	case openAITextDelta:
		if event.Delta != "" {
			return false, onDelta(event.Delta)
		}
	case openAIResponseCompleted, openAIResponseIncomplete:
		result.InputTokens = event.Response.Usage.InputTokens
		result.OutputTokens = event.Response.Usage.OutputTokens
		result.Truncated = event.Response.IncompleteDetails.Reason == openAIMaxOutputReason

		return true, nil
	case openAIResponseFailed:
		return false, &APIError{Provider: constants.ProviderOpenAI, Type: event.Response.Error.Type, Message: event.Response.Error.Message}
	case openAIErrorEvent:
		return false, &APIError{Provider: constants.ProviderOpenAI, Type: event.Error.Type, Message: event.Error.Message}
	}

	return false, nil
}
