package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/pg"
)

func init() {
	_ = viper.BindEnv(TypeViperKey, typeEnv)
	viper.SetDefault(TypeViperKey, typeDefault)
	_ = viper.BindEnv(InitViperKey, initEnv)
	viper.SetDefault(InitViperKey, initDefault)
}

const (
	postgresType = "postgres"
)

var availableTypes = []string{
	postgresType,
}

const (
	typeFlag     = "chains-store-type"
	TypeViperKey = "chains.registry.type"
	typeDefault  = postgresType
	typeEnv      = "CHAIN_REGISTRY_TYPE"
)

// Type register flag for the Chain Registry to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Chain Registry (one of %q)
Environment variable: %q`, availableTypes, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(TypeViperKey, f.Lookup(typeFlag))
}

var (
	initFlag     = "chain-registry-init"
	InitViperKey = "chain-registry.init"
	initDefault  []string
	initEnv      = "CHAIN_REGISTRY_INIT"
)

// Init register flag for the Chain Registry to define initialization state
func InitRegistry(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Initialize Chain Registry
Environment variable: %q`, initEnv)
	f.StringSlice(initFlag, initDefault, desc)
	_ = viper.BindPFlag(InitViperKey, f.Lookup(initFlag))
}

type Config struct {
	Type     string
	Postgres *pg.Config
	Chains   []string
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Type:     vipr.GetString(TypeViperKey),
		Postgres: pg.NewConfig(vipr),
		Chains:   vipr.GetStringSlice(InitViperKey),
	}
}

func Flags(f *pflag.FlagSet) {
	Type(f)
	InitRegistry(f)
	pg.Flags(f)
}
