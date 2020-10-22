package alias

import (
	"context"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/client"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	txcrafter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-crafter"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sender"
	txsigner "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer"
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
		"chain-registry":            viper.GetString(chainregistry.URLViperKey),
		"chain-registry-metrics":    viper.GetString(chainregistry.MetricsURLViperKey),
		"contract-registry":         viper.GetString(contractregistry.GRPCURLViperKey),
		"contract-registry-metrics": viper.GetString(contractregistry.MetricsURLViperKey),
		"contract-registry-http":    viper.GetString(contractregistry.HTTPURLViperKey),
		"tx-scheduler":              viper.GetString(txscheduler.URLViperKey),
		"tx-scheduler-metrics":      viper.GetString(txscheduler.MetricsURLViperKey),
		"api-key":                   viper.GetString(key.APIKeyViperKey),
		"tx-crafter-metrics":        viper.GetString(txcrafter.MetricsURLViperKey),
		"tx-signer-metrics":         viper.GetString(txsigner.MetricsURLViperKey),
		"tx-sender-metrics":         viper.GetString(txsender.MetricsURLViperKey),
		"tx-listener-metrics":       viper.GetString(txlistener.MetricsURLViperKey),
		"identity-manager":          viper.GetString(identitymanager.URLViperKey),
		"identity-manager-metrics":  viper.GetString(identitymanager.MetricsURLViperKey),
		"key-manager":               viper.GetString(keymanager.URLViperKey),
		"key-manager-metrics":       viper.GetString(txscheduler.MetricsURLViperKey),
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
