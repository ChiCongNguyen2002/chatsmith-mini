package configs

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"

	"chatsmith-mini/pkg/failure"

	"github.com/spf13/viper"
)

const (
	envConfigFile = "CONFIG_FILE"
	keyDelimiter  = "."
	envDelimiter  = "_"
)

var explicitEnv = map[string]string{
	KeyServerPort:             "PORT",
	KeyDatabaseURL:            "DATABASE_URL",
	KeyProviderOpenAIKey:      "OPENAI_API_KEY",
	KeyProviderAnthropicKey:   "ANTHROPIC_API_KEY",
	KeyAuthGoogleClientID:     "AUTH_GOOGLE_CLIENT_ID",
	KeyAuthGoogleClientSecret: "AUTH_GOOGLE_CLIENT_SECRET",
	KeyAuthGitHubClientID:     "AUTH_GITHUB_CLIENT_ID",
	KeyAuthGitHubClientSecret: "AUTH_GITHUB_CLIENT_SECRET",
	KeyAuthFacebookAppID:      "AUTH_FACEBOOK_APP_ID",
	KeyAuthFacebookAppSecret:  "AUTH_FACEBOOK_APP_SECRET",
	KeyAuthPublicBaseURL:      "AUTH_PUBLIC_BASE_URL",
	KeyAuthSessionSecret:      "AUTH_SESSION_SECRET",
}

func Load(ctx context.Context) (*Config, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	if err := readFile(v); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(keyDelimiter, envDelimiter))
	v.AutomaticEnv()
	for key, name := range explicitEnv {
		if err := v.BindEnv(key, name); err != nil {
			return nil, failure.ErrWithTrace(err)
		}
	}
	if err := loadSecrets(ctx, v); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func readFile(v *viper.Viper) error {
	path := os.Getenv(envConfigFile)
	explicit := path != ""
	if !explicit {
		path = DefaultConfigFile
	}

	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil && (explicit || !errors.Is(err, fs.ErrNotExist)) {
		return failure.ErrWithTrace(err)
	}

	return nil
}
