package user

import (
	"context"
)

type Service struct {
	repo *Repo
}

func newService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID string) (*User, error) {
	return s.repo.Get(ctx, userID)
}

func (s *Service) Login(ctx context.Context, identity Identity) (*User, error) {
	found, err := s.repo.FindByIdentity(ctx, string(identity.Provider), identity.Subject)
	if err != nil {
		return nil, err
	}
	if found != nil {
		if err := s.repo.RefreshProfile(ctx, found.ID, identity); err != nil {
			return nil, err
		}

		return found, nil
	}

	if shouldLink(identity) {
		match, err := s.repo.FindByVerifiedEmail(ctx, identity.Email)
		if err != nil {
			return nil, err
		}
		if match != nil {
			if err := s.repo.AddIdentity(ctx, match.ID, identity); err != nil {
				return nil, err
			}

			return match, nil
		}
	}

	return s.repo.CreateWithIdentity(ctx, identity)
}

func shouldLink(identity Identity) bool {
	return identity.EmailVerified && identity.Email != ""
}
