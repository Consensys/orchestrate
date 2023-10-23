package api

import (
	"github.com/consensys/orchestrate/pkg/types/entities"
)

type RegisterChainRequest struct {
	Name             string                   `json:"name" validate:"required" example:"mainnet"`                                                                                             // Name of the chain. Must be unique.
	URLs             []string                 `json:"urls" pg:"urls,array" validate:"required,min=1,unique,dive,url" example:"https://mainnet.infura.io/v3/a73136601e6f4924a0baa4ed880b535e"` // List of URLs of Ethereum nodes to connect to.
	Listener         RegisterListenerRequest  `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
	Headers          map[string]string        `json:"headers,omitempty" validate:"omitempty"`
	Labels           map[string]string        `json:"labels,omitempty"` // List of custom labels. Useful for adding custom information to the chain.
}

type RegisterListenerRequest struct {
	Depth             uint64 `json:"depth,omitempty" example:"0"`                                            // Block depth after which the Transaction Listener considers a block final and processes it (default 0).
	FromBlock         string `json:"fromBlock,omitempty" example:"latest"`                                   // Block from which the Transaction Listener should start processing transactions (default `latest`).
	BackOffDuration   string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration" example:"1s"` // Time to wait before trying to fetch a new mined block (for example `1s` or `1m`, default is `5s`).
	ExternalTxEnabled *bool  `json:"externalTxEnabled,omitempty" example:"false"`                            // Whether to listen to external transactions not crafted by Orchestrate (default `false`).
}

type UpdateChainRequest struct {
	Name             string                   `json:"name,omitempty" example:"mainnet"`
	Listener         *UpdateListenerRequest   `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
	Labels           map[string]string        `json:"labels,omitempty"`
	Headers          map[string]string        `json:"headers,omitempty" validate:"omitempty"`
}

type UpdateListenerRequest struct {
	Depth             uint64 `json:"depth,omitempty" example:"0"`
	BackOffDuration   string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration" example:"1s"`
	ExternalTxEnabled *bool  `json:"externalTxEnabled,omitempty" example:"false"`
	CurrentBlock      uint64 `json:"currentBlock,omitempty" example:"1"`
}

type PrivateTxManagerRequest struct {
	URL  string                        `json:"url" validate:"required,url" example:"http://tessera:3000"`         // Transaction manager endpoint.
	Type entities.PrivateTxManagerType `json:"type" validate:"required,isPrivateTxManagerType" example:"Tessera"` // Currently supports `Tessera` and `EEA``.
}
