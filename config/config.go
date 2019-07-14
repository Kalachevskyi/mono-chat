package config

import (
	"github.com/pkg/errors"
)

type Config struct {
	Token        string
	Debug        bool
	Offset       int
	Timeout      int
	EncodingLog  string // Valid values are "json" and "console",
	MonoApiToken string
	RedisUrl     string // Example localhost:6379
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New(`config parameter "token" can't be empty`)
	}

	if c.Timeout == 0 {
		return errors.New(`config parameter "timeout" can't be empty`)
	}

	if c.MonoApiToken == "" {
		return errors.New(`config parameter "mono_api_token" can't be empty`)
	}

	if c.RedisUrl == "" {
		return errors.New(`config parameter "redis_url" can't be empty`)
	}

	return nil
}
