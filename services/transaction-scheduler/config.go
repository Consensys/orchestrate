package transactionscheduler

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/multi"
)

type Config struct {
	App          *app.Config
	Store        *multi.Config
	Multitenancy bool
}

func NewConfig(vipr *viper.Viper) *Config {
	return &Config{
		App:          app.NewConfig(vipr),
		Store:        multi.NewConfig(vipr),
		Multitenancy: viper.GetBool(multitenancy.EnabledViperKey),
	}
}

// Flags register flags for Postgres database
func Flags(f *pflag.FlagSet) {
	multi.Flags(f)
	http.Flags(f)
}
