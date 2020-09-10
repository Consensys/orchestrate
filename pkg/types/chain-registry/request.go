package chainregistry

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Envelope holds information for a Faucet candidate
type Request struct {
	Chain       *models.Chain
	Beneficiary ethcommon.Address
	Candidates  map[string]Faucet
}
