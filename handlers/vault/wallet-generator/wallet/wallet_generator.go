package wallet

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

// Generator creates and handler responsible to generate wallets
func Generator(s keystore.KeyStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"keygen": "got a keygen request",
			"id":     txctx.Envelope.GetMetadata().GetId(),
		})

		add, err := s.GenerateWallet(txctx.Context())
		if err != nil {
			e := txctx.AbortWithError(err)
			txctx.Logger.WithError(e).Errorf("keygen: could not generate wallet")
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
