package ethereum

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/signer"
)

// Signer produce a handler executing Tessera signature
func Signer(k, onetime keystore.KeyStore) engine.HandlerFunc {
	return signer.GenerateSignerHandler(
		signTx,
		k,
		onetime,
		"Successfully signed transaction for ethereum public transaction",
		"Public ethereum signer: could not sign the transaction: ",
	)
}

func signTx(s keystore.KeyStore, txctx *engine.TxContext, sender ethcommon.Address, t *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	txctx.Logger.Tracef("Start signTx from ethereum")
	b, hash, err := s.SignTx(txctx.Context(), txctx.Envelope.ChainID, sender, t)
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
