package entities

import (
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Faucet struct {
	UUID            string
	Name            string
	TenantID        string
	ChainRule       string
	CreditorAccount ethcommon.Address
	MaxBalance      hexutil.Big
	Amount          hexutil.Big
	Cooldown        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
