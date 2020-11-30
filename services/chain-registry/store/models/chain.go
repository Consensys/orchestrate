package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

type Chain struct {
	tableName struct{} `pg:"chains"` // nolint:unused,structcheck // reason

	UUID                      string                   `json:"uuid" pg:",pk" example:"b4374e6f-b28a-4bad-b4fe-bda36eaf849c"`
	Name                      string                   `json:"name" validate:"required_with=UUID" example:"mainnet"`
	TenantID                  string                   `json:"tenantID" validate:"required_with=UUID" example:"tenant"`
	URLs                      []string                 `json:"urls" pg:"urls,array" validate:"required_with=UUID,unique,dive,url" example:"https://mainnet.infura.io/v3/a73136601e6f4924a0baa4ed880b535e"`
	ChainID                   string                   `json:"chainID" validate:"required_with=UUID" example:"1"`
	ListenerDepth             *uint64                  `json:"listenerDepth,omitempty" example:"0"`
	ListenerCurrentBlock      *uint64                  `json:"listenerCurrentBlock,string,omitempty" example:"0"`
	ListenerStartingBlock     *uint64                  `json:"listenerStartingBlock,string,omitempty" example:"0"`
	ListenerBackOffDuration   *string                  `json:"listenerBackOffDuration,omitempty" validate:"required_with=UUID,omitempty,isDuration" example:"1s"`
	ListenerExternalTxEnabled *bool                    `json:"listenerExternalTxEnabled,omitempty" example:"false"`
	CreatedAt                 *time.Time               `json:"createdAt" pg:"default:now()" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt                 *time.Time               `json:"updatedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	PrivateTxManagers         []*PrivateTxManagerModel `json:"privateTxManagers,omitempty" pg:"-" validate:"omitempty,dive,required"`
}

func (c *Chain) IsValid() bool {
	err := c.Validate(true)
	return err == nil
}

// nolint:gocritic
func (c Chain) Validate(isNewChain bool) error {
	// Remove UUID for chain updates so we do not validate required fields
	if !isNewChain {
		c.UUID = ""
	}

	err := utils.GetValidator().Struct(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chain) SetDefault() {
	if c.UUID == "" {
		c.UUID = uuid.Must(uuid.NewV4()).String()
	}
	if c.TenantID == "" {
		c.TenantID = multitenancy.DefaultTenant
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

	c.SetPrivateTxManagersDefault()
}

func (c *Chain) SetPrivateTxManagersDefault() {
	for idx, privTxManager := range c.PrivateTxManagers {
		if privTxManager.ChainUUID == "" && c.UUID != "" {
			c.PrivateTxManagers[idx].ChainUUID = c.UUID
		}
		if privTxManager.UUID == "" {
			c.PrivateTxManagers[idx].UUID = uuid.Must(uuid.NewV4()).String()
		}
	}
}
