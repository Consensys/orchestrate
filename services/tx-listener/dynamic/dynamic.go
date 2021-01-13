package dynamic

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	Provider      string
	Configuration *Configuration
}

type Configuration struct {
	Chains map[string]*Chain
}

type Listener struct {
	StartingBlock     uint64
	CurrentBlock      uint64
	Depth             uint64
	Backoff           time.Duration
	ExternalTxEnabled bool
}

type Chain struct {
	UUID     string
	TenantID string
	Name     string
	URL      string
	ChainID  string
	Listener Listener
	Active   bool
}

func (n *Chain) SetDefault() {
	if n.UUID == "" {
		n.UUID = uuid.Must(uuid.NewV4()).String()
	}
}
