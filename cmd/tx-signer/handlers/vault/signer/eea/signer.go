package eea

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/cmd/tx-signer/handlers/vault/signer/generic"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/multi-vault/keystore"
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
	b, hash, err := s.SignPrivateEEATx(txctx.Envelope.GetChain(), sender, t, &types.PrivateArgs{
		PrivateFor:    txctx.Envelope.GetArgs().GetPrivate().GetPrivateFor(),
		PrivateFrom:   txctx.Envelope.GetArgs().GetPrivate().GetPrivateFrom(),
		PrivateTxType: txctx.Envelope.GetArgs().GetPrivate().GetPrivateTxType(),
	})
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
