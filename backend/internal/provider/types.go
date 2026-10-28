package provider

import (
	"fmt"

	"chatsmith-mini/constants"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SearchEvent struct {
	Type  string
	Query string
	URL   string
	Title string
}

type Request struct {
	Model            string
	Messages         []Message
	MaxOutputTokens  int
	SystemPrompt     string
	WebSearchEnabled bool
	OnSearch         func(SearchEvent) error
}

type Result struct {
	InputTokens       int64
	OutputTokens      int64
	ProviderRequestID string
	Truncated         bool
}

type APIError struct {
	Provider constants.Provider
	Status   int
	Type     string
	Message  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s api error status=%d type=%s message=%s", e.Provider, e.Status, e.Type, e.Message)
}
