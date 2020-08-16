package eea

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/signer"
)

// Signer produce a handler executing Tessera signature
func Signer(k, onetime keystore.KeyStore) engine.HandlerFunc {
	return signer.GenerateSignerHandler(
		generateSignTx(),
		k,
		onetime,
		"Successfully signed transaction for EEA private transaction",
		"EEA signer could not sign the transaction: ",
	)
}

// Besu Docs https://besu.hyperledger.org/en/stable/HowTo/Send-Transactions/Creating-Sending-Private-Transactions/
func generateSignTx() signer.TransactionSignerFunc {
	return func(vault keystore.KeyStore, txctx *engine.TxContext, sender ethcommon.Address, t *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
		// Step 1: Sign Private Transaction
		raw, _, err := vault.SignPrivateEEATx(
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

		// EEA Transaction Hash cannot be computed, so we return nil
		return raw, &ethcommon.Hash{}, err
	}
}
