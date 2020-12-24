package testutils

import (
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

// TODO: To be refactored properly when chain registry types are refactored into types, entities and models
func FakeChain() *models.Chain {
	return &models.Chain{
		UUID:     uuid.Must(uuid.NewV4()).String(),
		Name:     "chainName",
		TenantID: multitenancy.DefaultTenant,
		URLs:     []string{"http://ethereum-node:8545"},
		ChainID:  "888",
	}
}
