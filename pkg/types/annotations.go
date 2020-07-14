package types

type Annotations struct {
	OneTimeKey bool   `json:"oneTimeKey,omitempty" example:"true"`
	ChainID    string `json:"chainID,omitempty" example:"1 (mainnet)"`
}
