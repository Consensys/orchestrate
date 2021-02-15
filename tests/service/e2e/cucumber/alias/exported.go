package alias

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener"
	txsender "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-sender"
)

const GlobalAka = "global"
const ExternalTxLabel = "externalTx"

var (
	aliases  *Registry
	initOnce = &sync.Once{}
)

func Init(rawTestData string) {
	initOnce.Do(func() {
		if aliases != nil {
			return
		}

		aliases = NewAliasRegistry()

		// Register global aliases
		importGlobalAlias(rawTestData)
	})
}

// GlobalAliasRegistry returns global Alias registry
func GlobalAliasRegistry() *Registry {
	return aliases
}

func importGlobalAlias(rawTestData string) {
	// register internal aliases
	internal := map[string]interface{}{
		"api":                 viper.GetString(client.URLViperKey),
		"api-metrics":         viper.GetString(client.MetricsURLViperKey),
		"api-key":             viper.GetString(key.APIKeyViperKey),
		"tx-sender-metrics":   viper.GetString(txsender.MetricsURLViperKey),
		"tx-listener-metrics": viper.GetString(txlistener.MetricsURLViperKey),
		"key-manager":         viper.GetString(keymanager.URLViperKey),
		"key-manager-metrics": viper.GetString(keymanager.MetricsURLViperKey),
		"external-tx-label":   ExternalTxLabel,
	}

	// import aliases from environment variable
	global := make(map[string]interface{})
	err := json.Unmarshal([]byte(rawTestData), &global)
	if err != nil {
		log.WithError(err).Fatalf("could not parse and register global")
	}
	for k, v := range internal {
		if _, ok := global[k]; ok {
			log.Fatalf("the key '%s' is not allowed in global alias", k)
		}
		global[k] = v
	}

	aliases.Set(global, GlobalAka)
}
