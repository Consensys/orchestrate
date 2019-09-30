package generic

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

// TODO: should be moved to pkg. Should be a method on Envelope.GetTransaction() *ethtypes.Transaction

// TransactionFromTxContext extract a transaction from a transaction context
func TransactionFromTxContext(txctx *engine.TxContext) *ethtypes.Transaction {

	if txctx.Envelope.GetArgs().GetCall().IsConstructor() {
		// Create contract deployment transaction
		return ethtypes.NewContractCreation(
			txctx.Envelope.GetTx().GetTxData().GetNonce(),
			txctx.Envelope.GetTx().GetTxData().GetValueBig(),
			txctx.Envelope.GetTx().GetTxData().GetGas(),
			txctx.Envelope.GetTx().GetTxData().GetGasPriceBig(),
			txctx.Envelope.GetTx().GetTxData().GetDataBytes(),
		)
	}

	// Create transaction
	address := txctx.Envelope.GetTx().GetTxData().GetTo().Address()
	return ethtypes.NewTransaction(
		txctx.Envelope.GetTx().GetTxData().GetNonce(),
		address,
		txctx.Envelope.GetTx().GetTxData().GetValueBig(),
		txctx.Envelope.GetTx().GetTxData().GetGas(),
		txctx.Envelope.GetTx().GetTxData().GetGasPriceBig(),
		txctx.Envelope.GetTx().GetTxData().GetDataBytes(),
	)
}
