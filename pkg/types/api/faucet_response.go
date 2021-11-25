package api

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type FaucetResponse struct {
	UUID            string            `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Name            string            `json:"name" validate:"required" example:"faucet-mainnet"`
	TenantID        string            `json:"tenantID,omitempty" example:"foo"`
	ChainRule       string            `json:"chainRule,omitempty" example:"mainnet"`
	CreditorAccount ethcommon.Address `json:"creditorAccount"  example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
	MaxBalance      string            `json:"maxBalance,omitempty" validate:"required,isBig" example:"100000000000000000 (wei)"`
	Amount          string            `json:"amount,omitempty" validate:"required,isBig" example:"60000000000000000 (wei)"`
	Cooldown        string            `json:"cooldown,omitempty" validate:"required,isDuration" example:"10s"`
	CreatedAt       time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt       time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}
