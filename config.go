package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nuts-foundation/nuts-admin/oidc"
	"golang.org/x/crypto/ssh"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/pflag"
)

const defaultPrefix = "NUTS_"
const defaultDelimiter = "."
const configFileFlag = "configfile"
const defaultConfigFile = "config.yaml"

func defaultConfig() Config {
	return Config{
		HTTPPort: 1305,
		Node: Node{
			Address: "http://localhost:8081",
		},
		AccessLogs: true,
		OIDC:       oidc.DefaultConfig(),
	}
}

type Config struct {
	HTTPPort   int    `koanf:"port"`
	BaseURL    string `koanf:"url"`
	Node       Node   `koanf:"node"`
	AccessLogs bool   `koanf:"accesslogs"`
	apiKey     crypto.Signer
	OIDC       oidc.Config `koanf:"oidc"`
}

type Node struct {
	Address string   `koanf:"address"`
	Auth    NodeAuth `koanf:"auth"`
}

type NodeAuth struct {
	// string points to the private key used to sign JWTs. If empty Nuts node API security is not enabled
	KeyFile string `koanf:"keyfile"`
	// User contains the API key user that will go into the iss field. It must match the user with the public key from the authorized_keys file in the Nuts node
	User string `koanf:"user"`
	// Audience dictates the aud field of the created JWT
	Audience string `kaonf:"audience"`
}

func generateSessionKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return nil, err
	}
	return key, nil
}

func (c Config) Validate() error {
	if err := c.OIDC.Validate(); err != nil {
		return fmt.Errorf("oidc config error: %w", err)
	}

	if c.OIDC.Enabled && c.BaseURL == "" {
		return errors.New("url is required when oidc is enabled")
	}

	return nil
}

func (c Config) Print() {
	maskedCopy := c
	if len(maskedCopy.OIDC.Client.Secret) > 0 {
		maskedCopy.OIDC.Client.Secret = "*****"
	}
	data, _ := json.Marshal(maskedCopy)
	logger.Info().Msgf("Config: %s", string(data))
}

func loadConfig() Config {
	flagset := loadFlagSet(os.Args[1:])

	var k = koanf.New(".")

	// Prepare koanf for parsing the config file
	configFilePath := resolveConfigFile(flagset)
	// Check if the file exists
	if _, err := os.Stat(configFilePath); err == nil {
		logger.Info().Msgf("Loading config from file: %s", configFilePath)
		if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
			logger.Fatal().Msgf("error while loading config from file: %v", err)
		}
	} else {
		logger.Info().Msgf("Using default config because no file was found at: %s", configFilePath)
	}

	// load env flags, can't return error
	_ = k.Load(envProvider(), nil)

	config := defaultConfig()

	// Unmarshal values of the config file into the config struct, potentially replacing default values
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error while unmarshalling config: %v", err)
	}

	// Load the API key
	if len(config.Node.Auth.KeyFile) > 0 {
		bytes, err := os.ReadFile(config.Node.Auth.KeyFile)
		if err != nil {
			log.Fatalf("error while reading private key file: %v", err)
		}
		config.apiKey, err = pemToPrivateKey(bytes)
		if err != nil {
			log.Fatalf("error while decoding private key file: %v", err)
		}
		if len(config.Node.Auth.User) == 0 {
			log.Fatal("node.auth.user config is required with node.auth.keyfile")
		}
		if len(config.Node.Auth.Audience) == 0 {
			log.Fatal("node.auth.audience config is required with node.auth.keyfile")
		}
	}

	return config
}

func loadFlagSet(args []string) *pflag.FlagSet {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.String(configFileFlag, defaultConfigFile, "Application config file")
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.Parse(args)
	return f
}

// resolveConfigFile resolves the path of the config file using the following sources:
// 1. commandline params (using the given flags)
// 2. environment vars,
// 3. default location.
func resolveConfigFile(flagset *pflag.FlagSet) string {

	k := koanf.New(defaultDelimiter)

	// load env flags, can't return error
	_ = k.Load(envProvider(), nil)

	// load cmd flags, without a parser, no error can be returned
	_ = k.Load(posflag.Provider(flagset, defaultDelimiter, k), nil)

	configFile := k.String(configFileFlag)
	return configFile
}

func envProvider() *env.Env {
	return env.Provider(defaultPrefix, defaultDelimiter, func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, defaultPrefix)), "_", defaultDelimiter, -1)
	})
}

// pemToPrivateKey converts a PEM encoded private key to a Signer interface. It supports EC, RSA and PKIX PEM encoded strings
func pemToPrivateKey(bytes []byte) (signer crypto.Signer, err error) {
	key, _ := ssh.ParseRawPrivateKey(bytes)
	if key == nil {
		err = errors.New("failed to decode PEM file")
		return
	}

	switch k := key.(type) {
	case *rsa.PrivateKey:
		signer = k
	case *ecdsa.PrivateKey:
		signer = k
	default:
		err = fmt.Errorf("unsupported private key type: %T", k)
	}

	return
}
