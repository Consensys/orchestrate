package eea

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/signer"
)

// Signer produce a handler executing Tessera signature
func Signer(k, onetime keystore.KeyStore) engine.HandlerFunc {
	return signer.GenerateSignerHandler(
		signTx,
		k,
		onetime,
		"Successfully signed transaction for EEA private transaction",
		"EEA signer could not sign the transaction: ",
	)
}

func signTx(vault keystore.KeyStore, txctx *engine.TxContext, sender ethcommon.Address, t *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	b, hash, err := vault.SignPrivateEEATx(
		txctx.Context(),
		txctx.Envelope.GetChainID(),
		sender,
		t,
		&types.PrivateArgs{
			PrivateFor:     txctx.Envelope.GetPrivateFor(),
			PrivateFrom:    txctx.Envelope.GetPrivateFrom(),
			PrivacyGroupID: txctx.Envelope.GetPrivacyGroupID(),
			PrivateTxType:  "restricted",
		},
	)
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
