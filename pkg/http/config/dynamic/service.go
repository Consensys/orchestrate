package dynamic

import (
	"reflect"

	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

// +k8s:deepcopy-gen=true

type Service struct {
	Swagger      *Swagger      `json:"swagger,omitempty" toml:"swagger,omitempty" yaml:"swagger,omitempty"`
	ReverseProxy *ReverseProxy `json:"reverseProxy,omitempty" toml:"reverseProxy,omitempty" yaml:"reverseProxy,omitempty"`
	HealthCheck  *HealthCheck  `json:"healthcheck,omitempty" toml:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
	Prometheus   *Prometheus   `json:"prometheus,omitempty" toml:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Dashboard    *Dashboard    `json:"dashboard,omitempty" toml:"dashboard,omitempty" yaml:"dashboard,omitempty"`
	Chains       *Chains       `json:"chains,omitempty" toml:"chains,omitempty" yaml:"chains,omitempty"`
	Contracts    *Contracts    `json:"contracts,omitempty" toml:"contracts,omitempty" yaml:"contracts,omitempty"`
	Transactions *Transactions `json:"transactions,omitempty" toml:"transactions,omitempty" yaml:"transactions,omitempty"`
	Identity     *Identity     `json:"identity,omitempty" toml:"identity,omitempty" yaml:"identity,omitempty"`
	Mock         *Mock         `json:"mock,omitempty" toml:"mock,omitempty" yaml:"mock,omitempty"`
}

func (s *Service) Type() string {
	return utils.ExtractType(s)
}

func (s *Service) Field() (interface{}, error) {
	return utils.ExtractField(s)
}

// +k8s:deepcopy-gen=true

type ReverseProxy struct {
	LoadBalancer       *LoadBalancer                      `json:"loadBalancer,omitempty" toml:"loadBalancer,omitempty" yaml:"loadBalancer,omitempty"`
	ResponseForwarding *traefikdynamic.ResponseForwarding `json:"responseForwarding,omitempty" toml:"responseForwarding,omitempty" yaml:"responseForwarding,omitempty"`
	PassHostHeader     *bool                              `json:"passHostHeader" toml:"passHostHeader" yaml:"passHostHeader"`
}

func (p *ReverseProxy) SetDefaults() {
	if p.LoadBalancer == nil {
		p.LoadBalancer = &LoadBalancer{}
	}

	if p.PassHostHeader == nil {
		p.PassHostHeader = utils.Bool(true)
	}
}

func (p *ReverseProxy) Mergeable(proxy *ReverseProxy) bool {
	p.SetDefaults()
	proxy.SetDefaults()

	savedServers := p.LoadBalancer.Servers
	defer func() {
		p.LoadBalancer.Servers = savedServers
	}()
	p.LoadBalancer.Servers = nil

	savedServersLB := proxy.LoadBalancer.Servers
	defer func() {
		proxy.LoadBalancer.Servers = savedServersLB
	}()
	proxy.LoadBalancer.Servers = nil

	return reflect.DeepEqual(p, proxy)
}

func FromTraefikService(service *traefikdynamic.Service) *Service {
	if service == nil {
		return nil
	}

	if service.LoadBalancer == nil {
		return nil
	}

	var servers []*Server
	for _, srv := range service.LoadBalancer.Servers {
		servers = append(servers, &Server{URL: srv.URL})
	}

	var sticky *Sticky
	if service.LoadBalancer.Sticky != nil && service.LoadBalancer.Sticky.Cookie != nil {
		sticky = &Sticky{
			Cookie: &Cookie{
				Name:     service.LoadBalancer.Sticky.Cookie.Name,
				HTTPOnly: service.LoadBalancer.Sticky.Cookie.HTTPOnly,
				Secure:   service.LoadBalancer.Sticky.Cookie.Secure,
				SameSite: service.LoadBalancer.Sticky.Cookie.SameSite,
			},
		}
	}

	proxy := &ReverseProxy{
		PassHostHeader:     service.LoadBalancer.PassHostHeader,
		ResponseForwarding: service.LoadBalancer.ResponseForwarding,
		LoadBalancer: &LoadBalancer{
			Servers: servers,
			Sticky:  sticky,
		},
	}

	return &Service{ReverseProxy: proxy}
}

func ToTraefikService(service *Service) *traefikdynamic.Service {
	if service == nil {
		return nil
	}

	if service.ReverseProxy == nil || service.ReverseProxy.LoadBalancer == nil {
		return nil
	}

	var servers []traefikdynamic.Server
	for _, srv := range service.ReverseProxy.LoadBalancer.Servers {
		servers = append(servers, traefikdynamic.Server{URL: srv.URL})
	}
	var sticky *traefikdynamic.Sticky
	if service.ReverseProxy.LoadBalancer.Sticky != nil && service.ReverseProxy.LoadBalancer.Sticky.Cookie != nil {
		sticky = &traefikdynamic.Sticky{
			Cookie: &traefikdynamic.Cookie{
				Name:     service.ReverseProxy.LoadBalancer.Sticky.Cookie.Name,
				HTTPOnly: service.ReverseProxy.LoadBalancer.Sticky.Cookie.HTTPOnly,
				Secure:   service.ReverseProxy.LoadBalancer.Sticky.Cookie.Secure,
				SameSite: service.ReverseProxy.LoadBalancer.Sticky.Cookie.SameSite,
			},
		}
	}

	return &traefikdynamic.Service{
		LoadBalancer: &traefikdynamic.ServersLoadBalancer{
			PassHostHeader:     service.ReverseProxy.PassHostHeader,
			ResponseForwarding: service.ReverseProxy.ResponseForwarding,
			Servers:            servers,
			Sticky:             sticky,
		},
	}
}

// +k8s:deepcopy-gen=true

type Swagger struct {
	SpecsFile string `json:"specsFile,omitempty" toml:"specsFile,omitempty" yaml:"specsFile,omitempty"`
}

// +k8s:deepcopy-gen=true

type HealthCheck struct {
}

// +k8s:deepcopy-gen=true

type Prometheus struct{}

// +k8s:deepcopy-gen=true

type Dashboard struct{}

// +k8s:deepcopy-gen=true

type Chains struct{}

// +k8s:deepcopy-gen=true

type Contracts struct{}

// +k8s:deepcopy-gen=true

type Envelopes struct{}

// +k8s:deepcopy-gen=true

type Transactions struct{}

// +k8s:deepcopy-gen=true

type Identity struct{}

// +k8s:deepcopy-gen=true

type Mock struct{}
