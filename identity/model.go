package identity

import (
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-admin/discovery"
	"github.com/nuts-foundation/nuts-admin/model"
)

type Identity struct {
	Subject string   `json:"subject"`
	DIDs    []string `json:"dids"`
}

type IdentityDetails struct {
	Identity
	DIDDocuments      []did.Document               `json:"did_documents"`
	DiscoveryServices []discovery.DIDStatus        `json:"discovery_services"`
	WalletCredentials []model.VerifiableCredential `json:"wallet_credentials"`
}
