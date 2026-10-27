package user

import "context"

type Service struct {
	repo *Repo
}

func newService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(ctx context.Context, profile Profile) (*User, error) {
	return s.repo.Upsert(ctx, profile)
}

func (s *Service) Get(ctx context.Context, userID string) (*User, error) {
	return s.repo.Get(ctx, userID)
}
