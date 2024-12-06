package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	libDiscovery "github.com/nuts-foundation/go-nuts-client/nuts/discovery"
	"github.com/nuts-foundation/go-nuts-client/nuts/vcr"
	"github.com/nuts-foundation/go-nuts-client/nuts/vdr"
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/identity"
	"github.com/nuts-foundation/nuts-admin/issuer"
	oidc2 "github.com/nuts-foundation/nuts-admin/oidc"
	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/nuts-foundation/nuts-admin/api"
)

const assetPath = "web/dist"

//go:embed web/dist/*
var embeddedFiles embed.FS

var logger zerolog.Logger

func getFileSystem(useFS bool) http.FileSystem {
	if useFS {
		logger.Info().Msg("using live mode")
		return http.FS(os.DirFS(assetPath))
	}
	logger.Debug().Msg("using embed mode")
	fsys, err := fs.Sub(embeddedFiles, assetPath)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// Skip authorization middleware for the status endpoint
func authSkipper(c echo.Context) bool {
	return strings.HasPrefix(c.Request().URL.Path, "/status")
}

func main() {
	logger = zerolog.New(log.Writer()).With().Timestamp().Logger()
	config := loadConfig()
	config.Print()
	if err := config.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid config")
		os.Exit(1)
	}

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = httpErrorHandler
	e.HideBanner = true
	e.HidePort = true
	if config.AccessLogs {
		e.Use(accessLoggerMiddleware(func(c echo.Context) bool {
			return c.Request().URL.Path == "/status"
		}, logger))
	}

	nodeAddress, err := url.Parse(config.Node.Address)
	if err != nil {
		log.Fatalf("unable to parse node address: %s", err)
	}

	if config.OIDC.Enabled {
		_, err := oidc2.SetupOIDC(config.OIDC, config.BaseURL, e, authSkipper)
		if err != nil {
			log.Fatalf("unable to initialize oidc: %s", err)
		}
	}

	// API security
	// TODO
	//tokenGenerator := func() (string, error) {
	//	return "", nil
	//}
	//if config.apiKey != nil {
	//	tokenGenerator = createTokenGenerator(config)
	//}

	vdrClient, _ := vdr.NewClient(config.Node.Address)
	vcrClient, _ := vcr.NewClient(config.Node.Address)
	discoveryClient, _ := libDiscovery.NewClient(config.Node.Address)

	// Initialize wrapper
	discoveryService := discovery.Service{
		Client: discoveryClient,
	}
	identityService := identity.Service{
		VDRClient:        vdrClient,
		VCRClient:        vcrClient,
		DiscoveryService: discoveryService,
	}
	apiWrapper := api.Wrapper{
		Identity:  identityService,
		Discovery: discoveryService,
		IssuerService: issuer.Service{
			IdentityService: identityService,
			VCRClient:       vcrClient,
		},
	}

	api.RegisterHandlers(e, apiWrapper)
	api.ConfigureProxy(logger, e, nodeAddress)

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/status", func(context echo.Context) error {
		return context.String(http.StatusOK, "OK")
	})

	e.GET("/*", echo.WrapHandler(assetHandler))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}

// httpErrorHandler includes the err.Err() string in a { "error": "msg" } json hash
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	type Map map[string]interface{}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else {
		msg = err.Error()
	}

	if _, ok := msg.(string); ok {
		msg = Map{"error": msg}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

func requestsStatusEndpoint(context echo.Context) bool {
	return context.Request().RequestURI == "/status"
}

// createTokenGenerator generates valid API tokens for the Nuts node and signs them with the private key
func createTokenGenerator(config Config) func() (string, error) {
	return func() (string, error) {
		key, err := jwkKey(config.apiKey)
		if err != nil {
			return "", err
		}

		issuedAt := time.Now()
		notBefore := issuedAt
		expires := notBefore.Add(time.Second * time.Duration(5))
		token, err := jwt.NewBuilder().
			Issuer(config.Node.Auth.User).
			Audience([]string{config.Node.Auth.Audience}).
			IssuedAt(issuedAt).
			NotBefore(notBefore).
			Expiration(expires).
			JwtID(uuid.New().String()).
			Build()

		bytes, err := jwt.Sign(token, jwa.SignatureAlgorithm(key.Algorithm()), key)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

func jwkKey(signer crypto.Signer) (key jwk.Key, err error) {
	// ssh key format
	key, err = jwk.New(signer)
	if err != nil {
		return nil, err
	}

	switch k := signer.(type) {
	case *rsa.PrivateKey:
		key.Set(jwk.AlgorithmKey, jwa.PS512)
	case *ecdsa.PrivateKey:
		var alg jwa.SignatureAlgorithm
		alg, err = ecAlg(k)
		key.Set(jwk.AlgorithmKey, alg)
	default:
		err = fmt.Errorf("unsupported signing private key: %T", k)
		return
	}

	err = jwk.AssignKeyID(key)

	return
}

func ecAlg(key *ecdsa.PrivateKey) (alg jwa.SignatureAlgorithm, err error) {
	alg, err = ecAlgUsingPublicKey(key.PublicKey)
	return
}

func ecAlgUsingPublicKey(key ecdsa.PublicKey) (alg jwa.SignatureAlgorithm, err error) {
	switch key.Params().BitSize {
	case 256:
		alg = jwa.ES256
	case 384:
		alg = jwa.ES384
	case 521:
		alg = jwa.ES512
	default:
		err = errors.New("unsupported key")
	}
	return
}

// accessLoggerMiddleware returns middleware that logs metadata of HTTP requests.
// Should be added as the outer middleware to catch all errors and potential status rewrites
func accessLoggerMiddleware(skipper middleware.Skipper, logger zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:     skipper,
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			event := logger.Info().Fields(map[string]interface{}{
				"remote_ip": values.RemoteIP,
				"method":    values.Method,
				"uri":       values.URI,
				"status":    values.Status,
			})
			if logger.GetLevel() >= zerolog.DebugLevel {
				event.Fields(map[string]interface{}{
					"headers": values.Headers,
				})
			}
			event.Msg("HTTP request")
			return nil
		},
	})
}
