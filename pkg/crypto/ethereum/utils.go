package ethereum

import (
	"encoding/base64"
	"fmt"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func GetSignatureSender(signature, payload string) (*ethcommon.Address, error) {
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %s", err.Error())
	}

	hash := crypto.Keccak256([]byte(payload))
	pubKey, err := crypto.SigToPub(hash, signatureBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %s", err.Error())
	}

	address := crypto.PubkeyToAddress(*pubKey)
	return &address, nil
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

func GetEIP155Signer(chainID string) types.Signer {
	chainIDBigInt := new(big.Int)
	chainIDBigInt, _ = chainIDBigInt.SetString(chainID, 10)
	return types.NewEIP155Signer(chainIDBigInt)
}

func GetQuorumPrivateTxSigner() quorumtypes.Signer {
	return quorumtypes.QuorumPrivateTxSigner{}
}
