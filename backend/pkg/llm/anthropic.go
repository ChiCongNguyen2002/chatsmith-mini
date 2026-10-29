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
	anthropicMessagesPath = "/v1/messages"
	anthropicVersion      = "2023-06-01"
	headerAnthropicKey    = "x-api-key"
	headerAnthropicVer    = "anthropic-version"

	anthropicMessageStart = "message_start"
	anthropicBlockStart   = "content_block_start"
	anthropicBlockDelta   = "content_block_delta"
	anthropicMessageDelta = "message_delta"
	anthropicMessageStop  = "message_stop"
	anthropicError        = "error"
	anthropicTextDelta    = "text_delta"
	anthropicStopMaxToken = "max_tokens"
)

type Anthropic struct {
	baseURL string
	apiKey  string
	effort  string
	client  *http.Client
}

func NewAnthropic(cfg configs.ProviderConfig, client *http.Client) *Anthropic {
	return &Anthropic{baseURL: cfg.AnthropicBaseURL, apiKey: cfg.AnthropicAPIKey, effort: cfg.AnthropicEffort, client: client}
}

type anthropicBody struct {
	Model        string          `json:"model"`
	MaxTokens    int             `json:"max_tokens"`
	Stream       bool            `json:"stream"`
	System       string          `json:"system,omitempty"`
	Messages     []Message       `json:"messages"`
	Tools        []anthropicTool `json:"tools,omitempty"`
	OutputConfig anthropicOutput `json:"output_config"`
}

type anthropicTool struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	MaxUses int    `json:"max_uses,omitempty"`
}

type anthropicOutput struct {
	Effort string `json:"effort"`
}

type anthropicUsage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
}

type anthropicEvent struct {
	Type         string `json:"type"`
	Index        int    `json:"index"`
	ContentBlock struct {
		Type      string          `json:"type"`
		Name      string          `json:"name"`
		ToolUseID string          `json:"tool_use_id"`
		Content   json.RawMessage `json:"content"`
	} `json:"content_block"`
	Message struct {
		ID    string         `json:"id"`
		Usage anthropicUsage `json:"usage"`
	} `json:"message"`
	Delta struct {
		Type       string `json:"type"`
		Text       string `json:"text"`
		StopReason string `json:"stop_reason"`
	} `json:"delta"`
	Usage anthropicUsage `json:"usage"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (a *Anthropic) Stream(ctx context.Context, req Request, onDelta func(string) error) (Result, error) {
	var (
		result    Result
		completed bool
	)

	err := streamSSE(ctx, a.client, target{
		provider: constants.ProviderAnthropic,
		url:      a.baseURL + anthropicMessagesPath,
		headers:  map[string]string{headerAnthropicKey: a.apiKey, headerAnthropicVer: anthropicVersion},
		body:     a.body(req),
	}, func(_, data string) error {
		var event anthropicEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return failure.ErrWithTrace(err)
		}

		done, err := a.apply(event, req, onDelta, &result)
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

func (a *Anthropic) body(req Request) anthropicBody {
	body := anthropicBody{
		Model:        req.Model,
		MaxTokens:    req.MaxOutputTokens,
		Stream:       true,
		System:       req.SystemPrompt,
		Messages:     req.Messages,
		OutputConfig: anthropicOutput{Effort: a.effort},
	}
	if req.WebSearchEnabled {
		body.Tools = []anthropicTool{{Type: searchToolType, Name: searchToolName, MaxUses: maxSearchUses}}
	}

	return body
}

func (a *Anthropic) apply(event anthropicEvent, req Request, onDelta func(string) error, result *Result) (bool, error) {
	switch event.Type {
	case anthropicMessageStart:
		result.ProviderRequestID = event.Message.ID
		result.InputTokens = event.Message.Usage.InputTokens
	case anthropicBlockStart:
		return false, a.handleBlockStart(event, req.OnSearch)
	case anthropicBlockDelta:
		if event.Delta.Type == anthropicTextDelta && event.Delta.Text != "" {
			return false, onDelta(event.Delta.Text)
		}
	case anthropicMessageDelta:
		result.Truncated = event.Delta.StopReason == anthropicStopMaxToken
		result.OutputTokens = event.Usage.OutputTokens
		if event.Usage.InputTokens > 0 {
			result.InputTokens = event.Usage.InputTokens
		}
	case anthropicMessageStop:
		return true, nil
	case anthropicError:
		return false, &APIError{Provider: constants.ProviderAnthropic, Type: event.Error.Type, Message: event.Error.Message}
	}

	return false, nil
}
