package api

import (
	"github.com/nuts-foundation/nuts-admin/identity"
	"net/http"

	"github.com/labstack/echo/v4"
)

var _ ServerInterface = (*Wrapper)(nil)

type Wrapper struct {
	Auth     auth
	Identity identity.Service
}

func (w Wrapper) GetIdentities(ctx echo.Context) error {
	identities, err := w.Identity.List(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, identities)
}

func (w Wrapper) CreateIdentity(ctx echo.Context) error {
	identityRequest := CreateIdentityJSONRequestBody{}
	if err := ctx.Bind(&identityRequest); err != nil {
		return err
	}
	id, err := w.Identity.Create(ctx.Request().Context(), identityRequest.DidQualifier)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, id)
}

func (w Wrapper) CheckSession(ctx echo.Context) error {
	// If this function is reached, it means the session is still valid
	return ctx.NoContent(http.StatusNoContent)
}

func (w Wrapper) CreateSession(ctx echo.Context) error {
	sessionRequest := CreateSessionRequest{}
	if err := ctx.Bind(&sessionRequest); err != nil {
		return err
	}

	if !w.Auth.CheckCredentials(sessionRequest.Username, sessionRequest.Password) {
		return echo.NewHTTPError(http.StatusForbidden, "invalid credentials")
	}

	token, err := w.Auth.CreateJWT(sessionRequest.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, CreateSessionResponse{Token: string(token)})
}
