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
	KeyServerFrontendDir        = "SERVER.FRONTEND_DIR"
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
	KeyChatSystemPrompt         = "CHAT.SYSTEM_PROMPT"
	KeyChatWebSearch            = "CHAT.WEB_SEARCH"
	KeyAuthGoogleClientID       = "AUTH.GOOGLE_CLIENT_ID"
	KeyAuthGoogleClientSecret   = "AUTH.GOOGLE_CLIENT_SECRET"
	KeyAuthPublicBaseURL        = "AUTH.PUBLIC_BASE_URL"
	KeyAuthSessionSecret        = "AUTH.SESSION_SECRET"
	KeyAuthSessionTTL           = "AUTH.SESSION_TTL"
	KeyAuthStateTTL             = "AUTH.STATE_TTL"
	KeyAuthCookieSecure         = "AUTH.COOKIE_SECURE"
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
	DefaultFrontendDir       = "../frontend"
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
	DefaultSystemPrompt      = "Bạn là Chat Smith, trợ lý AI. Trả lời ngắn gọn, chính xác, đúng trọng tâm. Dùng markdown khi cần. Không lặp lại câu hỏi. Không mở đầu dài dòng."
	DefaultWebSearch         = true
	DefaultSessionTTL        = 168 * time.Hour
	DefaultStateTTL          = 5 * time.Minute
	DefaultCookieSecure      = true
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
	KeyServerFrontendDir:        DefaultFrontendDir,
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
	KeyChatSystemPrompt:         DefaultSystemPrompt,
	KeyChatWebSearch:            DefaultWebSearch,
	KeyAuthSessionTTL:           DefaultSessionTTL,
	KeyAuthStateTTL:             DefaultStateTTL,
	KeyAuthCookieSecure:         DefaultCookieSecure,
}
