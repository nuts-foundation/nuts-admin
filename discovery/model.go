package discovery

import "github.com/nuts-foundation/go-did/vc"

// DIDStatus represents the status of a DID in the discovery service
type DIDStatus struct {
	ServiceID    string                     `json:"id"`
	Active       bool                       `json:"active"`
	Presentation *vc.VerifiablePresentation `json:"vp"`
}
