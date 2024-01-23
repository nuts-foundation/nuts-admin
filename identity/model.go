package identity

import "github.com/nuts-foundation/go-did/did"

type Identity struct {
	Name string  `json:"name"`
	DID  did.DID `json:"did"`
}
