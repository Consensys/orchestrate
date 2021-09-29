package multi

import (
	"fmt"

	pgstore "github.com/consensys/orchestrate/services/api/store/postgres"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	typeFlag     = "api-store-type"
	TypeViperKey = "api.store.type"
	typeDefault  = postgresType
	typeEnv      = "API_STORE_TYPE"
)

func Flags(f *pflag.FlagSet) {
	storeType(f)
	pgstore.Flags(f)
}

// typeFlag register flag for the Transaction scheduler to select
func storeType(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of API Store (one of %q) Environment variable: %q`, availableTypes, typeEnv)
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
