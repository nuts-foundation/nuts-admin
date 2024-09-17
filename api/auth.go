package api

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"sync"
)

func AuthMiddleware(router EchoRouter, appBaseUrl *url.URL, issuerUrl *url.URL, clientId string, clientSecret string) func(next echo.HandlerFunc) echo.HandlerFunc {
	const redirectPath = "auth/callback"
	middleware := authMiddleware{
		issuerUrl:    issuerUrl.String(),
		redirectUrl:  appBaseUrl.JoinPath(redirectPath).String(),
		clientId:     clientId,
		clientSecret: clientSecret,
		mux:          &sync.Mutex{},
	}

	router.GET(redirectPath, middleware.handleRedirect)
	return middleware.authenticate
}

type authMiddleware struct {
	issuerUrl     string
	redirectUrl   string
	clientId      string
	clientSecret  string
	mux           *sync.Mutex
	provider      *oidc.Provider
	tokenVerifier *oidc.IDTokenVerifier
	oauth2Config  oauth2.Config
}

func (a authMiddleware) init(ctx context.Context) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	if a.provider != nil {
		// already initialized
		return nil
	}
	// OpenID Connect Discovery
	var err error
	a.provider, err = oidc.NewProvider(ctx, a.issuerUrl)
	if err != nil {
		return fmt.Errorf("unable t")
	}
	a.tokenVerifier = a.provider.Verifier(&oidc.Config{ClientID: a.clientId})
	a.oauth2Config = oauth2.Config{
		ClientID:     a.clientId,
		ClientSecret: a.clientSecret,
		RedirectURL:  a.redirectUrl,
		Endpoint:     a.provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}
	return nil
}

func (a authMiddleware) authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		httpRequest := c.Request()
		ctx := httpRequest.Context()
		if err := a.init(ctx); err != nil {
			log.Logger.Error().Err(err).Msg("OpenID initialization error")
			return fmt.Errorf("OpenID initialization error")
		}
		// Check if the user is authenticated
		oauth2Token, err := a.oauth2Config.Exchange(ctx, httpRequest.URL.Query().Get("code"))
		if err != nil {
			log.Logger.Warn().Err(err).Msg("OpenID code exchange failed")
			return c.NoContent(http.StatusUnauthorized)
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			// handle missing token
		}

		// Parse and verify ID Token payload.
		idToken, err := a.tokenVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			// handle error
		}

		// Extract custom claims
		var claims struct {
			Email    string `json:"email"`
			Verified bool   `json:"email_verified"`
		}
		if err := idToken.Claims(&claims); err != nil {
			// handle error
		}
	}
}

func (a authMiddleware) handleRedirect(c echo.Context) error {
	http.Redirect(c.Response(), c.Request(), a.oauth2Config.AuthCodeURL(state), http.StatusFound)
}
