package identity

import (
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-admin/discovery"
)

type Identity struct {
	Name string  `json:"name"`
	DID  did.DID `json:"did"`
}

type IdentityDetails struct {
	Identity
	DIDDocument       did.Document              `json:"did_document"`
	DiscoveryServices []discovery.DIDStatus     `json:"discovery_services"`
	WalletCredentials []vc.VerifiableCredential `json:"wallet_credentials"`
}
