package handlers

import (
	"fmt"
	"strings"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
)

// Decoder creates a decode handler
func Decoder(r services.ABIRegistry) worker.HandlerFunc {
	return func(ctx *worker.Context) {
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id": ctx.T.GetChain().GetId(),
			"tx.hash":  ctx.T.GetReceipt().GetTxHash(),
		})

		// For each log in receipt
		for _, l := range ctx.T.GetReceipt().GetLogs() {
			if len(l.GetTopics()) == 0 {
				// This scenario is not supposed to append
				err := fmt.Errorf("Invalid receipt (no topics in log)")
				ctx.Logger.WithError(err).Errorf("decoder: invalid receipt")
				ctx.AbortWithError(err)
				return
			}

			// Retrieve event ABI from registry
			event, err := r.GetEventBySig(l.Topics[0])
			if err != nil {
				ctx.Logger.WithError(err).Errorf("decoder: could not retrieve event ABI")
				ctx.AbortWithError(err)
				return
			}

			// Decode log
			mapping, err := abi.Decode(&event, l)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("decoder: could not decode log")
				ctx.AbortWithError(err)
				return
			}

			// Set decoded data on log
			l.DecodedData = mapping
			l.Event = GetAbi(event)

			ctx.Logger.WithFields(log.Fields{
				"log": mapping,
			}).Debug("decoder: log decoded")
		}
	}
}

// GetAbi creates a string ABI (format EventName(argType1, argType2)) from an event
func GetAbi(e ethAbi.Event) string {
	inputs := make([]string, len(e.Inputs))
	for i, input := range e.Inputs {
		inputs[i] = fmt.Sprintf("%v", input.Type)
	}
	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(inputs, ","))
}
