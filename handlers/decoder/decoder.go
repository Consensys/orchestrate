package decoder

import (
	"fmt"
	"strings"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
)

// Decoder creates a decode handler
func Decoder(r registry.Registry) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id": txctx.Envelope.GetChain().GetId(),
			"tx.hash":  txctx.Envelope.GetReceipt().GetTxHash(),
		})

		// For each log in receipt
		for _, l := range txctx.Envelope.GetReceipt().GetLogs() {
			if len(l.GetTopics()) == 0 {
				// This scenario is not supposed to append
				err := fmt.Errorf("invalid receipt (no topics in log)")
				txctx.Logger.WithError(err).Errorf("decoder: invalid receipt")
				_ = txctx.AbortWithError(err)
				return
			}

			// Retrieve event ABI from registry
			event, err := r.GetEventBySig(l.Topics[0])
			if err != nil {
				txctx.Logger.WithError(err).Errorf("decoder: could not retrieve event ABI")
				_ = txctx.AbortWithError(err)
				return
			}

			// Decode log
			mapping, err := decoder.Decode(event, l)
			if err != nil {
				txctx.Logger.WithError(err).Errorf("decoder: could not decode log")
				_ = txctx.AbortWithError(err)
				return
			}

			// Set decoded data on log
			l.DecodedData = mapping
			l.Event = GetAbi(event)

			txctx.Logger.WithFields(log.Fields{
				"log": mapping,
			}).Debug("decoder: log decoded")
		}
	}
}

// GetAbi creates a string ABI (format EventName(argType1, argType2)) from an event
func GetAbi(e *ethAbi.Event) string {
	inputs := make([]string, len(e.Inputs))
	for i := range e.Inputs {
		inputs[i] = fmt.Sprintf("%v", e.Inputs[i].Type)
	}
	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(inputs, ","))
}
