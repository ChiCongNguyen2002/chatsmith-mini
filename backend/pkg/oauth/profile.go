package oauth

import (
	"context"

	"chatsmith-mini/constants"
)

type Profile struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

type Provider interface {
	Name() constants.AuthProvider
	AuthURL(state, verifier string) string
	Identify(ctx context.Context, code, verifier string) (Profile, error)
}
