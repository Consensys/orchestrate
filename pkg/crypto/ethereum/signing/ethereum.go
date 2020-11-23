package signing

import (
	"crypto/ecdsa"
	"encoding/base64"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func GetEIP155Signer(chainID string) types.Signer {
	chainIDBigInt := new(big.Int)
	chainIDBigInt, _ = chainIDBigInt.SetString(chainID, 10)
	return types.NewEIP155Signer(chainIDBigInt)
}

func GetQuorumPrivateTxSigner() quorumtypes.Signer {
	return quorumtypes.QuorumPrivateTxSigner{}
}

func SignTransaction(tx *types.Transaction, privKey *ecdsa.PrivateKey, signer types.Signer) ([]byte, error) {
	h := signer.Hash(tx)
	decodedSignature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		errMessage := "failed to sign ethereum transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return decodedSignature, nil
}

func SignQuorumPrivateTransaction(tx *quorumtypes.Transaction, privKey *ecdsa.PrivateKey, signer quorumtypes.Signer) ([]byte, error) {
	h := signer.Hash(tx)
	decodedSignature, err := crypto.Sign(h[:], privKey)
	if err != nil {
		errMessage := "failed to sign quorum private transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return decodedSignature, nil
}

func SignEEATransaction(tx *types.Transaction, privateArgs *entities.PrivateETHTransactionParams, chainID string, privKey *ecdsa.PrivateKey) ([]byte, error) {
	chainIDBigInt, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		errMessage := "invalid chainID"
		log.WithField("chain_id", chainID).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	privateFromEncoded, err := GetEncodedPrivateFrom(privateArgs.PrivateFrom)
	if err != nil {
		return nil, err
	}

	privateRecipientEncoded, err := GetEncodedPrivateRecipient(privateArgs.PrivacyGroupID, privateArgs.PrivateFor)
	if err != nil {
		return nil, err
	}

	hash, err := rlp.Hash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		chainIDBigInt,
		uint(0),
		uint(0),
		privateFromEncoded,
		privateRecipientEncoded,
		privateArgs.PrivateTxType,
	})
	if err != nil {
		errMessage := "failed to hash eea transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	signature, err := crypto.Sign(hash[:], privKey)
	if err != nil {
		errMessage := "failed to sign eea transaction"
		log.WithError(err).Error(errMessage)
		return nil, errors.CryptoOperationError(errMessage)
	}

	return signature, err
}

func GetEncodedPrivateFrom(privateFrom string) ([]byte, error) {
	privateFromEncoded, err := base64.StdEncoding.DecodeString(privateFrom)
	if err != nil {
		errMessage := "invalid base64 privateFrom"
		log.WithError(err).WithField("private_from", privateFrom).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	return privateFromEncoded, nil
}

func GetEncodedPrivateRecipient(privacyGroupID string, privateFor []string) (interface{}, error) {
	var privateRecipientEncoded interface{}
	var err error
	if privacyGroupID != "" {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(privacyGroupID)
		if err != nil {
			errMessage := "invalid base64 privacyGroupId"
			log.WithError(err).WithField("privacy_group_id", privacyGroupID).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range privateFor {
			b, der := base64.StdEncoding.DecodeString(v)
			if der != nil {
				errMessage := "invalid base64 privateFor"
				log.WithError(der).WithField("Private_for", v).Error(errMessage)
				return nil, errors.InvalidParameterError(errMessage)
			}
			privateForByteSlice = append(privateForByteSlice, b)
		}
		privateRecipientEncoded = privateForByteSlice
	}

	return privateRecipientEncoded, nil
}
