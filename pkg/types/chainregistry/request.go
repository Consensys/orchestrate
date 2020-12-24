package chainregistry

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

// Envelope holds information for a Faucet candidate
type Request struct {
	Chain       *models.Chain
	Beneficiary string
	Candidates  map[string]*entities.Faucet
}
