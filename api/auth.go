package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const basePath = "auth"
const redirectPath = "callback"
const authSuccessPath = "success"

func AuthMiddleware(router EchoRouter, appBaseUrl *url.URL, issuerUrl *url.URL, clientId string, clientSecret string, securedEndpoints []string) func(next echo.HandlerFunc) echo.HandlerFunc {
	middleware := authMiddleware{
		issuerUrl:        issuerUrl.String(),
		applicationUrl:   appBaseUrl,
		redirectUrl:      appBaseUrl.JoinPath(basePath, redirectPath).String(),
		clientId:         clientId,
		clientSecret:     clientSecret,
		mux:              &sync.Mutex{},
		sessions:         map[string]sessionEntry{},
		securedEndpoints: securedEndpoints,
	}

	router.GET(basePath+"/"+redirectPath, middleware.handleCallback)
	router.GET(basePath+"/"+authSuccessPath, func(c echo.Context) error {
		// Show HTML page that says "successfully authenticated", and redirects the browser after 2 seconds
		return c.HTML(http.StatusOK, "<html><head><meta http-equiv=\"refresh\" content=\"2;url="+appBaseUrl.String()+"\" /></head><body>Successfully authenticated, redirecting you back...</body></html>")
	})
	return middleware.authenticate
}

type authMiddleware struct {
	issuerUrl        string
	redirectUrl      string
	applicationUrl   *url.URL
	clientId         string
	clientSecret     string
	mux              *sync.Mutex
	provider         *oidc.Provider
	tokenVerifier    *oidc.IDTokenVerifier
	oauth2Config     oauth2.Config
	sessions         map[string]sessionEntry
	securedEndpoints []string
}

type sessionEntry struct {
	idToken oidc.IDToken
}

func (a *authMiddleware) init(ctx context.Context) error {
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
		return fmt.Errorf("unable to create OpenID Connect provider: %w", err)
	}
	eps, _ := json.Marshal(a.provider.Endpoint())
	log.Logger.Info().Msgf("OpenID Connect provider: %s", string(eps))
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

func (a *authMiddleware) authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var secured bool
		for _, path := range a.securedEndpoints {
			if path == "/" {
				if c.Request().URL.Path == "/" {
					secured = true
					break
				}
			} else if strings.HasPrefix(c.Request().URL.Path, path) {
				secured = true
				break
			}
		}
		if !secured {
			// No need to authenticate this path
			return next(c)
		}
		log.Logger.Info().Msgf("Auth for: %s", c.Request().URL.Path)

		sessionID, err := c.Cookie("sid")
		if err == nil {
			a.mux.Lock()
			_, hasSession := a.sessions[sessionID.Value]
			a.mux.Unlock()
			if hasSession {
				// user is authenticated
				return next(c)
			}
		}
		// Not authenticated, need to login
		if err := a.init(c.Request().Context()); err != nil {
			log.Logger.Error().Err(err).Msg("OpenID initialization error")
			return errors.New("auth init failure")
		}

		redirectUrl := a.oauth2Config.AuthCodeURL("TODO:STATE")
		// Check redirect URL to prevent invalid/recursive redirects
		if !(strings.HasPrefix(redirectUrl, "https://") || strings.HasPrefix(redirectUrl, "http://")) {
			log.Logger.Error().Msgf("Invalid redirect URL: %s", redirectUrl)
			return errors.New("auth redirect failure")
		}
		log.Logger.Info().Msgf("User not authenticated, redirecting to OpenID login (%s)", redirectUrl)
		return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
	}
}

func (a *authMiddleware) handleCallback(c echo.Context) error {
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
		return errors.New("missing id_token claim in OAuth2 token")
	}

	// Parse and verify ID Token payload.
	idToken, err := a.tokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return fmt.Errorf("failed to verify ID Token: %w", err)
	}
	idTokenInfo, _ := json.MarshalIndent(idToken, "", "  ")
	println(string(idTokenInfo))

	log.Logger.Info().Msgf("User authenticated: %s", idToken.Subject)
	// Store session
	a.mux.Lock()
	defer a.mux.Unlock()
	sessionID := uuid.NewString()
	a.sessions[sessionID] = sessionEntry{
		idToken: *idToken,
	}
	c.SetCookie(&http.Cookie{
		Name:    "sid",
		Value:   sessionID,
		Expires: time.Now().Add(15 * time.Minute),
		//HttpOnly: true,
		Path: "/",
		// TODO: secure session
	})
	return c.Redirect(http.StatusTemporaryRedirect, authSuccessPath)
}
