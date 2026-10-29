package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"chatsmith-mini/configs"
)

const classifierMaxTokens = 8

type Classifier struct {
	baseURL string
	apiKey  string
	model   string
	prompt  string
	timeout time.Duration
	client  *http.Client
}

func NewClassifier(cfg configs.ProviderConfig, chatCfg configs.ChatConfig, client *http.Client) *Classifier {
	return &Classifier{
		baseURL: cfg.AnthropicBaseURL,
		apiKey:  cfg.AnthropicAPIKey,
		model:   chatCfg.ClassifierModel,
		prompt:  chatCfg.ClassifierPrompt,
		timeout: chatCfg.ClassifierTimeout,
		client:  client,
	}
}

type classifyBody struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []Message `json:"messages"`
}

type classifyResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (c *Classifier) NeedsSearch(ctx context.Context, lastMessage string) bool {
	classifyCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	body := classifyBody{
		Model:     c.model,
		MaxTokens: classifierMaxTokens,
		System:    c.prompt,
		Messages:  []Message{{Role: "user", Content: lastMessage}},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		slog.Warn("classifier_marshal_failed", "error", err)
		return true
	}

	req, err := http.NewRequestWithContext(classifyCtx, http.MethodPost, c.baseURL+anthropicMessagesPath, bytes.NewReader(payload))
	if err != nil {
		slog.Warn("classifier_request_failed", "error", err)
		return true
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set(headerAnthropicKey, c.apiKey)
	req.Header.Set(headerAnthropicVer, anthropicVersion)

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Warn("classifier_call_failed", "error", err)
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("classifier_status_error", "status", resp.StatusCode)
		return true
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		slog.Warn("classifier_read_failed", "error", err)
		return true
	}

	var result classifyResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		slog.Warn("classifier_parse_failed", "error", err)
		return true
	}

	if len(result.Content) == 0 {
		return true
	}

	return result.Content[0].Text == "search"
}
