package oidc

import "errors"

type Config struct {
	Enabled       bool     `koanf:"enabled"`
	MetadataURL   string   `koanf:"metadata_url"`
	ClientID      string   `koanf:"client_id"`
	ClientSecret  string   `koanf:"client_secret"`
	Scope         []string `koanf:"scope"`
	SessionSecret string   `koanf:"session_secret"`
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

	if c.SessionSecret == "" {
		return errors.New("session_secret is required")
	}

	if c.Scope == nil {
		return errors.New("scope is required")
	}

	for _, scope := range c.Scope {
		if scope == "" {
			return errors.New("scope cannot be \"\" required")
		}
	}

	return nil
}
