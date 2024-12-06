package oidc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

type OIDC struct {
	sessionStore *sessions.FilesystemStore
	provider     *openidConnect.Provider

	name        string
	baseURL     string
	signInUrl   string
	callbackURL string
}

func SetupOIDC(config Config, baseURL string, e *echo.Echo, skipper middleware.Skipper) (*OIDC, error) {
	o := &OIDC{
		name:    "openid-connect",
		baseURL: baseURL,
	}
	if o.baseURL == "" {
		o.baseURL = "/"
	}
	o.signInUrl = fmt.Sprintf("%s/auth/%s", baseURL, o.name)
	o.callbackURL = fmt.Sprintf("%s/auth/%s/callback", baseURL, o.name)

	// Setup a Session Store for Goth
	o.sessionStore = sessions.NewFilesystemStore("", []byte(config.SessionSecret))
	o.sessionStore.MaxLength(81920) // 8Kb is now maximum size of the session
	gothic.Store = o.sessionStore

	// Setup OIDC provider for Goth
	provider, err := openidConnect.New(
		config.ClientID,
		config.ClientSecret,
		o.callbackURL,
		config.MetadataURL,
		config.Scope...,
	)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, errors.New("nil pro provider")
	}
	goth.UseProviders(provider)
	o.provider = provider

	o.RegisterHandlers(e)

	e.Use(AuthWithConfig(AuthConfig{
		Skipper:     skipper,
		RedirectURL: o.signInUrl,
	}))

	return o, nil
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

		err = gothic.StoreInSession("NickName", user.NickName, req, res)
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
	Skipper     middleware.Skipper
	RedirectURL string
}

var DefaultAuthConfig = AuthConfig{
	Skipper:     middleware.DefaultSkipper,
	RedirectURL: "",
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
				return next(c)
			}

			if strings.HasPrefix(req.URL.Path, "/auth/") {
				return next(c)
			}

			nick, err := gothic.GetFromSession("NickName", req)
			if err != nil || len(nick) == 0 {
				if len(config.RedirectURL) > 0 {
					return c.Redirect(http.StatusSeeOther, config.RedirectURL)
				}

				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
