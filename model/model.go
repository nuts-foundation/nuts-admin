package model

import (
	"time"

	"github.com/nuts-foundation/go-did/vc"
)

type VerifiableCredential vc.VerifiableCredential

type CredentialProfile struct {
	Type   string `json:"type" koanf:"type"`
	Issuer string `json:"issuer" koanf:"issuer"`
}

type IssuedCredential struct {
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

// GetCredentialStatus determines if a credential is active, expired, or revoked
func GetCredentialStatus(credential vc.VerifiableCredential, isRevoked bool) string {
	if isRevoked {
		return "revoked"
	}
	if credential.ExpirationDate != nil && credential.ExpirationDate.Before(time.Now()) {
		return "expired"
	}
	return "active"
}
