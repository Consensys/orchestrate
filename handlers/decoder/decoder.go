package decoder

import (
	"fmt"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
)

// Decoder creates a decode handler
func Decoder(r registry.Registry) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		txctx.Logger = txctx.Logger.WithFields(log.Fields{
			"chain.id":    txctx.Envelope.GetChain().ID().String(),
			"tx.hash":     txctx.Envelope.GetReceipt().GetTxHash().Hex(),
			"metadata.id": txctx.Envelope.GetMetadata().GetId(),
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
			event, defaultEvents, err := r.GetEventsBySigHash(
				ethCommon.BytesToHash(l.Topics[0].GetRaw()),
				common.AccountInstance{
					Chain:   txctx.Envelope.GetChain(),
					Account: nil,
				},
				uint(len(l.Topics)-1),
			)
			if err != nil {
				txctx.Logger.WithError(err).Errorf("decoder: could not retrieve event ABI")
				_ = txctx.AbortWithError(err)
				return
			}

			if event == nil && len(defaultEvents) == 0 {
				txctx.Logger.Errorf("decoder: failed to load event ABI")
				_ = txctx.AbortWithError(err)
				return
			}

			var mapping map[string]string
			if event != nil {
				mapping, err = decoder.Decode(event, l)
			} else {
				for _, potentialEvent := range defaultEvents {
					mapping, err = decoder.Decode(potentialEvent, l)
					if err == nil {
						event = potentialEvent
						break
					}
				}
			}

			// Decode log
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
