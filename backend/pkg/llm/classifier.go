package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"chatsmith-mini/configs"
)

const (
	classifierMaxTokens = 8
	classifierMaxBody   = 1024
	verdictSearch       = "search"
	logClassifierFailed = "classifier_failed"
)

var errNoVerdict = errors.New("classifier returned no content")

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

	verdict, err := c.verdict(classifyCtx, lastMessage)
	if err != nil {
		slog.Warn(logClassifierFailed, "error", err)

		return true
	}

	return verdict == verdictSearch
}

func (c *Classifier) verdict(ctx context.Context, lastMessage string) (string, error) {
	req, err := c.newRequest(ctx, lastMessage)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, classifierMaxBody))
	if err != nil {
		return "", err
	}

	var result classifyResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", err
	}
	if len(result.Content) == 0 {
		return "", errNoVerdict
	}

	return result.Content[0].Text, nil
}

func (c *Classifier) newRequest(ctx context.Context, lastMessage string) (*http.Request, error) {
	payload, err := json.Marshal(classifyBody{
		Model:     c.model,
		MaxTokens: classifierMaxTokens,
		System:    c.prompt,
		Messages:  []Message{{Role: "user", Content: lastMessage}},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+anthropicMessagesPath, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set(headerAnthropicKey, c.apiKey)
	req.Header.Set(headerAnthropicVer, anthropicVersion)

	return req, nil
}
