package chainregistry

import (
	"fmt"
	"time"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/dgraph-io/ristretto"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	httpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/metrics"
	metricregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	tcpmetrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/multi"
)

func init() {
	_ = viper.BindEnv(InitViperKey, initEnv)
	viper.SetDefault(InitViperKey, initDefault)
	_ = viper.BindEnv(CacheTTLViperKey, cacheTTLEnv)
	viper.SetDefault(CacheTTLViperKey, cacheDefault)
	_ = viper.BindEnv(MaxIdleConnsPerHostViperKey, maxIdleConnsPerHostEnv)
	viper.SetDefault(MaxIdleConnsPerHostViperKey, maxIdleConnsPerHostDefault)
}

var (
	initFlag     = "chain-registry-init"
	InitViperKey = "chain-registry.init"
	initDefault  []string
	initEnv      = "CHAIN_REGISTRY_INIT"
)

var (
	cacheTTLFlag     = "chain-registry-cache-ttl"
	CacheTTLViperKey = "chain-registry.cache.ttl"
	cacheDefault     = 0 * time.Second
	cacheTTLEnv      = "CHAIN_REGISTRY_CACHE_TTL"
)

var (
	maxIdleConnsPerHostFlag     = "chain-registry-max-idle-connections-per-host"
	MaxIdleConnsPerHostViperKey = "chain-registry.max-idle-connections-per-host"
	maxIdleConnsPerHostDefault  = 50
	maxIdleConnsPerHostEnv      = "CHAIN_REGISTRY_MAXIDLECONNSPERHOST"
)

type Config struct {
	App              *app.Config
	Cache            *ristretto.Config
	ProxyCacheTTL    *time.Duration
	ServersTransport *traefikstatic.ServersTransport
	Store            *multi.Config
	EnvChains        []string // Chains defined in ENV
	Multitenancy     bool
}

// Init register flag for the Chain Registry to define initialization state
func cmdFlags(f *pflag.FlagSet) {
	initDesc := fmt.Sprintf(`Initialize Chain Registry
Environment variable: %q`, initEnv)
	f.StringSlice(initFlag, initDefault, initDesc)
	_ = viper.BindPFlag(InitViperKey, f.Lookup(initFlag))

	cacheDesc := fmt.Sprintf(`Chain Registry Proxy Cache TTL duration (Disabled by default)
Environment variable: %q`, cacheTTLEnv)
	f.Duration(cacheTTLFlag, cacheDefault, cacheDesc)
	_ = viper.BindPFlag(CacheTTLViperKey, f.Lookup(cacheTTLFlag))

	maxIdleConnsPerHostDesc := fmt.Sprintf(`Maximum number of open HTTP connections to a chain proxied
Environment variable: %q`, maxIdleConnsPerHostEnv)
	f.Int(maxIdleConnsPerHostFlag, maxIdleConnsPerHostDefault, maxIdleConnsPerHostDesc)
	_ = viper.BindPFlag(MaxIdleConnsPerHostViperKey, f.Lookup(maxIdleConnsPerHostFlag))
}

func Flags(f *pflag.FlagSet) {
	cmdFlags(f)
	http.Flags(f)
	multi.Flags(f)
	configwatcher.Flags(f)
	metricregistry.Flags(f, httpmetrics.ModuleName, tcpmetrics.ModuleName)
}

func NewConfig(vipr *viper.Viper) *Config {
	cfg := &Config{
		App:   app.NewConfig(vipr),
		Store: multi.NewConfig(vipr),
		Cache: &ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		},
		ServersTransport: &traefikstatic.ServersTransport{
			MaxIdleConnsPerHost: viper.GetInt(MaxIdleConnsPerHostViperKey),
			InsecureSkipVerify:  true,
		},
		EnvChains:    viper.GetStringSlice(InitViperKey),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}

	cacheStr := viper.GetDuration(CacheTTLViperKey)
	if cacheStr != time.Duration(0) {
		cfg.ProxyCacheTTL = &cacheStr
		log.WithoutContext().Info("chain registry proxy cache is enabled.")
	}

	return cfg
}
