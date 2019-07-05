package generator

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore"
)

// WalletGenerator creates and handler responsible to generate wallets
func WalletGenerator(s keystore.KeyStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"keygen": "got a keygen request",
			"id":     txctx.Envelope.GetMetadata().GetId(),
		})

		add, err := s.GenerateWallet()
		if err != nil {
			txctx.Logger.WithError(err).Warnf("keygen: could not generate key %v", err)
			_ = txctx.Error(err)
		}

		if txctx.Envelope.GetFrom() != nil {
			txctx.Envelope.GetFrom().SetAddress(add.Bytes())
		} else {
			txctx.Envelope.From = (&ethereum.Account{}).SetAddress(add.Bytes())
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"keygen":  "completed a key gen request",
			"id":      txctx.Envelope.GetMetadata().GetId(),
			"address": add.String(),
		})
	}
}
