package configs

import "time"

const (
	KeyServerPort              = "SERVER.PORT"
	KeyServerReadHeaderTimeout = "SERVER.READ_HEADER_TIMEOUT"
	KeyServerReadTimeout       = "SERVER.READ_TIMEOUT"
	KeyServerIdleTimeout       = "SERVER.IDLE_TIMEOUT"
	KeyServerMaxHeaderBytes    = "SERVER.MAX_HEADER_BYTES"
	KeyServerShutdownTimeout   = "SERVER.SHUTDOWN_TIMEOUT"
	KeyDatabaseURL             = "DATABASE.URL"
	KeyDatabaseMaxOpenConns    = "DATABASE.MAX_OPEN_CONNS"
	KeyDatabaseMaxIdleConns    = "DATABASE.MAX_IDLE_CONNS"
)

const (
	DefaultConfigFile        = "configs/dev.config.yaml"
	DefaultServerPort        = 8080
	DefaultReadHeaderTimeout = 5 * time.Second
	DefaultReadTimeout       = 15 * time.Second
	DefaultIdleTimeout       = 90 * time.Second
	DefaultMaxHeaderBytes    = 16 << 10
	DefaultShutdownTimeout   = 8 * time.Second
	DefaultMaxOpenConns      = 25
	DefaultMaxIdleConns      = 5
)

var defaults = map[string]any{
	KeyServerPort:              DefaultServerPort,
	KeyServerReadHeaderTimeout: DefaultReadHeaderTimeout,
	KeyServerReadTimeout:       DefaultReadTimeout,
	KeyServerIdleTimeout:       DefaultIdleTimeout,
	KeyServerMaxHeaderBytes:    DefaultMaxHeaderBytes,
	KeyServerShutdownTimeout:   DefaultShutdownTimeout,
	KeyDatabaseMaxOpenConns:    DefaultMaxOpenConns,
	KeyDatabaseMaxIdleConns:    DefaultMaxIdleConns,
}
