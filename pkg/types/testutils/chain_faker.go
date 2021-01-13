package testutils

import (
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func FakeChain() *entities.Chain {
	return &entities.Chain{
		UUID:                      uuid.Must(uuid.NewV4()).String(),
		Name:                      "ganache",
		TenantID:                  multitenancy.DefaultTenant,
		URLs:                      []string{"http://ethereum-node:8545"},
		ChainID:                   "888",
		ListenerDepth:             0,
		ListenerCurrentBlock:      0,
		ListenerStartingBlock:     0,
		ListenerBackOffDuration:   "5s",
		ListenerExternalTxEnabled: false,
		PrivateTxManager:          FakePrivateTxManager(),
	}
}

func FakePrivateTxManager() *entities.PrivateTxManager {
	return &entities.PrivateTxManager{
		UUID:      uuid.Must(uuid.NewV4()).String(),
		ChainUUID: uuid.Must(uuid.NewV4()).String(),
		URL:       "http://tessera:8545",
		Type:      "Tessera",
	}
}
