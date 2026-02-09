package issuer

import (
	"context"
	"strings"
	"time"

	"github.com/nuts-foundation/go-nuts-client/nuts"
	"github.com/nuts-foundation/go-nuts-client/nuts/vcr"
	"github.com/nuts-foundation/nuts-admin/identity"
	"github.com/nuts-foundation/nuts-admin/model"
)

type Service struct {
	IdentityService identity.Service
	VCRClient       *vcr.Client
}

func (s Service) GetIssuedCredentials(ctx context.Context, issuer string, credentialTypes []string) ([]model.IssuedCredential, error) {
	var result []model.IssuedCredential
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
			currentResult := model.IssuedCredential{
				VerifiableCredential: model.ToModel(searchResult.VerifiableCredential),
			}
			if searchResult.Revocation != nil {
				currentResult.Status = "revoked"
			} else if searchResult.VerifiableCredential.ExpirationDate != nil && searchResult.VerifiableCredential.ExpirationDate.Before(time.Now()) {
				currentResult.Status = "expired"
			} else {
				currentResult.Status = "active"
			}
			result = append(result, currentResult)
		}
	}
	return result, nil
}
