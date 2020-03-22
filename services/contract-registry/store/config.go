package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres"
)

func init() {
	viper.SetDefault(TypeViperKey, typeDefault)
	_ = viper.BindEnv(TypeViperKey, typeEnv)
}

const (
	postgresType = "postgres"
)

var availableTypes = []string{
	postgresType,
}

const (
	typeFlag     = "contract-registry-type"
	TypeViperKey = "contracts.store-type"
	typeDefault  = "postgres"
	typeEnv      = "CONTRACT_REGISTRY_TYPE"
)

// Type register flag for db user
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Contract Registry (one of %q)
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

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	Type(f)
	postgres.Flags(f)
}
