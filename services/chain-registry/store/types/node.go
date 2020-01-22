package types

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/tls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const component = "chain-registry.store"

type Node struct {
	tableName struct{} `sql:"nodes"` //nolint:unused,structcheck // reason

	ID                      string     `json:"id,omitempty" sql:",pk"`
	Name                    string     `json:"name,omitempty"`
	TenantID                string     `json:"tenantID,omitempty"`
	URLs                    []string   `json:"urls,omitempty" sql:"urls,array"`
	CreatedAt               *time.Time `json:"createdAt,omitempty"`
	UpdatedAt               *time.Time `json:"updatedAt,omitempty"`
	ListenerDepth           uint64     `json:"listenerDepth,omitempty"`
	ListenerBlockPosition   int64      `json:"listenerBlockPosition,string,omitempty"`
	ListenerFromBlock       int64      `json:"listenerFromBlock,string,omitempty"`
	ListenerBackOffDuration string     `json:"listenerBackOffDuration,omitempty"`
}

func (n *Node) IsValid() bool {
	return n.Name != "" && n.TenantID != "" && len(n.URLs) != 0 && n.ListenerBackOffDuration != ""
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

func BuildConfiguration(nodes []*Node) (*dynamic.Configuration, error) {
	config := NewConfig()

	for _, node := range nodes {
		nodeID := strings.Join([]string{node.TenantID, node.Name}, "-")

		config.HTTP.Routers[nodeID] = &dynamic.Router{
			EntryPoints: []string{"http"},
			Service:     nodeID,
			// We set path Rule with tenantID placeholder so it is parsed by gorilla mux
			// and can be used by middlewares (in particular the authentication middleware)
			Rule:        fmt.Sprintf("Path(`/%s`) || Path(`/{tenantID:%s}/%s`)", node.ID, node.TenantID, node.Name),
			Middlewares: []string{"orchestrate-auth"},
		}

		servers := make([]dynamic.Server, 0)
		for _, nodeURL := range node.URLs {
			u, err := url.Parse(nodeURL)
			if err != nil {
				return nil, errors.FromError(err).ExtendComponent(component)
			}

			servers = append(servers, dynamic.Server{
				Scheme: u.Scheme,
				URL:    nodeURL,
			})
		}

		config.HTTP.Services[nodeID] = &dynamic.Service{
			LoadBalancer: &dynamic.ServersLoadBalancer{
				PassHostHeader: func(v bool) *bool { return &v }(false),
				Servers:        servers,
			},
		}
	}

	return config, nil
}
