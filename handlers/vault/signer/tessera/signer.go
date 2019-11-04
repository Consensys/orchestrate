package tessera

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tessera"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer/generic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multi-vault/keystore"
)

// Signer produce a handler executing Tessera signature
func Signer(k keystore.KeyStore, t tessera.Client) engine.HandlerFunc {
	successMsg := "Successfully signed transaction for Tessera private transaction"
	errorMsg := "Tessera signer could not sign the transaction"
	return engine.CombineHandlers(
		txHashSetter(t),
		generic.GenerateSignerHandler(signTx, k, successMsg, errorMsg),
	)
}

func signTx(s keystore.KeyStore, txctx *engine.TxContext, sender common.Address, t *ethtypes.Transaction) ([]byte, *common.Hash, error) {
	b, hash, err := s.SignPrivateTesseraTx(txctx.Envelope.GetChain(), sender, t)
	if err != nil {
		return b, hash, errors.FromError(err).ExtendComponent(component)
	}
	return b, hash, nil
}
