package oidc

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fakeSignIn(expireInSeconds int64) (string, error) {
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()

	err := gothic.StoreInSession("ExpiresAt", strconv.FormatInt(time.Now().Unix()+expireInSeconds, 10), req, rec)
	if err != nil {
		return "", err
	}

	cookie := rec.Header().Get("Set-Cookie")
	cookieParts := strings.Split(cookie, ";")
	if len(cookieParts) == 0 {
		return "", errors.New("failed to split set-cookie header fields")
	}

	sessionParts := strings.Split(cookieParts[0], "=")
	if len(sessionParts) < 2 || sessionParts[0] != "_gothic_session" {
		return "", errors.New("failed to retrieve session cookie")
	}

	return cookieParts[0], nil
}

func request(e *echo.Echo, method, target, cookie string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	rec := httptest.NewRecorder()
	e.Server.Handler.ServeHTTP(rec, req)
	return rec
}

func TestOIDC(t *testing.T) {
	// build test echo app
	e := echo.New()
	e.HideBanner = true

	e.GET("/protected", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	e.GET("/unprotected", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	e.GET("/api/protected", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	// setup OIDC
	config := Config{
		Enabled: true,
		Client: ClientConfig{
			ID:     "client_id",
			Secret: "client_secret",
		},
		//FIXME: Find a better way to load an openid-configuration
		MetadataURL: "https://accounts.google.com/.well-known/openid-configuration",
		Scope:       []string{"openid", "profile", "email"},
	}

	baseURL := "http://localhost:8080"
	authConfig := AuthConfig{
		Skipper: func(c echo.Context) bool {
			// for the following, skip authorization
			return strings.HasPrefix(c.Request().URL.Path, "/unprotected")
		},
		RedirectSkipper: func(c echo.Context) bool {
			// for the following, don't redirect but return 401
			return strings.HasPrefix(c.Request().URL.Path, "/api/")
		},
	}

	err := Setup(config, baseURL, e, authConfig)
	require.NoError(t, err)

	t.Run("unauthorized to unprotected route gives 200 OK", func(t *testing.T) {
		rec := request(e, echo.GET, "/unprotected", "")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("unauthorized to protected route gives 303 redirect to login", func(t *testing.T) {
		rec := request(e, echo.GET, "/protected", "")
		assert.Equal(t, http.StatusSeeOther, rec.Code)
	})

	t.Run("unauthorized to protected api gives 401 Unauthorized", func(t *testing.T) {
		rec := request(e, echo.GET, "/api/protected", "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("authorized to protected api gives 200 OK", func(t *testing.T) {
		cookie, err := fakeSignIn(60)
		assert.NoError(t, err)
		rec := request(e, echo.GET, "/api/protected", cookie)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("expired session to protected api gives 401 Unauthorized", func(t *testing.T) {
		cookie, err := fakeSignIn(-60)
		assert.NoError(t, err)
		rec := request(e, echo.GET, "/api/protected", cookie)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
