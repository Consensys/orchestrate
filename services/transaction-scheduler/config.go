package transactionscheduler

import (
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	pkgmultitenancy "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
)

type Config struct {
	app          *app.Config
	store        *store.Config
	multitenancy bool
}

func NewConfig(appCfg *app.Config, storeCfg *store.Config, multitenancy bool) Config {
	return Config{
		app:          appCfg,
		store:        storeCfg,
		multitenancy: multitenancy,
	}
}

func NewConfigFromViper(vipr *viper.Viper) Config {
	return NewConfig(app.NewConfig(vipr),
		store.NewConfig(vipr),
		vipr.GetBool(pkgmultitenancy.EnabledViperKey),
	)
}
