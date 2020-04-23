package store

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/postgres"
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
	typeFlag     = "transaction-scheduler-type"
	TypeViperKey = "transaction-scheduler.store.type"
	typeDefault  = postgresType
	typeEnv      = "TRANSACTION_SCHEDULER_TYPE"
)

// Type register flag for the Transaction scheduler to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Transaction scheduler Store (one of %q)
Environment variable: %q`, availableTypes, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(TypeViperKey, f.Lookup(typeFlag))
}

type Config struct {
	Type     string
	Postgres *pgstore.Config
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
