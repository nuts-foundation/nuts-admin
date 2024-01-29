package api

import (
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/identity"
	"net/http"

	"github.com/labstack/echo/v4"
)

var _ ServerInterface = (*Wrapper)(nil)

type Wrapper struct {
	Identity  identity.Service
	Discovery discovery.Service
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

func (w Wrapper) GetDiscoveryServices(ctx echo.Context) error {
	services, err := w.Discovery.GetDiscoveryServices(ctx.Request().Context())
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, services)
}
