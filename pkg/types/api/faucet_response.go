package api

import (
	"encoding/json"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type FaucetResponse struct {
	UUID            string            `json:"uuid" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Name            string            `json:"name" validate:"required" example:"faucet-mainnet"`
	TenantID        string            `json:"tenantID,omitempty" example:"foo"`
	ChainRule       string            `json:"chainRule,omitempty" example:"mainnet"`
	CreditorAccount ethcommon.Address `json:"creditorAccount"  example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18" swaggertype:"string"`
	MaxBalance      hexutil.Big       `json:"maxBalance,omitempty" validate:"required" example:"0x16345785D8A0000" swaggertype:"string"`
	Amount          hexutil.Big       `json:"amount,omitempty" validate:"required" example:"0xD529AE9E860000" swaggertype:"string"`
	Cooldown        string            `json:"cooldown,omitempty" validate:"required,isDuration" example:"10s"`
	CreatedAt       time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt       time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}

type faucetResponseJSON struct {
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	TenantID        string    `json:"tenantID"`
	ChainRule       string    `json:"chainRule,omitempty"`
	CreditorAccount string    `json:"creditorAccount"`
	MaxBalance      string    `json:"maxBalance,omitempty"`
	Amount          string    `json:"amount,omitempty"`
	Cooldown        string    `json:"cooldown,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}

func (a *FaucetResponse) MarshalJSON() ([]byte, error) {
	res := &faucetResponseJSON{
		UUID:            a.UUID,
		Name:            a.Name,
		TenantID:        a.TenantID,
		ChainRule:       a.ChainRule,
		CreditorAccount: a.CreditorAccount.String(),
		MaxBalance:      a.MaxBalance.String(),
		Amount:          a.Amount.String(),
		Cooldown:        a.Cooldown,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}

	return json.Marshal(res)
}
