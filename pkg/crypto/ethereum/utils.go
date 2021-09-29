package ethereum

import (
	"encoding/base64"
	"math/big"

	"github.com/consensys/orchestrate/pkg/errors"
	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetSignatureSender(signature, payload string) (*ethcommon.Address, error) {
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return nil, errors.EncodingError("failed to decode signature").AppendReason(err.Error())
	}

	payloadBytes, err := hexutil.Decode(payload)
	if err != nil {
		return nil, errors.EncodingError("failed to decode payload").AppendReason(err.Error())
	}

	pubKey, err := crypto.SigToPub(crypto.Keccak256(payloadBytes), signatureBytes)
	if err != nil {
		return nil, errors.CryptoOperationError("failed to recover public key").AppendReason(err.Error())
	}

	address := crypto.PubkeyToAddress(*pubKey)
	return &address, nil
}

func GetEncodedPrivateFrom(privateFrom string) ([]byte, error) {
	privateFromEncoded, err := base64.StdEncoding.DecodeString(privateFrom)
	if err != nil {
		return nil, errors.EncodingError("invalid base64 value for 'privateFrom'").AppendReason(err.Error())
	}

	return privateFromEncoded, nil
}

func GetEncodedPrivateRecipient(privacyGroupID string, privateFor []string) (interface{}, error) {
	var privateRecipientEncoded interface{}
	var err error
	if privacyGroupID != "" {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(privacyGroupID)
		if err != nil {
			return nil, errors.EncodingError("invalid base64 value for 'privacyGroupId'").AppendReason(err.Error())
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range privateFor {
			b, der := base64.StdEncoding.DecodeString(v)
			if der != nil {
				return nil, errors.EncodingError("invalid base64 value for 'privateFor'").AppendReason(der.Error())
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
