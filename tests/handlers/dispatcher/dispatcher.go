package dispatcher

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber/alias"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
)

// KeyOfFunc should return channel key to dispatch envelope to
type KeyOfFunc func(txtcx *engine.TxContext) (string, error)

// Dispatcher dispatch envelopes into a channel registry
func Dispatcher(reg *chanregistry.ChanRegistry, keyOfs ...KeyOfFunc) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.In == nil {
			panic("input message is nil")
		}

		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"scenario.id": txctx.Envelope.GetContextLabelsValue("scenario.id"),
			"id":          txctx.Envelope.GetID(),
			"job_uuid":    txctx.Envelope.GetJobUUID(),
			"msg.topic":   txctx.In.Entrypoint(),
		})

		// Copy envelope before dispatching (it ensures that envelope can de manipulated in a concurrent safe manner once dispatched)
		envelope := *txctx.Envelope

		if envelope.GetJobUUID() == "" {
			key := "tx.decoded/" + alias.ExternalTxLabel
			err := reg.Send(key, &envelope)
			if err == nil {
				txctx.Logger.WithFields(log.Fields{"key": key}).Debug("dispatcher: external tx envelope dispatched")
			}
			return
		}

		// Loop over key functions until we succeed in dispatching envelope to channel indexed by key
		for _, keyOf := range keyOfs {
			key, err := keyOf(txctx)
			if err != nil {
				continue
			}

			err = reg.Send(key, &envelope)
			if err != nil {
				continue
			}

			txctx.Logger.WithFields(log.Fields{"key": key}).Debug("dispatcher: envelope dispatched")
			return
		}

		txctx.Logger.Warnf("dispatcher: untracked envelope not dispatched")
	}
}
