package conversation

import (
	"context"
	"strings"
	"unicode/utf8"

	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/pagination"
	"chatsmith-mini/pkg/uid"
)

type Service struct {
	repo *Repo
}

func newService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, userID, title string) (*Conversation, error) {
	normalized, err := createTitle(title)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, userID, normalized)
}

func (s *Service) Get(ctx context.Context, userID, conversationID string) (*Conversation, error) {
	if !uid.Valid(conversationID) {
		return nil, failure.New(failure.CodeNotFound)
	}

	return s.repo.Get(ctx, conversationID, userID)
}

func (s *Service) List(ctx context.Context, userID, rawCursor string, limit int) (*Page, error) {
	limit = clamp(limit, DefaultListLimit, MaxListLimit)

	var after *pagination.Keyset
	if rawCursor != "" {
		keyset, err := pagination.Decode(rawCursor)
		if err != nil {
			return nil, failure.New(failure.CodeBadRequest)
		}
		after = &keyset
	}

	items, err := s.repo.List(ctx, userID, limit+1, after)
	if err != nil {
		return nil, err
	}

	page := &Page{Items: items}
	if len(items) > limit {
		last := items[limit-1]
		page.Items = items[:limit]
		next := pagination.Encode(last.UpdatedAt.UnixMicro(), last.ID)
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
	normalized := strings.TrimSpace(title)
	if normalized == "" || utf8.RuneCountInString(normalized) > MaxTitleRunes {
		return failure.New(failure.CodeBadRequest)
	}
	if !uid.Valid(conversationID) {
		return failure.New(failure.CodeNotFound)
	}

	return s.repo.UpdateTitle(ctx, conversationID, userID, normalized)
}

func createTitle(raw string) (string, error) {
	title := strings.TrimSpace(raw)
	if utf8.RuneCountInString(title) > MaxTitleRunes {
		return "", failure.New(failure.CodeBadRequest)
	}

	return title, nil
}

func clamp(limit, fallback, ceiling int) int {
	if limit <= 0 {
		return fallback
	}

	return min(limit, ceiling)
}
