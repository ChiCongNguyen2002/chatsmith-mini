package oauth

import (
	"context"
	"strconv"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"

	"golang.org/x/oauth2"
)

const (
	githubAuthURL   = "https://github.com/login/oauth/authorize"
	githubTokenURL  = "https://github.com/login/oauth/access_token"
	githubUserURL   = "https://api.github.com/user"
	githubEmailsURL = "https://api.github.com/user/emails"
)

type GitHub struct {
	cfg oauth2.Config
}

type githubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func NewGitHub(clientID, clientSecret, redirectURL string) *GitHub {
	return &GitHub{cfg: oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     oauth2.Endpoint{AuthURL: githubAuthURL, TokenURL: githubTokenURL},
	}}
}

func (g *GitHub) Name() constants.AuthProvider {
	return constants.AuthProviderGitHub
}

func (g *GitHub) AuthURL(state, verifier string) string {
	return g.cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
}

func (g *GitHub) Identify(ctx context.Context, code, verifier string) (Profile, error) {
	token, err := g.cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return Profile{}, failure.New(failure.CodeUnauthorized)
	}

	client := g.cfg.Client(ctx, token)
	var account githubUser
	if err := getJSON(ctx, client, githubUserURL, &account); err != nil {
		return Profile{}, err
	}
	var emails []githubEmail
	if err := getJSON(ctx, client, githubEmailsURL, &emails); err != nil {
		return Profile{}, err
	}

	email, verified := primaryVerifiedEmail(emails)
	name := account.Name
	if name == "" {
		name = account.Login
	}

	return Profile{
		Subject:       strconv.FormatInt(account.ID, 10),
		Email:         email,
		EmailVerified: verified,
		Name:          name,
		Picture:       account.AvatarURL,
	}, nil
}

func primaryVerifiedEmail(emails []githubEmail) (string, bool) {
	for _, item := range emails {
		if item.Primary && item.Verified {
			return item.Email, true
		}
	}

	return "", false
}
