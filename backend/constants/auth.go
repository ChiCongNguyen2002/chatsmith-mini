package constants

import "strings"

type AuthProvider string

const (
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderGitHub   AuthProvider = "github"
	AuthProviderFacebook AuthProvider = "facebook"
)

func AuthCallbackPath(name AuthProvider) string {
	return strings.Replace(PathAuthCallback, ":"+ParamProvider, string(name), 1)
}
