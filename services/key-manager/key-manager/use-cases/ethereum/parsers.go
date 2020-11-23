package ethereum

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func NewECDSAFromPrivKey(privKey string) (*ecdsa.PrivateKey, error) {
	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		errMessage := "failed to parse secp256k1 private key"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage).ExtendComponent(signTransactionComponent)
	}

	return ecdsaPrivKey, nil
}
