package configs

import (
	"fmt"
	"time"

	"chatsmith-mini/pkg/failure"
)

const maxPort = 65535

type Config struct {
	Server   ServerConfig   `mapstructure:"SERVER"`
	Database DatabaseConfig `mapstructure:"DATABASE"`
}

type ServerConfig struct {
	Port              int           `mapstructure:"PORT"`
	ReadHeaderTimeout time.Duration `mapstructure:"READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `mapstructure:"READ_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"IDLE_TIMEOUT"`
	MaxHeaderBytes    int           `mapstructure:"MAX_HEADER_BYTES"`
	ShutdownTimeout   time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
	RequestTimeout    time.Duration `mapstructure:"REQUEST_TIMEOUT"`
}

type DatabaseConfig struct {
	URL          string `mapstructure:"URL"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
}

type rule struct {
	key string
	ok  bool
}

func (c *Config) validate() error {
	s, d := c.Server, c.Database
	rules := []rule{
		{KeyServerPort, s.Port > 0 && s.Port <= maxPort},
		{KeyServerReadHeaderTimeout, s.ReadHeaderTimeout > 0},
		{KeyServerReadTimeout, s.ReadTimeout > 0},
		{KeyServerIdleTimeout, s.IdleTimeout > 0},
		{KeyServerMaxHeaderBytes, s.MaxHeaderBytes > 0},
		{KeyServerShutdownTimeout, s.ShutdownTimeout > 0},
		{KeyServerRequestTimeout, s.RequestTimeout > 0},
		{KeyDatabaseURL, d.URL != ""},
		{KeyDatabaseMaxOpenConns, d.MaxOpenConns > 0},
		{KeyDatabaseMaxIdleConns, d.MaxIdleConns > 0 && d.MaxIdleConns <= d.MaxOpenConns},
	}
	for _, r := range rules {
		if !r.ok {
			return failure.ErrWithTrace(fmt.Errorf("%s is invalid", r.key))
		}
	}

	return nil
}
