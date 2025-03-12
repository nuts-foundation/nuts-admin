package main

import (
	"github.com/nuts-foundation/nuts-admin/oidc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MessageHook string

func (m *MessageHook) Run(_ *zerolog.Event, _ zerolog.Level, message string) {
	*m = MessageHook(message)
}

func TestConfig_Print(t *testing.T) {
	var capturedMessage MessageHook
	logger = log.Hook(&capturedMessage)

	t.Run("redacting", func(t *testing.T) {
		t.Run("no OIDC client secret set", func(t *testing.T) {
			config := Config{}
			config.Print()
			assert.Contains(t, string(capturedMessage), `"Secret":""`)
		})
		t.Run("OIDC client secret set", func(t *testing.T) {
			config := Config{
				OIDC: oidc.Config{
					Client: oidc.ClientConfig{
						Secret: "changeme",
					},
				},
			}
			config.Print()
			assert.NotContains(t, string(capturedMessage), "changeme")
			assert.Contains(t, string(capturedMessage), `"Secret":"*****"`)
		})
	})
}
