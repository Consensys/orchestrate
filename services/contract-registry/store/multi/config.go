package multi

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store/postgres"
)

type Config struct {
	Postgres *postgres.Config
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Postgres: postgres.NewConfig(vipr),
	}
}

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	postgres.Flags(f)
}
