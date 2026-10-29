package oauth

import (
	"context"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"

	"golang.org/x/oauth2"
)

const (
	facebookGraphVersion = "v23.0"
	facebookAuthURL      = "https://www.facebook.com/" + facebookGraphVersion + "/dialog/oauth"
	facebookTokenURL     = "https://graph.facebook.com/" + facebookGraphVersion + "/oauth/access_token"
	facebookMeURL        = "https://graph.facebook.com/" + facebookGraphVersion + "/me?fields=id,name,email,picture.type(large)"
)

type Facebook struct {
	cfg oauth2.Config
}

type facebookUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func NewFacebook(appID, appSecret, redirectURL string) *Facebook {
	return &Facebook{cfg: oauth2.Config{
		ClientID:     appID,
		ClientSecret: appSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     oauth2.Endpoint{AuthURL: facebookAuthURL, TokenURL: facebookTokenURL},
	}}
}

func (f *Facebook) Name() constants.AuthProvider {
	return constants.AuthProviderFacebook
}

func (f *Facebook) AuthURL(state, verifier string) string {
	return f.cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
}

func (f *Facebook) Identify(ctx context.Context, code, verifier string) (Profile, error) {
	token, err := f.cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return Profile{}, failure.New(failure.CodeUnauthorized)
	}

	var account facebookUser
	if err := getJSON(ctx, f.cfg.Client(ctx, token), facebookMeURL, &account); err != nil {
		return Profile{}, err
	}

	return Profile{
		Subject: account.ID,
		Email:   account.Email,
		Name:    account.Name,
		Picture: account.Picture.Data.URL,
	}, nil
}
