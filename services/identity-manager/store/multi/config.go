package multi

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	pgstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/postgres"
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
	typeFlag     = "identity-manager-type"
	TypeViperKey = "identity-manager.store.type"
	typeDefault  = postgresType
	typeEnv      = "IDENTITY_MANAGER_TYPE"
)

// Type register flag for the Account manager to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Account manager Store (one of %q)
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
	pgstore.Flags(f)
}
