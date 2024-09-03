package identity

import (
	"context"
	"errors"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/nuts"
	"github.com/nuts-foundation/nuts-admin/nuts/vcr"
	"github.com/nuts-foundation/nuts-admin/nuts/vdr"
	"slices"
	"strings"
)

type Service struct {
	VDRClient        *vdr.Client
	VCRClient        *vcr.Client
	DiscoveryService discovery.Service
}

func (i Service) Create(ctx context.Context, subject *string) (*Identity, error) {
	if subject != nil && *subject == "" {
		subject = nil
	}
	httpResponse, err := i.VDRClient.CreateDID(ctx, vdr.CreateDIDJSONRequestBody{
		Subject: subject,
	})
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := vdr.ParseCreateDIDResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to create new subject")
	}
	result := Identity{Subject: response.JSON200.Subject}
	for _, didDocument := range response.JSON200.Documents {
		result.DIDs = append(result.DIDs, didDocument.ID.String())
	}
	return &result, nil
}

func (i Service) List(ctx context.Context) ([]Identity, error) {
	//httpResponse, err := i.VDRClient.ListDIDs(ctx)
	//if err != nil {
	//	return nil, nuts.UnwrapAPIError(err)
	//}
	//response, err := vdr.ParseListDIDsResponse(httpResponse)
	//if err != nil {
	//	return nil, err
	//}
	//if response.JSON200 == nil {
	//	return nil, errors.New("unable to list DIDs")
	//}
	//identities := make([]Identity, 0)
	//for _, did := range response.JSON200 {
	//	identities = append(identities, Identity{
	//		Subject: did,
	//		DIDs:    []string{did},
	//	})
	//}
	// TODO: Implement, waiting for https://github.com/nuts-foundation/nuts-node/pull/3336
	identities := make([]Identity, 0)
	identities = append(identities, Identity{
		Subject: "sub1",
	}, Identity{
		Subject: "sub2",
	})
	return identities, nil
}

func (i Service) Get(ctx context.Context, subjectID string) (*IdentityDetails, error) {
	// Make sure it exists
	identity, err := i.getSubject(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	result := IdentityDetails{
		Identity:          *identity,
		DiscoveryServices: make([]discovery.DIDStatus, 0),
		WalletCredentials: make([]vc.VerifiableCredential, 0),
	}

	// Get DIDDocuments
	for _, currentDID := range identity.DIDs {
		resp, err := i.VDRClient.ResolveDID(ctx, currentDID)
		if err != nil {
			return nil, err
		}
		didResponse, err := vdr.ParseResolveDIDResponse(resp)
		if err != nil {
			return nil, err
		}
		result.DIDDocuments = append(result.DIDDocuments, didResponse.JSON200.Document)
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

	// Get WalletCredentials
	for _, currentDID := range identity.DIDs {
		credentials, err := i.credentialsInWallet(ctx, currentDID)
		if err != nil {
			return nil, err
		}
		result.WalletCredentials = append(result.WalletCredentials, credentials...)
	}

	return &result, nil
}

func (i Service) getSubject(ctx context.Context, subject string) (*Identity, error) {
	httpResponse, err := i.VDRClient.SubjectDIDs(ctx, subject)
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := vdr.ParseSubjectDIDsResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to resolve DID")
	}
	return &Identity{
		Subject: subject,
		DIDs:    *response.JSON200,
	}, nil
}

func (i Service) credentialsInWallet(ctx context.Context, id string) ([]vc.VerifiableCredential, error) {
	httpResponse, err := i.VCRClient.GetCredentialsInWallet(ctx, id)
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := vcr.ParseGetCredentialsInWalletResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to list credentials")
	}
	return *response.JSON200, nil
}
