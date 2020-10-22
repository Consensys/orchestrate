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
func Vault(signer engine.HandlerFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.In.Entrypoint() == viper.GetString(broker.TxSignerViperKey) {
			signer(txctx)
		}
	}
}
