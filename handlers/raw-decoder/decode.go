package rawdecoder

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const component = "handler.raw_decoder"

func RawDecoder(txctx *engine.TxContext) {
	var tx *types.Transaction

	// TODO: Able to decode Raw EeaSendPrivateTransactions
	if txctx.Envelope.IsEeaSendPrivateTransaction() {
		return
	}

	if txctx.Envelope.GetRaw() == "" {
		_ = txctx.AbortWithError(errors.DataError("no raw filled - could not decode")).ExtendComponent(component)
		return
	}

	err := rlp.DecodeBytes(txctx.Envelope.MustGetRawBytes(), &tx)
	if err != nil {
		_ = txctx.AbortWithError(errors.DataError("could not decode raw - got %v", err)).ExtendComponent(component)
		return
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		_ = txctx.AbortWithError(errors.DataError("could not find sender - got %v", err)).ExtendComponent(component)
		return
	}

	if mode := txctx.Envelope.GetContextLabelsValue("txMode"); mode == "raw" {
		_ = txctx.Envelope.
			SetFrom(msg.From()).
			SetData(tx.Data()).
			SetGas(tx.Gas()).
			SetGasPrice(tx.GasPrice()).
			SetValue(tx.Value()).
			SetNonce(tx.Nonce()).
			SetTxHash(tx.Hash())

		// If not contract creation
		if tx.To() != nil {
			_ = txctx.Envelope.SetTo(*tx.To())
		}
	}
}
