package auth

import (
	"context"
	"crypto/rand"

	"chatsmith-mini/configs"
	"chatsmith-mini/internal/user"
	"chatsmith-mini/pkg/failure"
)

type Service struct {
	google *Google
	codec  *Codec
	users  *user.Service
	cfg    configs.AuthConfig
}

func (s *Service) StartLogin() (string, string) {
	nonce := rand.Text()

	return s.google.AuthURL(nonce), s.codec.Seal(nonce, s.cfg.StateTTL)
}

func (s *Service) FinishLogin(ctx context.Context, stateToken, state, code string) (string, error) {
	nonce, err := s.codec.Open(stateToken)
	if err != nil || nonce != state {
		return "", failure.New(failure.CodeUnauthorized)
	}

	profile, err := s.google.Profile(ctx, code)
	if err != nil {
		return "", err
	}

	account, err := s.users.Login(ctx, profile)
	if err != nil {
		return "", err
	}

	return s.codec.Seal(account.ID, s.cfg.SessionTTL), nil
}

func (s *Service) VerifySession(token string) (string, error) {
	userID, err := s.codec.Open(token)
	if err != nil {
		return "", failure.New(failure.CodeUnauthorized)
	}

	return userID, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*user.User, error) {
	return s.users.Get(ctx, userID)
}
