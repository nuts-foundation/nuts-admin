package identity

import (
	"context"
	"errors"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-admin/identity/vdr"
	"github.com/nuts-foundation/nuts-admin/nuts"
	"strings"
)

type Service struct {
	Client *vdr.Client
}

func (i Service) Create(ctx context.Context, shortName string) (*Identity, error) {
	httpResponse, err := i.Client.CreateDID(ctx, vdr.CreateDIDJSONRequestBody{
		Id: &shortName,
	})
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := vdr.ParseCreateDIDResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to create new DID")
	}
	result := parseIdentity(response.JSON200.ID)
	return &result, nil
}

func (i Service) List(ctx context.Context) ([]Identity, error) {
	httpResponse, err := i.Client.ListDIDs(ctx)
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := vdr.ParseListDIDsResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to list DIDs")
	}
	var identities []Identity
	for _, didStr := range *response.JSON200 {
		curr, err := did.ParseDID(didStr)
		if err != nil {
			return nil, err
		}
		identities = append(identities, parseIdentity(*curr))
	}
	return identities, nil
}

func parseIdentity(id did.DID) Identity {
	result := Identity{
		DID: id,
	}
	if strings.Contains(id.ID, ":") {
		result.Name = id.ID[strings.LastIndex(id.ID, ":")+1:]
	} else {
		result.Name = id.String()
	}
	return result
}
