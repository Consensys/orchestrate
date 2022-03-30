package testutils

import (
	"math/big"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/gofrs/uuid"
)

func FakeRegisterChainRequest() *api.RegisterChainRequest {
	return &api.RegisterChainRequest{
		Name: "chain-" + common.RandString(5),
		URLs: []string{"http://chain:8545"},
		Listener: api.RegisterListenerRequest{
			FromBlock:         "latest",
			ExternalTxEnabled: utils.ToPtr(false).(*bool),
		},
		PrivateTxManager: &api.PrivateTxManagerRequest{
			URL:  "http://tessera-eea:8545",
			Type: entities.EEAChainType,
		},
		Labels: map[string]string{
			"label1": common.RandString(5),
			"label2": common.RandString(5),
		},
	}
}

func FakeUpdateChainRequest() *api.UpdateChainRequest {
	return &api.UpdateChainRequest{
		Name: "chain-" + common.RandString(5),
		Listener: &api.UpdateListenerRequest{
			CurrentBlock: 55,
		},
		Labels: map[string]string{
			"label3": common.RandString(5),
		},
	}
}

func FakeChainResponse() *api.ChainResponse {
	return &api.ChainResponse{
		UUID:                      uuid.Must(uuid.NewV4()).String(),
		Name:                      "chain-" + common.RandString(5),
		TenantID:                  multitenancy.DefaultTenant,
		URLs:                      []string{"http://ethereum-node:8545"},
		ChainID:                   big.NewInt(888),
		ListenerDepth:             0,
		ListenerCurrentBlock:      0,
		ListenerStartingBlock:     0,
		ListenerBackOffDuration:   "5s",
		ListenerExternalTxEnabled: utils.ToPtr(false).(*bool),
		PrivateTxManager:          FakePrivateTxManager(),
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
}
