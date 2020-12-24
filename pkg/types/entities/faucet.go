package entities

import "time"

type Faucet struct {
	UUID            string
	Name            string
	TenantID        string
	ChainRule       string
	CreditorAccount string
	MaxBalance      string
	Amount          string
	Cooldown        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
