package tessera

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// If we need to send a transaction to Tessera enclave we first need to send a transaction data to Tessera
// to get a hash of this data. Then we need to replace data in a transaction object with a hash returned by
// Tessera enclave. We then need to sign the updated transaction
func txHashSetter(tesseraClient tessera.Client) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetTx().GetTxData() == nil {
			err := errors.DataError("can not send transaction with no data to Tessera").SetComponent(component)
			txctx.Logger.WithError(err).Errorf("failed to get transaction hash from Tessera")
			_ = txctx.AbortWithError(err)
			return
		}

		txHash, err := tesseraClient.StoreRaw(
			txctx.Envelope.GetChain().ID().String(),
			txctx.Envelope.GetTx().GetTxData().GetDataBytes(),
			txctx.Envelope.GetArgs().GetPrivate().GetPrivateFrom(),
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("failed to get transaction hash from Tessera")
			return
		}

		txctx.Envelope.GetTx().GetTxData().SetData(txHash)
	}
}
