package handlers

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"

	InfEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
)

func LogDecoder(ctx *types.Context, r services.ABIRegistry, log *types.Log, i int) {
	event, err := r.GetEventByID(log.Topics[0].Hex())
	if err != nil {
		e := types.Error{
			Err:  err,
			Type: 0, // TODO: add an error type ErrorTypeABIGet
		}
		// Abort execution
		ctx.AbortWithError(e)
		return
	}

	mapping, _ := InfEth.Decode(&event, &log.Log)
	ctx.T.Receipt().Logs[i].SetDecodedData(mapping)

}

// Decoder creates a decode handler
func TransactionDecoder(r services.ABIRegistry) types.HandlerFunc {
	return func(ctx *types.Context) {

		queue := make(chan map[string]string, len(ctx.T.Receipt().Logs))

		for i, log := range ctx.T.Receipt().Logs {

			go LogDecoder(ctx, r, log, i)

		}

		return
	}
}
