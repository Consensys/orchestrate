package models

import (
	"time"

	genuuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

type Chain struct {
	tableName struct{} `pg:"chains"` // nolint:unused,structcheck // reason

	UUID                      string                   `json:"uuid" pg:",pk"`
	Name                      string                   `json:"name" validate:"required_with=UUID"`
	TenantID                  string                   `json:"tenantID" validate:"required_with=UUID"`
	URLs                      []string                 `json:"urls" pg:"urls,array" validate:"required_with=UUID,unique,dive,url"`
	ListenerDepth             *uint64                  `json:"listenerDepth,omitempty"`
	ListenerCurrentBlock      *uint64                  `json:"listenerCurrentBlock,string,omitempty"`
	ListenerStartingBlock     *uint64                  `json:"listenerStartingBlock,string,omitempty"`
	ListenerBackOffDuration   *string                  `json:"listenerBackOffDuration,omitempty" validate:"required_with=UUID,omitempty,isDuration"`
	ListenerExternalTxEnabled *bool                    `json:"listenerExternalTxEnabled,omitempty"`
	CreatedAt                 *time.Time               `json:"createdAt"`
	UpdatedAt                 *time.Time               `json:"updatedAt,omitempty"`
	PrivateTxManagers         []*PrivateTxManagerModel `json:"privateTxManagers,omitempty"`
}

func (c *Chain) IsValid() bool {
	err := c.Validate(true)
	return err == nil
}

//nolint:gocritic
func (c Chain) Validate(isNewChain bool) error {
	// Remove UUID for chain updates so we do not validate required fields
	if !isNewChain {
		c.UUID = ""
	}

	err := utils.GetValidator().Struct(c)
	return err
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
