package oauth

import (
	"context"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

const (
	googleAuthURL   = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL  = "https://oauth2.googleapis.com/token"
	extraIDToken    = "id_token"
	claimEmail      = "email"
	claimEmailVerif = "email_verified"
	claimName       = "name"
	claimPicture    = "picture"
)

type Google struct {
	cfg oauth2.Config
}

func NewGoogle(clientID, clientSecret, redirectURL string) *Google {
	return &Google{cfg: oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     oauth2.Endpoint{AuthURL: googleAuthURL, TokenURL: googleTokenURL},
	}}
}

func (g *Google) Name() constants.AuthProvider {
	return constants.AuthProviderGoogle
}

func (g *Google) AuthURL(state, verifier string) string {
	return g.cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
}

func (g *Google) Identify(ctx context.Context, code, verifier string) (Profile, error) {
	token, err := g.cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return Profile{}, failure.New(failure.CodeUnauthorized)
	}

	raw, _ := token.Extra(extraIDToken).(string)
	if raw == "" {
		return Profile{}, failure.New(failure.CodeUnauthorized)
	}
	payload, err := idtoken.Validate(ctx, raw, g.cfg.ClientID)
	if err != nil {
		return Profile{}, failure.New(failure.CodeUnauthorized)
	}
	if verified, _ := payload.Claims[claimEmailVerif].(bool); !verified {
		return Profile{}, failure.New(failure.CodeUnauthorized)
	}

	return Profile{
		Subject:       payload.Subject,
		Email:         claimString(payload.Claims, claimEmail),
		EmailVerified: true,
		Name:          claimString(payload.Claims, claimName),
		Picture:       claimString(payload.Claims, claimPicture),
	}, nil
}

func claimString(claims map[string]any, key string) string {
	value, _ := claims[key].(string)

	return value
}
