package configs

import (
	"fmt"
	"math"
	"time"

	"chatsmith-mini/pkg/failure"
)

const (
	maxPort        = 65535
	microPerUSD    = 1_000_000
	minSecretBytes = 32
)

type Config struct {
	Server   ServerConfig   `mapstructure:"SERVER"`
	Database DatabaseConfig `mapstructure:"DATABASE"`
	Provider ProviderConfig `mapstructure:"PROVIDER"`
	Chat     ChatConfig     `mapstructure:"CHAT"`
	Auth     AuthConfig     `mapstructure:"AUTH"`
}

type AuthConfig struct {
	GoogleClientID     string        `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string        `mapstructure:"GOOGLE_CLIENT_SECRET"`
	PublicBaseURL      string        `mapstructure:"PUBLIC_BASE_URL"`
	SessionSecret      string        `mapstructure:"SESSION_SECRET"`
	SessionTTL         time.Duration `mapstructure:"SESSION_TTL"`
	StateTTL           time.Duration `mapstructure:"STATE_TTL"`
	CookieSecure       bool          `mapstructure:"COOKIE_SECURE"`
}

type ServerConfig struct {
	Port              int           `mapstructure:"PORT"`
	ReadHeaderTimeout time.Duration `mapstructure:"READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `mapstructure:"READ_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"IDLE_TIMEOUT"`
	MaxHeaderBytes    int           `mapstructure:"MAX_HEADER_BYTES"`
	ShutdownTimeout   time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
	RequestTimeout    time.Duration `mapstructure:"REQUEST_TIMEOUT"`
	StreamGrace       time.Duration `mapstructure:"STREAM_GRACE"`
	FrontendDir       string        `mapstructure:"FRONTEND_DIR"`
}

type DatabaseConfig struct {
	URL          string `mapstructure:"URL"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
}

type ProviderConfig struct {
	OpenAIAPIKey          string `mapstructure:"OPENAI_API_KEY"`
	OpenAIBaseURL         string `mapstructure:"OPENAI_BASE_URL"`
	OpenAIReasoningEffort string `mapstructure:"OPENAI_REASONING_EFFORT"`
	AnthropicAPIKey       string `mapstructure:"ANTHROPIC_API_KEY"`
	AnthropicBaseURL      string `mapstructure:"ANTHROPIC_BASE_URL"`
	AnthropicEffort       string `mapstructure:"ANTHROPIC_EFFORT"`
}

type ChatConfig struct {
	DefaultModel      string        `mapstructure:"DEFAULT_MODEL"`
	MaxMessageChars   int           `mapstructure:"MAX_MESSAGE_CHARS"`
	MaxBodyBytes      int64         `mapstructure:"MAX_BODY_BYTES"`
	MaxOutputTokens   int           `mapstructure:"MAX_OUTPUT_TOKENS"`
	ContextTokens     int           `mapstructure:"CONTEXT_TOKENS"`
	Timeout           time.Duration `mapstructure:"TIMEOUT"`
	DailyBudgetUSD    float64       `mapstructure:"DAILY_BUDGET_USD"`
	SystemPrompt      string        `mapstructure:"SYSTEM_PROMPT"`
	WebSearch         bool          `mapstructure:"WEB_SEARCH"`
	ClassifierModel   string        `mapstructure:"CLASSIFIER_MODEL"`
	ClassifierPrompt  string        `mapstructure:"CLASSIFIER_PROMPT"`
	ClassifierTimeout time.Duration `mapstructure:"CLASSIFIER_TIMEOUT"`
}

func (c ChatConfig) DailyBudgetMicro() int64 {
	return int64(math.Round(c.DailyBudgetUSD * microPerUSD))
}

type rule struct {
	key string
	ok  bool
}

func (c *Config) validate() error {
	s, d, p, h := c.Server, c.Database, c.Provider, c.Chat
	rules := []rule{
		{KeyServerPort, s.Port > 0 && s.Port <= maxPort},
		{KeyServerReadHeaderTimeout, s.ReadHeaderTimeout > 0},
		{KeyServerReadTimeout, s.ReadTimeout > 0},
		{KeyServerIdleTimeout, s.IdleTimeout > 0},
		{KeyServerMaxHeaderBytes, s.MaxHeaderBytes > 0},
		{KeyServerShutdownTimeout, s.ShutdownTimeout > 0},
		{KeyServerRequestTimeout, s.RequestTimeout > 0},
		{KeyServerStreamGrace, s.StreamGrace >= 0},
		{KeyDatabaseURL, d.URL != ""},
		{KeyDatabaseMaxOpenConns, d.MaxOpenConns > 0},
		{KeyDatabaseMaxIdleConns, d.MaxIdleConns > 0 && d.MaxIdleConns <= d.MaxOpenConns},
		{KeyProviderOpenAIBaseURL, p.OpenAIBaseURL != ""},
		{KeyProviderAnthropicBaseURL, p.AnthropicBaseURL != ""},
		{KeyChatDefaultModel, h.DefaultModel != ""},
		{KeyChatMaxMessageChars, h.MaxMessageChars > 0},
		{KeyChatMaxBodyBytes, h.MaxBodyBytes > 0},
		{KeyChatMaxOutputTokens, h.MaxOutputTokens > 0},
		{KeyChatContextTokens, h.ContextTokens > 0},
		{KeyChatTimeout, h.Timeout > 0},
		{KeyChatDailyBudgetUSD, h.DailyBudgetUSD >= 0},
		{KeyAuthProviders, c.Auth.HasProvider()},
		{KeyAuthPublicBaseURL, c.Auth.PublicBaseURL != ""},
		{KeyAuthSessionSecret, len(c.Auth.SessionSecret) >= minSecretBytes},
		{KeyAuthSessionTTL, c.Auth.SessionTTL > 0},
		{KeyAuthStateTTL, c.Auth.StateTTL > 0},
	}
	for _, r := range rules {
		if !r.ok {
			return failure.ErrWithTrace(fmt.Errorf("%s is invalid", r.key))
		}
	}

	return nil
}

func (a AuthConfig) HasProvider() bool {
	return a.GoogleClientID != "" && a.GoogleClientSecret != ""
}
