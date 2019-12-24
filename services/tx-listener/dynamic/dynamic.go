package dynamic

import (
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
	BlockPosition   uint64
	Depth           uint64
	BackOffDuration time.Duration
}

type Node struct {
	ID       string
	TenantID string
	Name     string
	URL      string

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
