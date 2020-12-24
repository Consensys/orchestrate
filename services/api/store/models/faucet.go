package models

import (
	"time"
)

type Faucet struct {
	tableName struct{} `pg:"faucets"` // nolint:unused,structcheck // reason

	UUID            string `pg:",pk"`
	Name            string
	TenantID        string
	ChainRule       string
	CreditorAccount string
	MaxBalance      string
	Amount          string
	Cooldown        string
	CreatedAt       time.Time `pg:"default:now()"`
	UpdatedAt       time.Time `pg:"default:now()"`
}
