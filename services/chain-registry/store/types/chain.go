package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
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
	err := c.Validate()
	return err == nil
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("chain name cannot be empty")
	}
	if c.TenantID == "" {
		return fmt.Errorf("chain tenantID cannot be empty")
	}

	if c.ListenerBackOffDuration == nil || *c.ListenerBackOffDuration == "" {
		return fmt.Errorf("chain backOffDuration cannot be empty")
	} else if _, err := time.ParseDuration(*c.ListenerBackOffDuration); err != nil {
		return err
	}

	if len(c.URLs) == 0 {
		return fmt.Errorf("chain URLs cannot be an empty list")
	}

	for _, uri := range c.URLs {
		if err := validateChainURL(uri); err != nil {
			return err
		}
	}

	if _, err := uuid.Parse(c.UUID); c.UUID != "" && err != nil {
		return err
	}

	return nil
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

func validateChainURL(str string) error {
	var re = regexp.MustCompile(`(?m)((https?://)?.*):(\d*)\/?(.*)`)
	if !re.MatchString(str) {
		return fmt.Errorf("malformed URL %s", str)
	}

	return nil
}
