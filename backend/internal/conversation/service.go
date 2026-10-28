package conversation

import (
	"context"
	"strings"
	"unicode/utf8"

	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/uid"
)

type Service struct {
	repo *Repo
}

func newService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID, title string) (*Conversation, error) {
	title = strings.TrimSpace(title)
	if utf8.RuneCountInString(title) > MaxTitleRunes {
		return nil, failure.New(failure.CodeBadRequest)
	}

	return s.repo.Create(ctx, userID, title)
}

func (s *Service) Get(ctx context.Context, userID, conversationID string) (*Conversation, error) {
	if !uid.Valid(conversationID) {
		return nil, failure.New(failure.CodeNotFound)
	}

	return s.repo.Get(ctx, conversationID, userID)
}

func (s *Service) List(ctx context.Context, userID, rawCursor string, limit int) (*ConversationPage, error) {
	limit = clamp(limit, DefaultListLimit, MaxListLimit)

	var after *cursor
	if rawCursor != "" {
		decoded, err := decodeCursor(rawCursor)
		if err != nil {
			return nil, err
		}
		after = decoded
	}

	items, err := s.repo.List(ctx, userID, limit+1, after)
	if err != nil {
		return nil, err
	}

	page := &ConversationPage{Items: items}
	if len(items) > limit {
		page.Items = items[:limit]
		next := encodeCursor(page.Items[limit-1])
		page.NextCursor = &next
	}

	return page, nil
}

func (s *Service) Delete(ctx context.Context, userID, conversationID string) error {
	if !uid.Valid(conversationID) {
		return failure.New(failure.CodeNotFound)
	}

	return s.repo.Delete(ctx, conversationID, userID)
}

func (s *Service) Rename(ctx context.Context, userID, conversationID, title string) error {
	title = strings.TrimSpace(title)
	if title == "" || utf8.RuneCountInString(title) > MaxTitleRunes {
		return failure.New(failure.CodeBadRequest)
	}
	if !uid.Valid(conversationID) {
		return failure.New(failure.CodeNotFound)
	}

	return s.repo.UpdateTitle(ctx, conversationID, userID, title)
}

func (s *Service) Messages(ctx context.Context, userID, conversationID string, afterSeq int64, limit int) ([]Message, error) {
	if _, err := s.Get(ctx, userID, conversationID); err != nil {
		return nil, err
	}

	return s.repo.Messages(ctx, conversationID, afterSeq, clamp(limit, DefaultMessageLimit, MaxMessageLimit))
}

func clamp(limit, fallback, ceiling int) int {
	if limit <= 0 {
		return fallback
	}

	return min(limit, ceiling)
}
