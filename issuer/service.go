package issuer

import (
	"context"
	"sort"
	"strings"

	"github.com/nuts-foundation/go-nuts-client/nuts"
	"github.com/nuts-foundation/go-nuts-client/nuts/vcr"
	"github.com/nuts-foundation/nuts-admin/identity"
	"github.com/nuts-foundation/nuts-admin/model"
)

type Service struct {
	IdentityService identity.Service
	VCRClient       *vcr.Client
}

func (s Service) GetIssuedCredentials(ctx context.Context, issuer string, credentialTypes []string) ([]model.CredentialWithStatus, error) {
	var result []model.CredentialWithStatus
	for _, credentialType := range credentialTypes {
		credentialType = strings.TrimSpace(credentialType)
		if credentialType == "" {
			continue
		}
		httpResponse, err := s.VCRClient.SearchIssuedVCs(ctx, &vcr.SearchIssuedVCsParams{
			CredentialType: credentialType,
			Issuer:         issuer,
		})
		response, err := nuts.ParseResponse(err, httpResponse, vcr.ParseSearchIssuedVCsResponse)
		if err != nil {
			return nil, err
		}
		for _, searchResult := range response.JSON200.VerifiableCredentials {
			result = append(result, model.SearchResultToModel(searchResult))
		}
	}
	// Sort by issuance date, descending (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].VerifiableCredential.IssuanceDate.After(result[j].VerifiableCredential.IssuanceDate)
	})
	return result, nil
}
