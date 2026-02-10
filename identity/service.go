package identity

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/go-nuts-client/nuts"
	"github.com/nuts-foundation/go-nuts-client/nuts/vcr"
	"github.com/nuts-foundation/go-nuts-client/nuts/vdr"
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/model"
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
	httpResponse, err := i.VDRClient.CreateSubject(ctx, vdr.CreateSubjectJSONRequestBody{
		Subject: subject,
	})
	response, err := nuts.ParseResponse(err, httpResponse, vdr.ParseCreateSubjectResponse)
	if err != nil {
		return nil, err
	}
	result := Identity{Subject: response.JSON200.Subject}
	for _, didDocument := range response.JSON200.Documents {
		result.DIDs = append(result.DIDs, didDocument.ID.String())
	}
	return &result, nil
}

func (i Service) List(ctx context.Context) ([]Identity, error) {
	httpResponse, err := i.VDRClient.ListSubjects(ctx)
	response, err := nuts.ParseResponse(err, httpResponse, vdr.ParseListSubjectsResponse)
	if err != nil {
		return nil, err
	}
	identities := make([]Identity, 0)
	for subject, subjectDIDs := range *response.JSON200 {
		identities = append(identities, Identity{
			Subject: subject,
			DIDs:    subjectDIDs,
		})
	}
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
		WalletCredentials: make([]model.CredentialWithStatus, 0),
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
	result.WalletCredentials, err = i.credentialsInWallet(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (i Service) getSubject(ctx context.Context, subject string) (*Identity, error) {
	httpResponse, err := i.VDRClient.SubjectDIDs(ctx, subject)
	response, err := nuts.ParseResponse(err, httpResponse, vdr.ParseSubjectDIDsResponse)
	if err != nil {
		return nil, err
	}
	return &Identity{
		Subject: subject,
		DIDs:    *response.JSON200,
	}, nil
}

func (i Service) credentialsInWallet(ctx context.Context, subjectID string) ([]model.CredentialWithStatus, error) {
	httpResponse, err := i.VCRClient.SearchCredentialsInWallet(ctx, subjectID)
	response, err := nuts.ParseResponse(err, httpResponse, vcr.ParseSearchCredentialsInWalletResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to list credentials")
	}
	result := make([]model.CredentialWithStatus, 0)
	for _, searchResult := range response.JSON200.VerifiableCredentials {
		result = append(result, model.SearchResultToModel(searchResult))
	}
	return result, nil
}

func (i Service) getRevocationStatus(ctx context.Context, credential *vc.VerifiableCredential) (string, error) {
	allowUntrustedIssuer := true
	httpResponse, err := i.VCRClient.VerifyVC(ctx, vcr.VerifyVCJSONRequestBody{
		VerifiableCredential: *credential,
		VerificationOptions: &vcr.VCVerificationOptions{
			AllowUntrustedIssuer: &allowUntrustedIssuer,
		},
	})
	if err != nil {
		return "", fmt.Errorf("checking credential validity: %w", err)
	}
	response, err := nuts.ParseResponse(err, httpResponse, vcr.ParseVerifyVCResponse)
	if err != nil {
		return "", err
	}
	if response.JSON200 != nil {
		if response.JSON200.Validity {
			return "", nil
		}
		if response.JSON200.Message != nil {
			return *response.JSON200.Message, nil
		}
		// Not valid, but no message
		return "invalid", nil
	} else {
		return "", errors.New("check failure")
	}
}
