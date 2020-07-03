package account

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
)

// Generator creates and handler responsible to generate accounts
func Generator(s keystore.KeyStore) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"keygen": "got a keygen request",
		})

		add, err := s.GenerateAccount(txctx.Context())
		if err != nil {
			e := txctx.AbortWithError(err)
			txctx.Logger.WithError(e).Errorf("keygen: could not generate account")
		}

		_ = txctx.Envelope.SetFrom(add)

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"keygen":  "completed a key gen request",
			"address": add.String(),
		})
	}
}
