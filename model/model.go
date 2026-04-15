package model

import (
	"time"

	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/go-nuts-client/nuts/vcr"
)

type VerifiableCredential vc.VerifiableCredential

type CredentialProfile struct {
	Type   string `json:"type" koanf:"type"`
	Issuer string `json:"issuer" koanf:"issuer"`
}

type CredentialWithStatus struct {
	VerifiableCredential
	Status string `json:"status"`
}

func ListToModel(vcs []vc.VerifiableCredential) []VerifiableCredential {
	result := make([]VerifiableCredential, 0)
	for _, credential := range vcs {
		currentCredential := ToModel(credential)
		result = append(result, currentCredential)
	}
	return result
}

func ToModel(credential vc.VerifiableCredential) VerifiableCredential {
	currentCredential := VerifiableCredential(credential)
	if credential.Format() == vc.JWTCredentialProofFormat {
		currentCredential.Proof = []interface{}{
			"jwt",
		}
	}
	return currentCredential
}

func SearchResultToModel(searchResult vcr.SearchVCResult) CredentialWithStatus {
	currentResult := CredentialWithStatus{
		VerifiableCredential: ToModel(searchResult.VerifiableCredential),
	}
	if searchResult.Revocation != nil {
		currentResult.Status = "revoked"
	} else if searchResult.VerifiableCredential.ExpirationDate != nil && searchResult.VerifiableCredential.ExpirationDate.Before(time.Now()) {
		currentResult.Status = "expired"
	} else {
		currentResult.Status = "active"
	}
	return currentResult
}
