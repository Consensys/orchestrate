package types

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	genuuid "github.com/satori/go.uuid"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/tls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const component = "chain-registry.store"

type Chain struct {
	tableName struct{} `pg:"chains"` // nolint:unused,structcheck // reason

	UUID                    string     `json:"uuid,omitempty" pg:",pk"`
	Name                    string     `json:"name,omitempty"`
	TenantID                string     `json:"tenantID,omitempty"`
	URLs                    []string   `json:"urls,omitempty" pg:"urls,array"`
	CreatedAt               *time.Time `json:"createdAt,omitempty"`
	UpdatedAt               *time.Time `json:"updatedAt,omitempty"`
	ListenerDepth           *uint64    `json:"listenerDepth,omitempty"`
	ListenerBlockPosition   *int64     `json:"listenerBlockPosition,string,omitempty"`
	ListenerFromBlock       *int64     `json:"listenerFromBlock,string,omitempty"`
	ListenerBackOffDuration *string    `json:"listenerBackOffDuration,omitempty"`
}

func (c *Chain) IsValid() bool {
	return c.Name != "" && c.TenantID != "" && len(c.URLs) != 0 && c.ListenerBackOffDuration != nil && *c.ListenerBackOffDuration != ""
}

func (c *Chain) SetDefault() {
	if c.UUID == "" {
		c.UUID = genuuid.NewV4().String()
	}
	if !viper.GetBool(multitenancy.EnabledViperKey) && c.TenantID == "" {
		c.TenantID = multitenancy.DefaultTenantIDName
	}
	if c.ListenerDepth == nil {
		depth := uint64(0)
		c.ListenerDepth = &depth
	}
	if c.ListenerBlockPosition == nil {
		blockPosition := int64(-1)
		c.ListenerBlockPosition = &blockPosition
	}
	if c.ListenerFromBlock == nil {
		fromBlock := *c.ListenerBlockPosition
		c.ListenerFromBlock = &fromBlock
	}
	if c.ListenerBackOffDuration == nil {
		backOffDuration := "1s"
		c.ListenerBackOffDuration = &backOffDuration
	}
}

func NewConfig() *dynamic.Configuration {
	return &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:     make(map[string]*dynamic.Router),
			Middlewares: make(map[string]*dynamic.Middleware),
			Services:    make(map[string]*dynamic.Service),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
	}
}

func BuildConfiguration(chains []*Chain) (*dynamic.Configuration, error) {
	config := NewConfig()

	for _, chain := range chains {
		chainKey := strings.Join([]string{chain.TenantID, chain.Name}, "-")

		config.HTTP.Routers[chainKey] = &dynamic.Router{
			EntryPoints: []string{"orchestrate"},
			Priority:    math.MaxInt32,
			Service:     chainKey,
			// We set path Rule with tenantID placeholder so it is parsed by gorilla mux
			// and can be used by middlewares (in particular the authentication middleware)
			Rule:        fmt.Sprintf("Path(`/%s`) || Path(`/{tenantID:%s}/%s`)", chain.UUID, chain.TenantID, chain.Name),
			Middlewares: []string{"orchestrate-auth"},
		}

		servers := make([]dynamic.Server, 0)
		for _, chainURL := range chain.URLs {
			u, err := url.Parse(chainURL)
			if err != nil {
				return nil, errors.FromError(err).ExtendComponent(component)
			}

			servers = append(servers, dynamic.Server{
				Scheme: u.Scheme,
				URL:    chainURL,
			})
		}

		config.HTTP.Services[chainKey] = &dynamic.Service{
			LoadBalancer: &dynamic.ServersLoadBalancer{
				PassHostHeader: func(v bool) *bool { return &v }(false),
				Servers:        servers,
			},
		}
	}

	return config, nil
}
