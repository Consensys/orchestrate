package postgres

import (
	"github.com/ConsenSys/orchestrate/pkg/database/postgres"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	PG *postgres.Config
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		PG: postgres.NewConfig(vipr),
	}
}

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	postgres.PGFlags(f)
}
