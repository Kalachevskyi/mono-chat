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
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("config parameter Token can't be empty")
	}

	if c.Timeout == 0 {
		return errors.New("config parameter Timeout can't be empty")
	}

	if c.MonoApiToken == "" {
		return errors.New("config parameter MonoApiToken can't be empty")
	}

	return nil
}
