package postgres

import (
	"github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
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
