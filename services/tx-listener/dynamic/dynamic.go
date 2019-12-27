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
	Nodes map[string]*Node
}

type Listener struct {
	BlockPosition int64
	Depth         uint64
	Backoff       time.Duration
}

type Node struct {
	ID       string
	TenantID string
	Name     string
	URL      string
	ChainID  *big.Int

	Listener *Listener

	Active bool
}

func (n *Node) SetDefault() {
	if n.ID == "" {
		n.ID = uuid.NewV4().String()
	}

	if n.Listener == nil {
		n.Listener = &Listener{}
	}
}
