package oidc

import (
	"errors"
	"strings"
)

type Config struct {
	Enabled      bool     `koanf:"enabled"`
	MetadataURL  string   `koanf:"metadata"`
	ClientID     string   `koanf:"id"`
	ClientSecret string   `koanf:"secret"`
	Scope        []string `koanf:"scope"`
}

func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Scope: []string{
			"openid",
			"profile",
			"email",
		},
	}
}

func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.MetadataURL == "" {
		return errors.New("metadata_url is required")
	}

	if c.ClientID == "" {
		return errors.New("client_id is required")
	}

	if c.ClientSecret == "" {
		return errors.New("client_secret is required")
	}

	scopeCount := 0
	for _, scope := range c.Scope {
		if strings.TrimSpace(scope) == "" {
			return errors.New("scope cannot be an empty string")
		}
		scopeCount++
	}
	if scopeCount == 0 {
		return errors.New("at lease once scope is required")
	}

	return nil
}
