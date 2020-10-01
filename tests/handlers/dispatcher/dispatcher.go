package dispatcher

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/chanregistry"
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

		// Loop over key functions until we succeed in dispatching envelope to channel indexed by key
		for _, keyOf := range keyOfs {
			// Compute envelope key
			key, err := keyOf(txctx)
			if err != nil {
				// Could not compute key
				continue
			}

			// Dispatch envelope
			err = reg.Send(
				key,
				&envelope,
			)
			if err != nil {
				// Could not dispatch
				continue
			}

			txctx.Logger.WithFields(log.Fields{
				"key": key,
			}).Debugf("dispatcher: envelope dispatched")

			// Envelope has been successfully dispatched so we return
			return
		}

		// Failed in dispatching envelope thus we ignore
		txctx.Logger.Warnf("dispatcher: untracked envelope not dispatched")
	}
}
