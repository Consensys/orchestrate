package dynamic

import (
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	Provider      string
	Configuration *Configuration
}

type Configuration struct {
	Chains map[string]*Chain
}

type Listener struct {
	BlockPosition     int64
	Depth             uint64
	Backoff           time.Duration
	ExternalTxEnabled bool
}

type Chain struct {
	UUID     string
	TenantID string
	Name     string
	URL      string
	ChainID  *big.Int

	Listener *Listener

	Active bool
}

func (n *Chain) SetDefault() {
	if n.UUID == "" {
		n.UUID = uuid.NewV4().String()
	}

	if n.Listener == nil {
		n.Listener = &Listener{}
	}
}
