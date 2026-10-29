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

	return providers
}

func callbackURL(cfg configs.AuthConfig, name constants.AuthProvider) string {
	return cfg.PublicBaseURL + constants.AuthCallbackPath(name)
}
