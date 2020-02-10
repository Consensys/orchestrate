package eea

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/generic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
)

// Signer produce a handler executing Tessera signature
func Signer(k keystore.KeyStore) engine.HandlerFunc {
	return generic.GenerateSignerHandler(
		signTx,
		k,
		"Successfully signed transaction for EEA private transaction",
		"EEA signer could not sign the transaction: ",
	)
}

func signTx(s keystore.KeyStore, txctx *engine.TxContext, sender common.Address, t *ethtypes.Transaction) ([]byte, *common.Hash, error) {
	b, hash, err := s.SignPrivateEEATx(txctx.Context(), txctx.Builder.ChainID, sender, t, &types.PrivateArgs{
		PrivateFor:    txctx.Builder.PrivateFor,
		PrivateFrom:   txctx.Builder.PrivateFrom,
		PrivateTxType: txctx.Builder.PrivateTxType,
	})
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
