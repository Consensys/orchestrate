package models

import (
	"time"
)

type Faucet struct {
	tableName struct{} `pg:"faucets"` // nolint:unused,structcheck // reason

	UUID      string     `json:"uuid,omitempty" pg:",pk"`
	Name      string     `json:"name,omitempty"`
	TenantID  string     `json:"tenantID,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	ChainRule       string `json:"chainRule"`
	CreditorAccount string `json:"creditorAccountAddress,omitempty" validate:"omitempty,eth_addr"`
	MaxBalance      string `json:"maxBalance,omitempty"`
	Amount          string `json:"amount,omitempty"`
	Cooldown        string `json:"cooldown,omitempty"`
}
