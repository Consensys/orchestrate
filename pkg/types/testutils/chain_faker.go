package testutils

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/gofrs/uuid"
)

func FakeChain() *entities.Chain {
	return &entities.Chain{
		UUID:                      uuid.Must(uuid.NewV4()).String(),
		Name:                      "ganache",
		TenantID:                  multitenancy.DefaultTenant,
		URLs:                      []string{"http://ethereum-node:8545"},
		ChainID:                   big.NewInt(888),
		ListenerDepth:             0,
		ListenerCurrentBlock:      0,
		ListenerStartingBlock:     0,
		ListenerBackOffDuration:   "5s",
		ListenerExternalTxEnabled: utils.ToPtr(false).(*bool),
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
