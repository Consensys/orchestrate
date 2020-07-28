package types

type Annotations struct {
	OneTimeKey bool   `json:"oneTimeKey,omitempty" example:"true"`
	ChainID    string `json:"chainID,omitempty" example:"1 (mainnet)"`
	Priority   string `json:"priority,omitempty" validate:"isPriority" example:"very-high"`
}
