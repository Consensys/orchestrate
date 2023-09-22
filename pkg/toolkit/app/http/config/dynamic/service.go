package dynamic

import (
	"reflect"

	"github.com/consensys/orchestrate/pkg/utils"
	traefikdynamic "github.com/traefik/traefik/v2/pkg/config/dynamic"
)

// +k8s:deepcopy-gen=true

type Service struct {
	Swagger      *Swagger      `json:"swagger,omitempty" toml:"swagger,omitempty" yaml:"swagger,omitempty"`
	ReverseProxy *ReverseProxy `json:"reverseProxy,omitempty" toml:"reverseProxy,omitempty" yaml:"reverseProxy,omitempty"`
	HealthCheck  *HealthCheck  `json:"healthcheck,omitempty" toml:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
	Prometheus   *Prometheus   `json:"prometheus,omitempty" toml:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Dashboard    *Dashboard    `json:"dashboard,omitempty" toml:"dashboard,omitempty" yaml:"dashboard,omitempty"`
	API          *API          `json:"api,omitempty" toml:"api,omitempty" yaml:"api,omitempty"`
	KeyManager   *KeyManager   `json:"keyManager,omitempty" toml:"keymanager,omitempty" yaml:"keymanager,omitempty"`
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
		//FIXME CUSTOM HEADER change to false for custom header
		//p.PassHostHeader = utils.Bool(false)
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

type API struct{}

// +k8s:deepcopy-gen=true

type KeyManager struct{}

// +k8s:deepcopy-gen=true

type Mock struct{}
