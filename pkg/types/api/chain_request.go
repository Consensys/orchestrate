package api

type RegisterChainRequest struct {
	Name             string                   `json:"name" validate:"required" example:"mainnet"`
	URLs             []string                 `json:"urls" pg:"urls,array" validate:"required,min=1,unique,dive,url" example:"https://mainnet.infura.io/v3/a73136601e6f4924a0baa4ed880b535e"`
	Listener         RegisterListenerRequest  `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
}

type RegisterListenerRequest struct {
	Depth             uint64 `json:"depth,omitempty" example:"0"`
	FromBlock         string `json:"fromBlock,omitempty" example:"latest"`
	BackOffDuration   string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration" example:"1s"`
	ExternalTxEnabled bool   `json:"externalTxEnabled,omitempty" example:"false"`
}

type UpdateChainRequest struct {
	Name             string                   `json:"name,omitempty" example:"mainnet"`
	Listener         *UpdateListenerRequest   `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
}

type UpdateListenerRequest struct {
	Depth             uint64 `json:"depth,omitempty" example:"0"`
	BackOffDuration   string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration" example:"1s"`
	ExternalTxEnabled bool   `json:"externalTxEnabled,omitempty" example:"false"`
	CurrentBlock      uint64 `json:"currentBlock,omitempty" example:"1"`
}

type PrivateTxManagerRequest struct {
	URL  string `json:"url" validate:"required,url" example:"http://tessera:3000"`
	Type string `json:"type" validate:"required,isPrivateTxManagerType" example:"Tessera"`
}
