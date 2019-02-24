package handlers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
)

// Decoder creates a decode handler
func Decoder(r services.ABIRegistry) types.HandlerFunc {
	return func(ctx *types.Context) {
		ctx.Logger = ctx.Logger.WithFields(log.Fields{
			"chain.id": ctx.T.Chain().ID.Text(16),
			"tx.hash":  ctx.T.Receipt().TxHash.Hex(),
		})

		// For each log in receipt
		for _, l := range ctx.T.Receipt().Logs {
			if len(l.Topics) == 0 {
				// This scenario is not supposed to append
				err := fmt.Errorf("Invalid receipt (no topics in log)")
				ctx.Logger.WithError(err).Errorf("decoder: invalid receipt")
				ctx.AbortWithError(err)
				return
			}

			// Retrieve event ABI from registry
			event, err := r.GetEventBySig(l.Topics[0].Hex())
			if err != nil {
				ctx.Logger.WithError(err).Errorf("decoder: could not retrieve event ABI")
				ctx.AbortWithError(err)
				return
			}

			// Decode log
			mapping, err := ethereum.Decode(&event, &l.Log)
			if err != nil {
				ctx.Logger.WithError(err).Errorf("decoder: could not decode log")
				ctx.AbortWithError(err)
				return
			}

			// Set decoded data on log
			l.SetDecodedData(mapping)

			ctx.Logger.WithFields(log.Fields{
				"log": mapping,
			}).Debug("decoder: log decoded")
		}
	}
}
