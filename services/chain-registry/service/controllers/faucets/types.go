package faucets

type PostRequest struct {
	Name            string `json:"name" validate:"required" example:"faucet-mainnet"`
	ChainRule       string `json:"chainRule,omitempty" validate:"required" example:"mainnet"`
	CreditorAccount string `json:"creditorAccount,omitempty" validate:"required,eth_addr" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	MaxBalance      string `json:"maxBalance,omitempty" validate:"required,isBig" example:"100000000000000000 (wei)"`
	Amount          string `json:"amount,omitempty" validate:"required,isBig" example:"60000000000000000 (wei)"`
	Cooldown        string `json:"cooldown,omitempty" validate:"required,isDuration" example:"10s"`
}

type PatchRequest struct {
	Name            string `json:"name,omitempty" validate:"omitempty" example:"faucet-mainnet"`
	ChainRule       string `json:"chainRule,omitempty" validate:"omitempty" example:"mainnet"`
	CreditorAccount string `json:"creditorAccount,omitempty" validate:"omitempty,eth_addr" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	MaxBalance      string `json:"maxBalance,omitempty" validate:"omitempty,isBig" example:"100000000000000000 (wei)"`
	Amount          string `json:"amount,omitempty" validate:"omitempty,isBig" example:"60000000000000000 (wei)"`
	Cooldown        string `json:"cooldown,omitempty" validate:"omitempty,isDuration" example:"10s"`
}
