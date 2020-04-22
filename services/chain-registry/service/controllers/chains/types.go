package chains

type ListenerRequest struct {
	Depth             *uint64 `json:"depth,omitempty"`
	FromBlock         *string `json:"fromBlock,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}

type PrivateTxManagerRequest struct {
	URL  string `json:"url" validate:"omitempty,url"`
	Type string `json:"type" validate:"omitempty,isPrivateTxManagerType"`
}

type PostRequest struct {
	Name             string                   `json:"name" validate:"required"`
	URLs             []string                 `json:"urls" pg:"urls,array" validate:"required,min=1,unique,dive,url"`
	Listener         *ListenerPostRequest     `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
}

type ListenerPostRequest struct {
	Depth             *uint64 `json:"depth,omitempty"`
	FromBlock         *string `json:"fromBlock,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}

type PatchRequest struct {
	Name             string                   `json:"name,omitempty"`
	URLs             []string                 `json:"urls,omitempty" pg:"urls,array" validate:"omitempty,min=1,unique,dive,url"`
	Listener         *ListenerPatchRequest    `json:"listener,omitempty"`
	PrivateTxManager *PrivateTxManagerRequest `json:"privateTxManager,omitempty"`
}

type ListenerPatchRequest struct {
	Depth             *uint64 `json:"depth,omitempty"`
	CurrentBlock      *uint64 `json:"currentBlock,string,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty" validate:"omitempty,isDuration"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}
