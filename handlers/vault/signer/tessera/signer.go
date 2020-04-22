package tessera

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/generic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
)

// Signer produce a handler executing Tessera signature
func Signer(k keystore.KeyStore) engine.HandlerFunc {
	successMsg := "Successfully signed transaction for Tessera private transaction"
	errorMsg := "Tessera signer could not sign the transaction"
	return generic.GenerateSignerHandler(
		signTx,
		k,
		successMsg,
		errorMsg,
	)
}

func signTx(s keystore.KeyStore, txctx *engine.TxContext, sender ethcommon.Address, t *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
	b, hash, err := s.SignPrivateTesseraTx(txctx.Context(), txctx.Envelope.GetChainID(), sender, t)
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
