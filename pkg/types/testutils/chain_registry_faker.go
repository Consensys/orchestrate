package testutils

import (
	"time"

	"github.com/consensys/orchestrate/pkg/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/gofrs/uuid"
)

func FakeRegisterChainRequest() *api.RegisterChainRequest {
	return &api.RegisterChainRequest{
		Name: "mainnet",
		URLs: []string{"http://chain:8545"},
		Listener: api.RegisterListenerRequest{
			FromBlock:         "latest",
			ExternalTxEnabled: false,
		},
		PrivateTxManager: &api.PrivateTxManagerRequest{
			URL:  "http://orion:8545",
			Type: entities.OrionChainType,
		},
		Labels: map[string]string{
			"label1": "val1",
			"label2": "val2",
		},
	}
}

func FakeUpdateChainRequest() *api.UpdateChainRequest {
	return &api.UpdateChainRequest{
		Name: "mainnet",
		Listener: &api.UpdateListenerRequest{
			CurrentBlock: 55,
		},
		Labels: map[string]string{
			"label3": "val3",
		},
	}
}

func FakeChainResponse() *api.ChainResponse {
	return &api.ChainResponse{
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
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
}
