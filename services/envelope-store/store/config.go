package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/postgres"
)

func init() {
	viper.SetDefault(TypeViperKey, typeDefault)
	_ = viper.BindEnv(TypeViperKey, typeEnv)
}

const (
	postgresType = "postgres"
	inMemoryType = "in-memory"
)

var availableTypes = []string{
	postgresType,
	inMemoryType,
}

const (
	typeFlag     = "envelope-store-type"
	TypeViperKey = "envelopes.store.type"
	typeDefault  = postgresType
	typeEnv      = "ENVELOPE_STORE_TYPE"
)

// Type register flag for the Envelope Store to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Envelope Store (one of %q)
Environment variable: %q`, availableTypes, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(TypeViperKey, f.Lookup(typeFlag))
}

type Config struct {
	Type     string
	Postgres *pgstore.Config
}

func DefaultConfig() *Config {
	return NewConfig(viper.New())
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		Type:     vipr.GetString(TypeViperKey),
		Postgres: pgstore.NewConfig(vipr),
	}
}

func Flags(f *pflag.FlagSet) {
	Type(f)
	pgstore.Flags(f)
}
