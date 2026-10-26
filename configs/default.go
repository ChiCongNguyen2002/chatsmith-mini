package configs

import "time"

const (
	KeyServerPort               = "SERVER.PORT"
	KeyServerReadHeaderTimeout  = "SERVER.READ_HEADER_TIMEOUT"
	KeyServerReadTimeout        = "SERVER.READ_TIMEOUT"
	KeyServerIdleTimeout        = "SERVER.IDLE_TIMEOUT"
	KeyServerMaxHeaderBytes     = "SERVER.MAX_HEADER_BYTES"
	KeyServerShutdownTimeout    = "SERVER.SHUTDOWN_TIMEOUT"
	KeyServerRequestTimeout     = "SERVER.REQUEST_TIMEOUT"
	KeyServerStreamGrace        = "SERVER.STREAM_GRACE"
	KeyDatabaseURL              = "DATABASE.URL"
	KeyDatabaseMaxOpenConns     = "DATABASE.MAX_OPEN_CONNS"
	KeyDatabaseMaxIdleConns     = "DATABASE.MAX_IDLE_CONNS"
	KeyProviderOpenAIKey        = "PROVIDER.OPENAI_API_KEY"
	KeyProviderOpenAIBaseURL    = "PROVIDER.OPENAI_BASE_URL"
	KeyProviderOpenAIEffort     = "PROVIDER.OPENAI_REASONING_EFFORT"
	KeyProviderAnthropicKey     = "PROVIDER.ANTHROPIC_API_KEY"
	KeyProviderAnthropicBaseURL = "PROVIDER.ANTHROPIC_BASE_URL"
	KeyProviderAnthropicEffort  = "PROVIDER.ANTHROPIC_EFFORT"
	KeyChatDefaultModel         = "CHAT.DEFAULT_MODEL"
	KeyChatMaxMessageChars      = "CHAT.MAX_MESSAGE_CHARS"
	KeyChatMaxBodyBytes         = "CHAT.MAX_BODY_BYTES"
	KeyChatMaxOutputTokens      = "CHAT.MAX_OUTPUT_TOKENS"
	KeyChatContextTokens        = "CHAT.CONTEXT_TOKENS"
	KeyChatTimeout              = "CHAT.TIMEOUT"
	KeyChatDailyBudgetUSD       = "CHAT.DAILY_BUDGET_USD"
)

const (
	DefaultConfigFile        = "configs/dev.config.yaml"
	DefaultServerPort        = 8080
	DefaultReadHeaderTimeout = 5 * time.Second
	DefaultReadTimeout       = 15 * time.Second
	DefaultIdleTimeout       = 90 * time.Second
	DefaultMaxHeaderBytes    = 16 << 10
	DefaultShutdownTimeout   = 8 * time.Second
	DefaultRequestTimeout    = 10 * time.Second
	DefaultStreamGrace       = 2 * time.Second
	DefaultMaxOpenConns      = 25
	DefaultMaxIdleConns      = 5
	DefaultOpenAIBaseURL     = "https://api.openai.com"
	DefaultAnthropicBaseURL  = "https://api.anthropic.com"
	DefaultProviderEffort    = "low"
	DefaultChatModel         = "claude-haiku-5-5"
	DefaultMaxMessageChars   = 8000
	DefaultMaxBodyBytes      = 65536
	DefaultMaxOutputTokens   = 2048
	DefaultContextTokens     = 24000
	DefaultChatTimeout       = 120 * time.Second
	DefaultDailyBudgetUSD    = 2.0
)

var defaults = map[string]any{
	KeyServerPort:               DefaultServerPort,
	KeyServerReadHeaderTimeout:  DefaultReadHeaderTimeout,
	KeyServerReadTimeout:        DefaultReadTimeout,
	KeyServerIdleTimeout:        DefaultIdleTimeout,
	KeyServerMaxHeaderBytes:     DefaultMaxHeaderBytes,
	KeyServerShutdownTimeout:    DefaultShutdownTimeout,
	KeyServerRequestTimeout:     DefaultRequestTimeout,
	KeyServerStreamGrace:        DefaultStreamGrace,
	KeyDatabaseMaxOpenConns:     DefaultMaxOpenConns,
	KeyDatabaseMaxIdleConns:     DefaultMaxIdleConns,
	KeyProviderOpenAIBaseURL:    DefaultOpenAIBaseURL,
	KeyProviderOpenAIEffort:     DefaultProviderEffort,
	KeyProviderAnthropicBaseURL: DefaultAnthropicBaseURL,
	KeyProviderAnthropicEffort:  DefaultProviderEffort,
	KeyChatDefaultModel:         DefaultChatModel,
	KeyChatMaxMessageChars:      DefaultMaxMessageChars,
	KeyChatMaxBodyBytes:         DefaultMaxBodyBytes,
	KeyChatMaxOutputTokens:      DefaultMaxOutputTokens,
	KeyChatContextTokens:        DefaultContextTokens,
	KeyChatTimeout:              DefaultChatTimeout,
	KeyChatDailyBudgetUSD:       DefaultDailyBudgetUSD,
}
