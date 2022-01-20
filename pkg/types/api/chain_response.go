package api

import (
	"math/big"
	"time"

	"github.com/consensys/orchestrate/pkg/types/entities"
)

type ChainResponse struct {
	UUID                      string                     `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`                          // UUID of the registered chain.
	Name                      string                     `json:"name" example:"mainnet"`                                                       // Name of the chain.
	TenantID                  string                     `json:"tenantID" example:"tenant"`                                                    // ID of the tenant executing the API.
	OwnerID                   string                     `json:"ownerID,omitempty" example:"foo"`                                              // ID of the chain owner.
	URLs                      []string                   `json:"urls" example:"https://mainnet.infura.io/v3/a73136601e6f4924a0baa4ed880b535e"` // URLs of Ethereum nodes connected to.
	ChainID                   *big.Int                   `json:"chainID" example:"2445" swaggertype:"string"`                                  // [Ethereum chain ID](https://besu.hyperledger.org/en/latest/Concepts/NetworkID-And-ChainID/).
	ListenerDepth             uint64                     `json:"listenerDepth" example:"0"`                                                    // Block depth after which the Transaction Listener considers a block final and processes it.
	ListenerCurrentBlock      uint64                     `json:"listenerCurrentBlock" example:"0"`                                             // Current block.
	ListenerStartingBlock     uint64                     `json:"listenerStartingBlock" example:"5000"`                                         // Block at which the Transaction Listener starts processing transactions
	ListenerBackOffDuration   string                     `json:"listenerBackOffDuration" example:"5s"`                                         // Time to wait before trying to fetch a new mined block.
	ListenerExternalTxEnabled bool                       `json:"listenerExternalTxEnabled" example:"false"`                                    // Whether the chain listens for external transactions not crafted by Orchestrate.
	PrivateTxManager          *entities.PrivateTxManager `json:"privateTxManager,omitempty"`
	Labels                    map[string]string          `json:"labels,omitempty"`                                // List of custom labels.
	CreatedAt                 time.Time                  `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"` // Date and time at which the chain was registered.
	UpdatedAt                 time.Time                  `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"` // Date and time at which the chain details were updated.
}
