package postgres

import (
	"github.com/go-pg/pg/v9"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

type Config struct {
	PG *pg.Options
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		PG: postgres.NewOptions(vipr),
	}
}

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	postgres.PGFlags(f)
}
