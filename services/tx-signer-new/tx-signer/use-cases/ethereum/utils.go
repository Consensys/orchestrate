package ethereum

import (
	"github.com/consensys/quorum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func GetSignedRawTransaction(transaction *types.Transaction, signature []byte, signer types.Signer) (string, error) {
	signedTx, err := transaction.WithSignature(signer, signature)
	if err != nil {
		errMessage := "failed to set transaction signature"
		log.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage).ExtendComponent(signEEATransactionComponent)
	}

	signedRaw, err := rlp.Encode(signedTx)
	if err != nil {
		errMessage := "failed to RLP encode signed transaction"
		log.WithError(err).Error(errMessage)
		return "", errors.CryptoOperationError(errMessage).ExtendComponent(signEEATransactionComponent)
	}

	return hexutil.Encode(signedRaw), nil
}
