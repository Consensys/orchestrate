package chains

type ListenerRequest struct {
	Depth             *uint64 `json:"depth,omitempty"`
	FromBlock         *string `json:"fromBlock,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}

type PrivateTxManagerRequest struct {
	URL  string `json:"url" validate:"omitempty,url" example:"http://tessera:3000"`
	Type string `json:"type" validate:"omitempty,isPrivateTxManagerType" example:"Tessera"`
}

type PostRequest struct {
	Name             string                   `json:"name" validate:"required" example:"mainnet"`
	URLs             []string                 `json:"urls" pg:"urls,array" validate:"required,min=1,unique,dive,url" example:"https://mainnet.infura.io/v3/a73136601e6f4924a0baa4ed880b535e"`
	Listener         *ListenerPostRequest     `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
}

type ListenerPostRequest struct {
	Depth             *uint64 `json:"depth,omitempty" example:"0"`
	FromBlock         *string `json:"fromBlock,omitempty" example:"latest"`
	BackOffDuration   *string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration" example:"1s"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty" example:"false"`
}

type PatchRequest struct {
	Name             string                   `json:"name,omitempty" example:"mainnet"`
	URLs             []string                 `json:"urls,omitempty" pg:"urls,array" validate:"omitempty,min=1,unique,dive,url" example:"https://mainnet.infura.io/v3/a73136601e6f4924a0baa4ed880b535e"`
	Listener         *ListenerPatchRequest    `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
}

type ListenerPatchRequest struct {
	Depth             *uint64 `json:"depth,omitempty" example:"0"`
	CurrentBlock      *uint64 `json:"currentBlock,omitempty" example:"0"`
	BackOffDuration   *string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration" example:"1s"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty" example:"false"`
}
