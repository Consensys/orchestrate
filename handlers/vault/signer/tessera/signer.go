package tessera

import (
	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/signer"
)

// Signer produce a handler executing Tessera signature
func Signer(k, onetime keystore.KeyStore) engine.HandlerFunc {
	successMsg := "Successfully signed transaction for Tessera private transaction"
	errorMsg := "Tessera signer could not sign the transaction"
	return signer.GenerateSignerHandler(
		signTx,
		k,
		onetime,
		successMsg,
		errorMsg,
	)
}

func signTx(s keystore.KeyStore, txctx *engine.TxContext, sender ethcommon.Address, t *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {

	// Change *ethtypes.Transaction to *quorumtypes.Transaction in order to set V=37 or 38 when signing
	var quorumTx *quorumtypes.Transaction
	if t.To() == nil {
		quorumTx = quorumtypes.NewContractCreation(t.Nonce(), t.Value(), t.Gas(), t.GasPrice(), t.Data())
	} else {
		quorumTx = quorumtypes.NewTransaction(t.Nonce(), *t.To(), t.Value(), t.Gas(), t.GasPrice(), t.Data())
	}

	b, hash, err := s.SignPrivateTesseraTx(txctx.Context(), txctx.Envelope.GetChainID(), sender, quorumTx)
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
