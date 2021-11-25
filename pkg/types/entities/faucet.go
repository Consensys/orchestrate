package entities

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type Faucet struct {
	UUID            string
	Name            string
	TenantID        string
	ChainRule       string
	CreditorAccount ethcommon.Address
	MaxBalance      string
	Amount          string
	Cooldown        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
