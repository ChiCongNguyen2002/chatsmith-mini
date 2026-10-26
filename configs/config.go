package configs

import (
	"fmt"
	"time"

	"chatsmith-mini/pkg/failure"
)

const maxPort = 65535

type Config struct {
	Server ServerConfig `mapstructure:"SERVER"`
}

type ServerConfig struct {
	Port              int           `mapstructure:"PORT"`
	ReadHeaderTimeout time.Duration `mapstructure:"READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `mapstructure:"READ_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"IDLE_TIMEOUT"`
	MaxHeaderBytes    int           `mapstructure:"MAX_HEADER_BYTES"`
	ShutdownTimeout   time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

type rule struct {
	key string
	ok  bool
}

func (c *Config) validate() error {
	s := c.Server
	rules := []rule{
		{KeyServerPort, s.Port > 0 && s.Port <= maxPort},
		{KeyServerReadHeaderTimeout, s.ReadHeaderTimeout > 0},
		{KeyServerReadTimeout, s.ReadTimeout > 0},
		{KeyServerIdleTimeout, s.IdleTimeout > 0},
		{KeyServerMaxHeaderBytes, s.MaxHeaderBytes > 0},
		{KeyServerShutdownTimeout, s.ShutdownTimeout > 0},
	}
	for _, r := range rules {
		if !r.ok {
			return failure.ErrWithTrace(fmt.Errorf("%s is invalid", r.key))
		}
	}

	return nil
}
