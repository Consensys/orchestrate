package types

import (
	"time"

	genuuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

type Chain struct {
	tableName struct{} `pg:"chains"` // nolint:unused,structcheck // reason

	UUID                      string     `json:"uuid" pg:",pk"`
	Name                      string     `json:"name"`
	TenantID                  string     `json:"tenantID"`
	URLs                      []string   `json:"urls" pg:"urls,array"`
	ListenerDepth             *uint64    `json:"listenerDepth,omitempty"`
	ListenerCurrentBlock      *uint64    `json:"listenerCurrentBlock,string,omitempty"`
	ListenerStartingBlock     *uint64    `json:"listenerStartingBlock,string,omitempty"`
	ListenerBackOffDuration   *string    `json:"listenerBackOffDuration,omitempty"`
	ListenerExternalTxEnabled *bool      `json:"listenerExternalTxEnabled,omitempty"`
	CreatedAt                 *time.Time `json:"createdAt"`
	UpdatedAt                 *time.Time `json:"updatedAt,omitempty"`
}

func (c *Chain) IsValid() bool {
	return c.Name != "" && c.TenantID != "" && len(c.URLs) != 0 && c.ListenerBackOffDuration != nil && *c.ListenerBackOffDuration != ""
}

func (c *Chain) SetDefault() {
	if c.UUID == "" {
		c.UUID = genuuid.NewV4().String()
	}
	if c.TenantID == "" {
		c.TenantID = multitenancy.DefaultTenantIDName
	}
	if c.ListenerDepth == nil {
		depth := uint64(0)
		c.ListenerDepth = &depth
	}
	if c.ListenerStartingBlock == nil {
		startingBlock := uint64(0)
		c.ListenerStartingBlock = &startingBlock
	}
	if c.ListenerCurrentBlock == nil {
		c.ListenerCurrentBlock = c.ListenerStartingBlock
	}
	if c.ListenerBackOffDuration == nil {
		backOffDuration := "1s"
		c.ListenerBackOffDuration = &backOffDuration
	}
	if c.ListenerExternalTxEnabled == nil {
		externalTxEnabled := false
		c.ListenerExternalTxEnabled = &externalTxEnabled
	}
}
