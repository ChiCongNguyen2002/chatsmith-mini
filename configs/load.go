package configs

import (
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

var legacyEnv = map[string]string{
	KeyServerPort:           "PORT",
	KeyDatabaseURL:          "DATABASE_URL",
	KeyProviderOpenAIKey:    "OPENAI_API_KEY",
	KeyProviderAnthropicKey: "ANTHROPIC_API_KEY",
}

func Load() (*Config, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	if err := readFile(v); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(keyDelimiter, envDelimiter))
	v.AutomaticEnv()
	for key, name := range legacyEnv {
		if err := v.BindEnv(key, name); err != nil {
			return nil, failure.ErrWithTrace(err)
		}
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
