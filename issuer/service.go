package issuer

import (
	"context"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-admin/identity"
	"github.com/nuts-foundation/nuts-admin/nuts/vcr"
	"strings"
)

type Service struct {
	IdentityService identity.Service
	VCRClient       *vcr.Client
}

func (s Service) GetIssuedCredentials(ctx context.Context, issuer string, credentialTypes []string) ([]vc.VerifiableCredential, error) {
	var result []vc.VerifiableCredential
	for _, credentialType := range credentialTypes {
		credentialType = strings.TrimSpace(credentialType)
		if credentialType == "" {
			continue
		}
		searchHTTPResponse, err := s.VCRClient.SearchIssuedVCs(ctx, &vcr.SearchIssuedVCsParams{
			CredentialType: credentialType,
			Issuer:         issuer,
		})
		if err != nil {
			return nil, err
		}
		searchResponse, err := vcr.ParseSearchIssuedVCsResponse(searchHTTPResponse)
		if err != nil {
			return nil, err
		}
		for _, searchResult := range searchResponse.JSON200.VerifiableCredentials {
			result = append(result, searchResult.VerifiableCredential)
		}
	}
	return result, nil
}
