package chainregistry

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
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
	app          *app.Config
	store        *store.Config
	envChains    []string // Chains defined in ENV
	multitenancy bool
}

// Init register flag for the Chain Registry to define initialization state
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Initialize Chain Registry Environment variable: %q`, initEnv)
	f.StringSlice(initFlag, initDefault, desc)
	_ = viper.BindPFlag(InitViperKey, f.Lookup(initFlag))
}

func NewConfig(appCfg *app.Config, storeCfg *store.Config, chains []string, multi bool) Config {
	return Config{
		app:          appCfg,
		store:        storeCfg,
		envChains:    chains,
		multitenancy: multi,
	}
}

func NewConfigFromViper(vipr *viper.Viper) Config {
	return NewConfig(app.NewConfig(vipr),
		store.NewConfig(vipr),
		viper.GetStringSlice(InitViperKey),
		viper.GetBool(multitenancy.EnabledViperKey),
	)
}

func Flags(f *pflag.FlagSet) {
	Type(f)
	store.Flags(f)
}
