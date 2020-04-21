package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/postgres"
)

func init() {
	_ = viper.BindEnv(TypeViperKey, typeEnv)
	viper.SetDefault(TypeViperKey, typeDefault)
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

type Config struct {
	Type     string
	Postgres *postgres.Config
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Type:     vipr.GetString(TypeViperKey),
		Postgres: postgres.NewConfig(vipr),
	}
}

func Flags(f *pflag.FlagSet) {
	Type(f)
	postgres.Flags(f)
}
