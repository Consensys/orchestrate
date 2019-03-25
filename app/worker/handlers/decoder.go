package handlers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/abi"
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
			l.Event = event.String()

			ctx.Logger.WithFields(log.Fields{
				"log": mapping,
			}).Debug("decoder: log decoded")
		}
	}
}
