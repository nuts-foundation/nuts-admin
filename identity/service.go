package identity

import (
	"context"
	"errors"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/identity/vdr"
	"github.com/nuts-foundation/nuts-admin/nuts"
	"slices"
	"strings"
)

type Service struct {
	Client           *vdr.Client
	DiscoveryService discovery.Service
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
	identities := make([]Identity, 0)
	for _, didStr := range *response.JSON200 {
		curr, err := did.ParseDID(didStr)
		if err != nil {
			return nil, err
		}
		identities = append(identities, parseIdentity(*curr))
	}
	return identities, nil
}

func (i Service) Get(ctx context.Context, subjectID string) (*IdentityDetails, error) {
	// Make sure it exists
	// TODO: requires https://github.com/nuts-foundation/nuts-node/pull/2757
	//didDocument, err := i.resolveDID(ctx, did)
	//if err != nil {
	//	return nil, err
	//}
	result := IdentityDetails{
		Identity: parseIdentity(did.MustParseDID(subjectID)),
	}

	// Get DiscoveryService status
	allDiscoveryServices, err := i.DiscoveryService.GetDiscoveryServices(ctx)
	if err != nil {
		return nil, err
	}
	for _, service := range allDiscoveryServices {
		status, err := i.DiscoveryService.ActivationStatus(ctx, service.Id, subjectID)
		if err != nil {
			return nil, err
		}
		result.DiscoveryServices = append(result.DiscoveryServices, *status)
	}
	// Stable order for UI
	slices.SortFunc(result.DiscoveryServices, func(a, b discovery.DIDStatus) int {
		return strings.Compare(a.ServiceID, b.ServiceID)
	})

	// TODO: Get WalletCredentials
	// requires https://github.com/nuts-foundation/nuts-node/pull/2754
	return &result, nil
}

func (i Service) resolveDID(ctx context.Context, did string) (*did.Document, error) {
	httpResponse, err := i.Client.ResolveDID(ctx, did)
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := vdr.ParseResolveDIDResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to resolve DID")
	}
	return &response.JSON200.Document, nil
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
