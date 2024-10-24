package model

import "github.com/nuts-foundation/go-did/vc"

type VerifiableCredential vc.VerifiableCredential

func ToModel(vcs []vc.VerifiableCredential) []VerifiableCredential {
	result := make([]VerifiableCredential, 0)
	for _, credential := range vcs {
		currentCredential := VerifiableCredential(credential)
		if credential.Format() == vc.JWTCredentialProofFormat {
			currentCredential.Proof = []interface{}{
				"jwt",
			}
		}
		result = append(result, currentCredential)
	}
	return result
}
