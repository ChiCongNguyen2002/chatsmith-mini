package auth

import (
	"context"

	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/user"
	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/logger"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type Google struct {
	cfg oauth2.Config
}

func newGoogle(cfg configs.AuthConfig) *Google {
	return &Google{cfg: oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.PublicBaseURL + PathGoogleCallback,
		Scopes:       []string{scopeOpenID, scopeEmail, scopeProfile},
		Endpoint:     oauth2.Endpoint{AuthURL: googleAuthURL, TokenURL: googleTokenURL},
	}}
}

func (g *Google) AuthURL(state string) string {
	return g.cfg.AuthCodeURL(state)
}

func (g *Google) Profile(ctx context.Context, code string) (user.Profile, error) {
	log := logger.FromContext(ctx)
	token, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		log.Warn(logGoogleExchangeFailed, constants.LogKeyError, err)
		return user.Profile{}, failure.New(failure.CodeUnauthorized)
	}

	raw, _ := token.Extra(extraIDToken).(string)
	if raw == "" {
		log.Warn(logGoogleIDTokenMissing)
		return user.Profile{}, failure.New(failure.CodeUnauthorized)
	}
	payload, err := idtoken.Validate(ctx, raw, g.cfg.ClientID)
	if err != nil {
		log.Warn(logGoogleIDTokenInvalid, constants.LogKeyError, err)
		return user.Profile{}, failure.New(failure.CodeUnauthorized)
	}
	if verified, _ := payload.Claims[claimEmailVerified].(bool); !verified {
		log.Warn(logGoogleEmailUnverified)
		return user.Profile{}, failure.New(failure.CodeUnauthorized)
	}

	return user.Profile{
		GoogleSub: payload.Subject,
		Email:     claimString(payload.Claims, claimEmail),
		Name:      claimString(payload.Claims, claimName),
		Picture:   claimString(payload.Claims, claimPicture),
	}, nil
}

func claimString(claims map[string]any, key string) string {
	value, _ := claims[key].(string)

	return value
}
