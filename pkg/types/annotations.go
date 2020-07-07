package types

type Annotations struct {
	OneTimeKey bool   `json:"oneTimeKey,omitempty"`
	ChainID    string `json:"chainID,omitempty"`
}
