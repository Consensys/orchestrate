package tessera

import (
	"errors"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tessera"
)

// If we need to send a transaction to Tessera enclave we first need to send a transaction data to Tessera
// to get a hash of this data. Then we need to replace data in a transaction object with a hash returned by
// Tessera enclave. We then need to sign the updated transaction
func txHashSetter(tesseraClient tessera.Client) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Envelope.GetTx().GetTxData() == nil {
			_ = txctx.AbortWithError(
				errors.New("transaction sent to Tessera should have data field"),
			)
			return
		}

		txHash, err := tesseraClient.StoreRaw(txctx.Envelope.GetChain().ID().String(), txctx.Envelope.GetTx().GetTxData().GetDataBytes(), txctx.Envelope.GetArgs().GetPrivate().GetPrivateFrom())
		if err != nil {
			_ = txctx.AbortWithError(
				fmt.Errorf("failed to get transaction hash from Tessera: %q", err),
			)
			return
		}

		txctx.Envelope.GetTx().GetTxData().SetData(txHash)
	}
}
