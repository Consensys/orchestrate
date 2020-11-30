package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

type Faucet struct {
	tableName struct{} `pg:"faucets"` // nolint:unused,structcheck // reason

	UUID      string     `json:"uuid,omitempty" pg:",pk" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Name      string     `json:"name,omitempty" example:"faucet-mainnet"`
	TenantID  string     `json:"tenantID,omitempty" example:"tenant"`
	CreatedAt *time.Time `json:"createdAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`

	ChainRule       string `json:"chainRule" example:"mainnet"`
	CreditorAccount string `json:"creditorAccountAddress,omitempty" validate:"omitempty,eth_addr" example:"0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"`
	MaxBalance      string `json:"maxBalance,omitempty" example:"100000000000000000 (wei)"`
	Amount          string `json:"amount,omitempty" example:"60000000000000000 (wei)"`
	Cooldown        string `json:"cooldown,omitempty" example:"10s"`
}

func (f *Faucet) SetDefault() {
	if f.UUID == "" {
		f.UUID = uuid.Must(uuid.NewV4()).String()
	}

	if f.TenantID == "" {
		f.TenantID = multitenancy.DefaultTenant
	}
}
