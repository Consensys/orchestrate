package tessera

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// If we need to send a transaction to Tessera enclave we first need to send a transaction data to Tessera
// to get a hash of this data. Then we need to replace data in a transaction object with a hash returned by
// Tessera enclave. We then need to sign the updated transaction
func txHashSetter(tesseraClient tessera.Client) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		if txctx.Builder.GetData() == "" {
			err := errors.DataError("can not send transaction without data to Tessera").SetComponent(component)
			txctx.Logger.WithError(err).Errorf("failed to get transaction hash from Tessera")
			_ = txctx.AbortWithError(err)
			return
		}

		txHash, err := tesseraClient.StoreRaw(
			txctx.Builder.GetChainIDString(),
			txctx.Builder.MustGetDataBytes(),
			txctx.Builder.GetPrivateFrom(),
		)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("failed to get transaction hash from Tessera")
			return
		}

		_ = txctx.Builder.SetData(txHash)
		txctx.Logger.Debugf("Sent transaction body to 'storesaw' endpoint and get txHash to be signed: %s", hexutil.Encode(txHash))
	}
}
