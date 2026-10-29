package auth

import (
	"context"
	"strings"
	"time"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/user"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/oauth"
	"chatsmith-mini/pkg/sessionstate"
)

type Service struct {
	providers map[constants.AuthProvider]oauth.Provider
	users     *user.Service
	repo      *Repo
	codec     *sessionstate.Codec
	tokens    *sessionstate.Tokens
	cfg       configs.AuthConfig
}

func newService(repo *Repo, users *user.Service, providers []oauth.Provider, cfg configs.AuthConfig) *Service {
	byName := make(map[constants.AuthProvider]oauth.Provider, len(providers))
	for _, p := range providers {
		byName[p.Name()] = p
	}

	return &Service{
		providers: byName,
		users:     users,
		repo:      repo,
		codec:     sessionstate.NewCodec(cfg.SessionSecret),
		tokens:    sessionstate.NewTokens(),
		cfg:       cfg,
	}
}

func (s *Service) StartLogin(name constants.AuthProvider) (string, string, error) {
	provider, err := s.provider(name)
	if err != nil {
		return "", "", err
	}

	nonce, err := s.tokens.NewToken()
	if err != nil {
		return "", "", failure.ErrWithTrace(err)
	}
	verifier, err := s.tokens.NewToken()
	if err != nil {
		return "", "", failure.ErrWithTrace(err)
	}

	state := string(name) + stateSep + nonce + stateSep + verifier

	return provider.AuthURL(nonce, verifier), s.codec.Seal(state, s.cfg.StateTTL), nil
}

func (s *Service) FinishLogin(ctx context.Context, name constants.AuthProvider, stateToken, state, code string) (string, error) {
	provider, err := s.provider(name)
	if err != nil {
		return "", err
	}

	opened, err := s.codec.Open(stateToken)
	if err != nil {
		return "", failure.New(failure.CodeUnauthorized)
	}
	parts := strings.Split(opened, stateSep)
	if len(parts) != stateParts || parts[0] != string(name) || parts[1] != state {
		return "", failure.New(failure.CodeUnauthorized)
	}

	profile, err := provider.Identify(ctx, code, parts[2])
	if err != nil {
		return "", err
	}

	account, err := s.users.Login(ctx, user.Identity{
		Provider:      name,
		Subject:       profile.Subject,
		Email:         profile.Email,
		EmailVerified: profile.EmailVerified,
		Name:          profile.Name,
		Picture:       profile.Picture,
	})
	if err != nil {
		return "", err
	}

	return s.openSession(ctx, account.ID)
}

func (s *Service) VerifySession(ctx context.Context, token string) (string, error) {
	session, err := s.repo.Find(ctx, HashToken(token))
	if err != nil {
		return "", err
	}
	if session == nil || session.Expired(time.Now()) {
		return "", failure.New(failure.CodeUnauthorized)
	}

	return session.UserID, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.repo.Delete(ctx, HashToken(token))
}

func (s *Service) Me(ctx context.Context, userID string) (*user.User, error) {
	return s.users.Get(ctx, userID)
}

func (s *Service) openSession(ctx context.Context, userID string) (string, error) {
	token, err := s.tokens.NewToken()
	if err != nil {
		return "", failure.ErrWithTrace(err)
	}

	err = s.repo.Create(ctx, Session{
		TokenHash: HashToken(token),
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.cfg.SessionTTL),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) provider(name constants.AuthProvider) (oauth.Provider, error) {
	provider, ok := s.providers[name]
	if !ok {
		return nil, failure.New(failure.CodeNotFound)
	}

	return provider, nil
}
