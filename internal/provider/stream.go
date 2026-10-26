package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"
)

const (
	maxErrorBody          = 4096
	maxIdleConnsPerHost   = 256
	idleConnTimeout       = 90 * time.Second
	tlsHandshakeTimeout   = 10 * time.Second
	responseHeaderTimeout = 30 * time.Second
	headerContentType     = "Content-Type"
	contentTypeJSON       = "application/json"
)

var errStreamEnded = errors.New("stream ended before completion")

func NewHTTPClient() *http.Client {
	return &http.Client{Transport: &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		IdleConnTimeout:       idleConnTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
	}}
}

type target struct {
	provider constants.Provider
	url      string
	headers  map[string]string
	body     any
}

func streamSSE(ctx context.Context, client *http.Client, t target, handle sseHandler) error {
	payload, err := json.Marshal(t.body)
	if err != nil {
		return failure.ErrWithTrace(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url, bytes.NewReader(payload))
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	for name, value := range t.headers {
		req.Header.Set(name, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return statusError(t.provider, resp)
	}
	if err := readSSE(resp.Body, handle); err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}

func statusError(provider constants.Provider, resp *http.Response) error {
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))

	var payload struct {
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.Unmarshal(raw, &payload)

	return &APIError{Provider: provider, Status: resp.StatusCode, Type: payload.Error.Type, Message: payload.Error.Message}
}
