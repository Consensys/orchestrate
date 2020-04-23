package faucets

type PostRequest struct {
	Name            string `json:"name" validate:"required"`
	ChainRule       string `json:"chainRule,omitempty" validate:"required"`
	CreditorAccount string `json:"creditorAccount,omitempty" validate:"required,eth_addr"`
	MaxBalance      string `json:"maxBalance,omitempty" validate:"required,isBig"`
	Amount          string `json:"amount,omitempty" validate:"required,isBig"`
	Cooldown        string `json:"cooldown,omitempty" validate:"required,isDuration"`
}

type PatchRequest struct {
	Name            string `json:"name,omitempty" validate:"omitempty"`
	ChainRule       string `json:"chainRule,omitempty" validate:"omitempty"`
	CreditorAccount string `json:"creditorAccount,omitempty" validate:"omitempty,eth_addr"`
	MaxBalance      string `json:"maxBalance,omitempty" validate:"omitempty,isBig"`
	Amount          string `json:"amount,omitempty" validate:"omitempty,isBig"`
	Cooldown        string `json:"cooldown,omitempty" validate:"omitempty,isDuration"`
}
