package chainregistry

import (
	"fmt"
	"time"

	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/ping"
	traefiktypes "github.com/containous/traefik/v2/pkg/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
)

const DefaultInternalEntryPointName = "orchestrate"

func init() {
	viper.SetDefault(ProxyAddressViperKey, proxyAddressDefault)
	_ = viper.BindEnv(ProxyAddressViperKey, proxyAddressEnv)
	viper.SetDefault(AddressViperKey, addressDefault)
	_ = viper.BindEnv(AddressViperKey, addressEnv)
	viper.SetDefault(ProvidersThrottleDurationViperKey, providersThrottleDurationDefault)
	_ = viper.BindEnv(ProvidersThrottleDurationViperKey, providersThrottleDurationEnv)
}

func Flags(f *pflag.FlagSet) {
	ProxyAddress(f)
	RegistryAddress(f)
	ProvidersThrottleDuration(f)
}

const (
	proxyAddressFlag     = "chain-proxy-addr"
	ProxyAddressViperKey = "chain.proxy.addr"
	proxyAddressDefault  = ":8080"
	proxyAddressEnv      = "CHAIN_PROXY_ADDRESS"
)

// ProxyAddress register flag for chain proxy address
func ProxyAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address to expose Chain-Registry Proxy to blockchain nodes
Environment variable: %q`, proxyAddressEnv)
	f.String(proxyAddressFlag, proxyAddressDefault, desc)
	_ = viper.BindPFlag(ProxyAddressViperKey, f.Lookup(proxyAddressFlag))
}

const (
	addressFlag     = "chain-registry-addr"
	AddressViperKey = "chain.registry.addr"
	addressDefault  = ":8081"
	addressEnv      = "CHAIN_REGISTRY_ADDRESS"
)

// RegistryAddress register flag for chain proxy address
func RegistryAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address to expose Chain-Registry Registry to blockchain nodes
Environment variable: %q`, addressEnv)
	f.String(addressFlag, addressDefault, desc)
	_ = viper.BindPFlag(AddressViperKey, f.Lookup(addressFlag))
}

const (
	providersThrottleDurationFlag     = "providers-throttle-duration"
	ProvidersThrottleDurationViperKey = "providers.throttle.duration"
	providersThrottleDurationDefault  = time.Second
	providersThrottleDurationEnv      = "PROVIDERS_THROTTLE_DURATION"
)

// ProvidersThrottleDuration register flag for throttle time duration
func ProvidersThrottleDuration(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Duration to wait for, after a configuration reload, before taking into account any new configuration
Environment variable: %q`, providersThrottleDurationEnv)
	f.Duration(providersThrottleDurationFlag, providersThrottleDurationDefault, desc)
	_ = viper.BindPFlag(ProvidersThrottleDurationViperKey, f.Lookup(providersThrottleDurationFlag))
}

func NewConfig() *static.Configuration {
	orchestrateEp := &static.EntryPoint{
		Address: viper.GetString(AddressViperKey),
	}
	orchestrateEp.SetDefaults()

	httpEp := &static.EntryPoint{
		Address: viper.GetString(ProxyAddressViperKey),
	}
	httpEp.SetDefaults()
	httpEp.ProxyProtocol = &static.ProxyProtocol{
		Insecure: true,
	}
	httpEp.ForwardedHeaders = &static.ForwardedHeaders{
		Insecure: true,
	}

	return &static.Configuration{
		Providers: &static.Providers{
			// TODO: make it configurable
			ProvidersThrottleDuration: traefiktypes.Duration(time.Second),
		},
		EntryPoints: static.EntryPoints{
			"http":                        httpEp,
			DefaultInternalEntryPointName: orchestrateEp,
		},
		API: &static.API{
			Dashboard: true,
			// Insecure:  true,
		},
		Ping: &ping.Handler{
			EntryPoint: "orchestrate",
		},
		Metrics: &traefiktypes.Metrics{
			Prometheus: &traefiktypes.Prometheus{
				EntryPoint:           "orchestrate",
				Buckets:              []float64{0.1, 0.3, 1.2, 5},
				AddEntryPointsLabels: true,
				AddServicesLabels:    true,
			},
		},
		ServersTransport: &static.ServersTransport{
			MaxIdleConnsPerHost: 200,
			InsecureSkipVerify:  true,
		},
		Log: &traefiktypes.TraefikLog{
			Level: viper.GetString(logger.LogLevelViperKey),
		},
		AccessLog: &traefiktypes.AccessLog{
			Format: "json",
		},
	}
}
