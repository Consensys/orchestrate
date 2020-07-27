package eea

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/keystore/signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
)

// Signer produce a handler executing Tessera signature
func Signer(k, onetime keystore.KeyStore, ec ethclient.Client) engine.HandlerFunc {
	return signer.GenerateSignerHandler(
		generateSignTx(ec),
		k,
		onetime,
		"Successfully signed transaction for EEA private transaction",
		"EEA signer could not sign the transaction: ",
	)
}

// @TODO Reminder, this code is provisional till we can implement a proper solution with two job https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/pegasyseng/orchestrate/253
// Besu Docs https://besu.hyperledger.org/en/stable/HowTo/Send-Transactions/Creating-Sending-Private-Transactions/
func generateSignTx(ec ethclient.Client) signer.TransactionSignerFunc {
	return func(vault keystore.KeyStore, txctx *engine.TxContext, sender ethcommon.Address, t *ethtypes.Transaction) ([]byte, *ethcommon.Hash, error) {
		// Step 1: Sign Private Transaction
		privRaw, _, err := vault.SignPrivateEEATx(
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
			return nil, nil, errors.FromError(err).ExtendComponent(component)
		}

		// Step 2: Send Private Transaction and obtain enclaveKey
		url, err := proxy.GetURL(txctx)
		if err != nil {
			return nil, nil, errors.FromError(err).ExtendComponent(component)
		}

		enclaveKey, err := ec.PrivDistributeRawTransaction(txctx.Context(), url, hexutil.Encode(privRaw))
		if err != nil {
			return nil, nil, errors.CryptoOperationError(err.Error()).SetComponent(component)
		}
		_ = txctx.Envelope.SetEnclaveKey(enclaveKey.Hex())

		txctx.Logger.WithField("enclavekey", enclaveKey.String()).WithField("nonce", t.Nonce()).
			Warnf("signer: private tx was sent")

		// Step3: Sign marking transaction
		privPContractAddr, err := ec.EEAPrivPrecompiledContractAddr(txctx.Context(), url)
		if err != nil {
			return nil, nil, errors.FromError(err).ExtendComponent(component)
		}

		markingTxNonce, err := txctx.Envelope.GetEEAMarkingNonce()
		if err != nil {
			return nil, nil, errors.FromError(err).ExtendComponent(component)
		}

		markingTx := ethtypes.NewTransaction(
			markingTxNonce,
			privPContractAddr,
			t.Value(),
			t.Gas(),
			t.GasPrice(),
			enclaveKey.Bytes(),
		)

		b, hash, err := vault.SignTx(
			txctx.Context(),
			txctx.Envelope.GetChainID(),
			sender,
			markingTx,
		)

		if err != nil {
			return nil, nil, errors.FromError(err).ExtendComponent(component)
		}

		return b, hash, nil
	}

}
