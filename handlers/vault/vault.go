package vault

import (
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// Vault creates a Vault handler
//
// Vault is a fork handler that allows either to sign a transaction
// or generate a new key depending on the input entrypoint.
func Vault(signer, generator engine.HandlerFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		switch txctx.In.Entrypoint() {
		case viper.GetString(broker.TxSignerViperKey):
			signer(txctx)
		case viper.GetString(broker.AccountGeneratorViperKey):
			generator(txctx)
		}
	}
}
