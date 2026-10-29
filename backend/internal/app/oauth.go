package app

import (
	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/oauth"
)

func newAuthProviders(cfg configs.AuthConfig) []oauth.Provider {
	var providers []oauth.Provider
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		providers = append(providers, oauth.NewGoogle(cfg.GoogleClientID, cfg.GoogleClientSecret, callbackURL(cfg, constants.AuthProviderGoogle)))
	}
	if cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "" {
		providers = append(providers, oauth.NewGitHub(cfg.GitHubClientID, cfg.GitHubClientSecret, callbackURL(cfg, constants.AuthProviderGitHub)))
	}
	if cfg.FacebookAppID != "" && cfg.FacebookAppSecret != "" {
		providers = append(providers, oauth.NewFacebook(cfg.FacebookAppID, cfg.FacebookAppSecret, callbackURL(cfg, constants.AuthProviderFacebook)))
	}

	return providers
}

func callbackURL(cfg configs.AuthConfig, name constants.AuthProvider) string {
	return cfg.PublicBaseURL + constants.AuthCallbackPath(name)
}
