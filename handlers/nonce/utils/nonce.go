package utils

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
)

func GetNonce(ctx context.Context, ec ethclient.ChainStateReader, e *tx.Envelope, url string) (uint64, error) {
	return ec.PendingNonceAt(ctx, url, e.MustGetFromAddress())
}

func EEAGetNonce(ctx context.Context, ec ethclient.EEAChainStateReader, e *tx.Envelope, url string) (uint64, error) {
	switch {
	case e.IsEeaSendPrivateTransactionPrivacyGroup():
		return ec.PrivNonce(ctx, url, e.MustGetFromAddress(), e.GetPrivacyGroupID())
	case e.IsEeaSendPrivateTransactionPrivateFor():
		return ec.PrivEEANonce(ctx, url, e.MustGetFromAddress(), e.GetPrivateFrom(), e.GetPrivateFor())
	default:
		return 0, errors.InternalError("invalid EEA envelope type")
	}
}
