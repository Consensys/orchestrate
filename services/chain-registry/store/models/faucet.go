package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
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

func (f *Faucet) SetDefault() {
	if f.UUID == "" {
		f.UUID = uuid.Must(uuid.NewV4()).String()
	}

	if f.TenantID == "" {
		f.TenantID = multitenancy.DefaultTenant
	}
}
