package provider

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

	anthropicMessageStart    = "message_start"
	anthropicBlockStart      = "content_block_start"
	anthropicBlockDelta      = "content_block_delta"
	anthropicMessageDelta    = "message_delta"
	anthropicMessageStop     = "message_stop"
	anthropicError           = "error"
	anthropicTextDelta       = "text_delta"
	anthropicStopMaxToken    = "max_tokens"
	anthropicServerToolUse   = "server_tool_use"
	anthropicSearchResult    = "web_search_tool_result"
	anthropicSearchResultRow = "web_search_result"

	searchToolType = "web_search_20250305"
	searchToolName = "web_search"
	maxSearchUses  = 3
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
	Type     string             `json:"type"`
	Name     string             `json:"name"`
	MaxUses  int                `json:"max_uses,omitempty"`
	Location *anthropicLocation `json:"user_location,omitempty"`
}

type anthropicLocation struct {
	Type    string `json:"type"`
	Country string `json:"country"`
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

type searchResultItem struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (a *Anthropic) Stream(ctx context.Context, req Request, onDelta func(string) error) (Result, error) {
	var (
		result    Result
		completed bool
	)

	body := anthropicBody{
		Model:        req.Model,
		MaxTokens:    req.MaxOutputTokens,
		Stream:       true,
		System:       req.SystemPrompt,
		Messages:     req.Messages,
		OutputConfig: anthropicOutput{Effort: a.effort},
	}
	if req.WebSearchEnabled {
		body.Tools = []anthropicTool{{
			Type:    searchToolType,
			Name:    searchToolName,
			MaxUses: maxSearchUses,
			Location: &anthropicLocation{
				Type:    "approximate",
				Country: "VN",
			},
		}}
	}

	err := streamSSE(ctx, a.client, target{
		provider: constants.ProviderAnthropic,
		url:      a.baseURL + anthropicMessagesPath,
		headers:  map[string]string{headerAnthropicKey: a.apiKey, headerAnthropicVer: anthropicVersion},
		body:     body,
	}, func(_, data string) error {
		var event anthropicEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			return failure.ErrWithTrace(err)
		}

		switch event.Type {
		case anthropicMessageStart:
			result.ProviderRequestID = event.Message.ID
			result.InputTokens = event.Message.Usage.InputTokens
		case anthropicBlockStart:
			return a.handleBlockStart(event, req.OnSearch)
		case anthropicBlockDelta:
			if event.Delta.Type == anthropicTextDelta && event.Delta.Text != "" {
				return onDelta(event.Delta.Text)
			}
		case anthropicMessageDelta:
			result.Truncated = event.Delta.StopReason == anthropicStopMaxToken
			result.OutputTokens = event.Usage.OutputTokens
			if event.Usage.InputTokens > 0 {
				result.InputTokens = event.Usage.InputTokens
			}
		case anthropicMessageStop:
			completed = true
		case anthropicError:
			return &APIError{Provider: constants.ProviderAnthropic, Type: event.Error.Type, Message: event.Error.Message}
		}

		return nil
	})
	if err != nil {
		return result, err
	}
	if !completed {
		return result, failure.ErrWithTrace(errStreamEnded)
	}

	return result, nil
}

func (a *Anthropic) handleBlockStart(event anthropicEvent, onSearch func(SearchEvent) error) error {
	if onSearch == nil {
		return nil
	}

	switch event.ContentBlock.Type {
	case anthropicServerToolUse:
		if event.ContentBlock.Name == searchToolName {
			return onSearch(SearchEvent{Type: "started"})
		}
	case anthropicSearchResult:
		sources := a.parseSearchResults(event.ContentBlock.Content)
		for _, src := range sources {
			if err := onSearch(SearchEvent{Type: "source", URL: src.URL, Title: src.Title}); err != nil {
				return err
			}
		}
		return onSearch(SearchEvent{Type: "completed"})
	}

	return nil
}

func (a *Anthropic) parseSearchResults(raw json.RawMessage) []searchResultItem {
	if raw == nil {
		return nil
	}
	var items []searchResultItem
	_ = json.Unmarshal(raw, &items)
	var results []searchResultItem
	for _, item := range items {
		if item.Type == anthropicSearchResultRow && item.URL != "" {
			results = append(results, item)
		}
	}

	return results
}
