package chainregistry

import (
	"fmt"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/dgraph-io/ristretto"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func init() {
	_ = viper.BindEnv(InitViperKey, initEnv)
	viper.SetDefault(InitViperKey, initDefault)
}

var (
	initFlag     = "chain-registry-init"
	InitViperKey = "chain-registry.init"
	initDefault  []string
	initEnv      = "CHAIN_REGISTRY_INIT"
)

type Config struct {
	App              *app.Config
	Cache            *ristretto.Config
	ServersTransport *traefikstatic.ServersTransport
	Store            *store.Config
	EnvChains        []string // Chains defined in ENV
	Multitenancy     bool
}

// Init register flag for the Chain Registry to define initialization state
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Initialize Chain Registry
Environment variable: %q`, initEnv)
	f.StringSlice(initFlag, initDefault, desc)
	_ = viper.BindPFlag(InitViperKey, f.Lookup(initFlag))
}

func Flags(f *pflag.FlagSet) {
	Type(f)
	http.Flags(f)
	store.Flags(f)
	configwatcher.Flags(f)
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:   app.NewConfig(vipr),
		Store: store.NewConfig(vipr),
		Cache: &ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		},
		ServersTransport: &traefikstatic.ServersTransport{
			MaxIdleConnsPerHost: 200,
			InsecureSkipVerify:  true,
		},
		EnvChains:    viper.GetStringSlice(InitViperKey),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}
}
