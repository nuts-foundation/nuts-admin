package oidc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/quasoft/memstore"
)

type OIDC struct {
	baseURL     string
	signInUrl   string
	callbackURL string
}

func Setup(config Config, baseURL string, e *echo.Echo, authConfig AuthConfig) error {
	const name = "openid-connect"

	o := &OIDC{
		baseURL: baseURL,
	}

	if o.baseURL == "" {
		o.baseURL = "/"
	}
	o.signInUrl = fmt.Sprintf("%s/auth/%s", o.baseURL, name)
	o.callbackURL = fmt.Sprintf("%s/auth/%s/callback", o.baseURL, name)

	authConfig.redirectURL = o.signInUrl

	// Set up a Session Store for Goth
	sessionStore := memstore.NewMemStore(
		securecookie.GenerateRandomKey(32),
		securecookie.GenerateRandomKey(32),
	)
	gothic.Store = sessionStore

	// Setup OIDC provider for Goth
	provider, err := openidConnect.New(
		config.Client.ID,
		config.Client.Secret,
		o.callbackURL,
		config.MetadataURL,
		config.Scope...,
	)
	if err != nil {
		return err
	}
	if provider == nil {
		return errors.New("oidc provider failed to initialize")
	}
	goth.UseProviders(provider)

	o.RegisterHandlers(e)

	e.Use(AuthWithConfig(authConfig))

	return nil
}

type UserInfo struct {
	NickName string `json:"NickName"`
}

func (o *OIDC) RegisterHandlers(e *echo.Echo) {
	e.GET("/auth/:provider", func(c echo.Context) error {
		provider := c.Param("provider")
		if provider == "" {
			return c.String(http.StatusBadRequest, "Provider not specified")
		}

		q := c.Request().URL.Query()
		q.Add("provider", c.Param("provider"))
		c.Request().URL.RawQuery = q.Encode()

		req := c.Request()
		res := c.Response().Writer
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			return c.JSON(http.StatusOK, gothUser)
		}
		gothic.BeginAuthHandler(res, req)
		return nil
	})

	e.GET("/auth/:provider/callback", func(c echo.Context) error {
		req := c.Request()
		res := c.Response().Writer
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = gothic.StoreInSession("ExpiresAt", strconv.FormatInt(user.ExpiresAt.Unix(), 10), req, res)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.Redirect(http.StatusTemporaryRedirect, o.baseURL)
	})

	e.GET("/auth/:provider/info", func(c echo.Context) error {
		req := c.Request()

		nick, err := gothic.GetFromSession("NickName", req)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, UserInfo{
			NickName: nick,
		})
	})
}

type AuthConfig struct {
	Skipper         middleware.Skipper
	RedirectSkipper middleware.Skipper
	redirectURL     string
}

var DefaultAuthConfig = AuthConfig{
	Skipper:         middleware.DefaultSkipper,
	RedirectSkipper: middleware.DefaultSkipper,
}

func Auth() echo.MiddlewareFunc {
	c := DefaultAuthConfig
	return AuthWithConfig(c)
}

func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			if config.Skipper(c) {
				return next(c) // Skip authentication
			}

			if strings.HasPrefix(req.URL.Path, "/auth/") {
				return next(c) // Skip authentication
			}

			validSession := false

			// Check if the ExpiresAt in the session is not expired
			expiresAt, _ := gothic.GetFromSession("ExpiresAt", req)
			if len(expiresAt) > 0 {
				expiresAtInt, err := strconv.ParseInt(expiresAt, 10, 64)
				if err == nil {
					if expiresAtInt > time.Now().Unix() {
						validSession = true
					}
				}
			}

			// If authorization failed, redirect to login or return 401
			if !validSession {
				if len(config.redirectURL) > 0 && !config.RedirectSkipper(c) {
					return c.Redirect(http.StatusSeeOther, config.redirectURL)
				}

				return echo.ErrUnauthorized
			}

			// Authorized, continue
			return next(c)
		}
	}
}
