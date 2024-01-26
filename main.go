package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/nuts-foundation/nuts-admin/identity"
	"github.com/nuts-foundation/nuts-admin/nuts/client/vdr"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/nuts-foundation/nuts-admin/api"
)

const assetPath = "web/dist"

//go:embed web/dist/*
var embeddedFiles embed.FS

func getFileSystem(useFS bool) http.FileSystem {
	if useFS {
		log.Print("using live mode")
		return http.FS(os.DirFS(assetPath))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(embeddedFiles, assetPath)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func main() {
	config := loadConfig()
	config.Print(log.Writer())

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = httpErrorHandler

	// API security
	// TODO
	//tokenGenerator := func() (string, error) {
	//	return "", nil
	//}
	//if config.apiKey != nil {
	//	tokenGenerator = createTokenGenerator(config)
	//}

	vdrClient, _ := vdr.NewClient(config.Node.Address)

	// Initialize wrapper
	apiWrapper := api.Wrapper{
		Identity: identity.Service{
			Client: vdrClient,
		},
	}

	api.RegisterHandlers(e, apiWrapper)

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
