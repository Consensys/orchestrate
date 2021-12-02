package entities

import (
	"math/big"
	"time"
)

type Chain struct {
	UUID                      string
	Name                      string
	TenantID                  string
	OwnerID                   string
	URLs                      []string
	ChainID                   *big.Int
	ListenerDepth             uint64
	ListenerCurrentBlock      uint64
	ListenerStartingBlock     uint64
	ListenerBackOffDuration   string
	ListenerExternalTxEnabled bool
	PrivateTxManager          *PrivateTxManager
	Labels                    map[string]string
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
