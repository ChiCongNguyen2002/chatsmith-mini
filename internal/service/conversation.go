package service

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"chatsmith-mini/internal/repository"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/uid"
)

const (
	maxTitleLen         = 200
	defaultListLimit    = 20
	maxListLimit        = 100
	defaultMessageLimit = 50
	maxMessageLimit     = 200
)

type ConversationPage struct {
	Items      []repository.ConversationItem `json:"items"`
	NextCursor *string                       `json:"next_cursor"`
}

type Conversation struct {
	repo *repository.Conversation
}

func NewConversation(repo *repository.Conversation) *Conversation {
	return &Conversation{repo: repo}
}

func (s *Conversation) Create(ctx context.Context, userID, title string) (*repository.ConversationItem, error) {
	title = strings.TrimSpace(title)
	if utf8.RuneCountInString(title) > maxTitleLen {
		return nil, failure.ErrWithTrace(failure.ErrBadRequest)
	}

	return s.repo.Create(ctx, userID, title)
}

func (s *Conversation) List(ctx context.Context, userID, rawCursor string, limit int) (*ConversationPage, error) {
	limit = clamp(limit, defaultListLimit, maxListLimit)

	var (
		afterCreatedAt *time.Time
		afterID        string
	)
	if rawCursor != "" {
		createdAt, id, err := decodeCursor(rawCursor)
		if err != nil {
			return nil, err
		}
		afterCreatedAt, afterID = &createdAt, id
	}

	items, err := s.repo.List(ctx, userID, limit+1, afterCreatedAt, afterID)
	if err != nil {
		return nil, err
	}

	page := &ConversationPage{Items: items}
	if len(items) > limit {
		page.Items = items[:limit]
		last := page.Items[limit-1]
		cursor := encodeCursor(last.CreatedAt, last.ID)
		page.NextCursor = &cursor
	}

	return page, nil
}

func (s *Conversation) Messages(ctx context.Context, userID, conversationID string, afterSeq int64, limit int) ([]repository.MessageItem, error) {
	if !uid.Valid(conversationID) {
		return nil, failure.ErrWithTrace(failure.ErrNotFound)
	}

	owned, err := s.repo.Owned(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, failure.ErrWithTrace(failure.ErrNotFound)
	}

	return s.repo.Messages(ctx, conversationID, afterSeq, clamp(limit, defaultMessageLimit, maxMessageLimit))
}

func clamp(limit, fallback, ceiling int) int {
	if limit <= 0 {
		return fallback
	}
	if limit > ceiling {
		return ceiling
	}

	return limit
}

func encodeCursor(createdAt time.Time, id string) string {
	raw := strconv.FormatInt(createdAt.UnixMicro(), 10) + "_" + id

	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(cursor string) (time.Time, string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, "", failure.ErrWithTrace(failure.ErrInvalidCursor)
	}

	micro, id, found := strings.Cut(string(raw), "_")
	if !found || !uid.Valid(id) {
		return time.Time{}, "", failure.ErrWithTrace(failure.ErrInvalidCursor)
	}

	value, err := strconv.ParseInt(micro, 10, 64)
	if err != nil {
		return time.Time{}, "", failure.ErrWithTrace(failure.ErrInvalidCursor)
	}

	return time.UnixMicro(value).UTC(), id, nil
}
