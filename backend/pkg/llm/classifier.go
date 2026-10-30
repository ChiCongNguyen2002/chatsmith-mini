package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/logger"
)

const (
	classifierMaxTokens = 8
	classifierMaxBody   = 64 << 10
	classifierHeadBytes = 200
	verdictSearch       = "search"
	logClassifierFailed = "classifier_failed"
	contentTypeText     = "text"
	thinkingDisabled    = "disabled"
)

var errNoVerdict = errors.New("classifier returned no text block")

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
	Model     string        `json:"model"`
	MaxTokens int           `json:"max_tokens"`
	System    string        `json:"system"`
	Messages  []Message     `json:"messages"`
	Thinking  thinkingParam `json:"thinking"`
}

type thinkingParam struct {
	Type string `json:"type"`
}

type classifyResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (c *Classifier) NeedsSearch(ctx context.Context, lastMessage string) bool {
	classifyCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	verdict, err := c.verdict(classifyCtx, lastMessage)
	if err != nil {
		logger.FromContext(ctx).Warn(logClassifierFailed, constants.LogKeyError, err)

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

	raw, err := io.ReadAll(io.LimitReader(resp.Body, classifierMaxBody))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", bodyError(resp.StatusCode, raw, errStatus)
	}

	var result classifyResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", bodyError(resp.StatusCode, raw, err)
	}
	for _, block := range result.Content {
		if block.Type == contentTypeText {
			return block.Text, nil
		}
	}

	return "", bodyError(resp.StatusCode, raw, errNoVerdict)
}

var errStatus = errors.New("unexpected status")

func bodyError(status int, raw []byte, cause error) error {
	head := raw
	if len(head) > classifierHeadBytes {
		head = head[:classifierHeadBytes]
	}

	return fmt.Errorf("status=%d body_len=%d body_head=%q: %w", status, len(raw), head, cause)
}

func (c *Classifier) newRequest(ctx context.Context, lastMessage string) (*http.Request, error) {
	payload, err := json.Marshal(classifyBody{
		Model:     c.model,
		MaxTokens: classifierMaxTokens,
		System:    c.prompt,
		Messages:  []Message{{Role: "user", Content: lastMessage}},
		Thinking:  thinkingParam{Type: thinkingDisabled},
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
