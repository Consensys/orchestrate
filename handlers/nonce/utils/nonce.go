package utils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
)

func GetNonce(ec ethclient.ChainStateReader, txctx *engine.TxContext, url string) (uint64, error) {
	switch {
	case txctx.Envelope.IsEeaSendPrivateTransactionPrivacyGroup():
		return ec.PrivNonce(txctx.Context(), url, txctx.Envelope.MustGetFromAddress(), txctx.Envelope.GetPrivacyGroupID())
	case txctx.Envelope.IsEeaSendPrivateTransactionPrivateFor():
		return ec.PrivEEANonce(txctx.Context(), url, txctx.Envelope.MustGetFromAddress(), txctx.Envelope.GetPrivateFrom(), txctx.Envelope.GetPrivateFor())
	default:
		return ec.PendingNonceAt(txctx.Context(), url, txctx.Envelope.MustGetFromAddress())
	}
}
