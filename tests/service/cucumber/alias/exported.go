package alias

import (
	"context"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

const GlobalAka = "global"

var (
	aliases  *Registry
	initOnce = &sync.Once{}
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if aliases != nil {
			return
		}

		aliases = NewAliasRegistry()

		// Register global aliases
		importGlobalAlias(viper.GetString(cucumberAliasesViperKey))
	})
}

// GlobalAliasRegistry returns global Alias registry
func GlobalAliasRegistry() *Registry {
	return aliases
}

func importGlobalAlias(rawAliases string) {
	// import aliases from environment variable
	global := make(map[string]interface{})
	err := json.Unmarshal([]byte(rawAliases), &global)
	if err != nil {
		log.Fatalf("could not parse and register alias - got %v", err)
	}
	// register internal aliases
	internal := map[string]interface{}{
		"chain-registry":            viper.GetString(chainregistry.ChainRegistryURLViperKey),
		"chain-registry-metrics":    viper.GetString(chainregistry.ChainRegistryMetricsURLViperKey),
		"contract-registry":         viper.GetString(contractregistry.ContractRegistryURLViperKey),
		"contract-registry-metrics": viper.GetString(contractregistry.ContractRegistryMetricsURLViperKey),
		"contract-registry-http":    viper.GetString(contractregistry.ContractRegistryHTTPURLViperKey),
		"envelope-store":            viper.GetString(envelopestore.EnvelopeStoreURLViperKey),
		"envelope-store-metrics":    viper.GetString(envelopestore.EnvelopeStoreMetricsURLViperKey),
		"envelope-store-http":       viper.GetString(envelopestore.EnvelopeStoreHTTPURLViperKey),
		"tx-scheduler":              viper.GetString(txscheduler.TxSchedulerURLViperKey),
		"tx-scheduler-metrics":      viper.GetString(txscheduler.TxSchedulerMetricsURLViperKey),
	}
	for k, v := range internal {
		if _, ok := global[k]; ok {
			log.Fatalf("the key '%s' is not allowed in global alias", k)
		}
		global[k] = v
	}

	aliases.Set(global, GlobalAka)
	log.WithFields(log.Fields{
		"aka":   GlobalAka,
		"value": global,
	}).Infof("parser: global alias registered")
}
