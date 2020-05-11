package session

import (
	"encoding/base64"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/types"
	"golang.org/x/crypto/sha3"
)

func privateTxHash(tx *ethtypes.Transaction, privateArgs *types.PrivateArgs, chain *big.Int) (ethcommon.Hash, error) {
	privateFromEncoded, privateRecipientEncoded, err := privateArgsEncoded(privateArgs)
	if err != nil {
		return ethcommon.Hash{}, err
	}

	return rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		chain,
		uint(0),
		uint(0),
		privateFromEncoded,
		privateRecipientEncoded,
		privateArgs.PrivateTxType,
	})
}

func privateArgsEncoded(privateArgs *types.PrivateArgs) (privateFromEncoded, privateRecipientEncoded interface{}, err error) {
	if len(privateArgs.PrivateFor) > 0 && privateArgs.PrivacyGroupID != "" {
		return nil, nil, errors.DataError("privacyGroupId and privateFor fields are mutually exclusive")
	}

	privateFromEncoded, err = base64.StdEncoding.DecodeString(privateArgs.PrivateFrom)
	if err != nil {
		return nil, nil, errors.DataError("invalid base64 for privateFrom - got %v", err)
	}

	if privateArgs.PrivacyGroupID != "" {
		privateRecipientEncoded, err = base64.StdEncoding.DecodeString(privateArgs.PrivacyGroupID)
		if err != nil {
			return nil, nil, errors.DataError("invalid base64 for privacyGroupId - got %v", err)
		}
	} else {
		var privateForByteSlice [][]byte
		for _, v := range privateArgs.PrivateFor {
			b, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return nil, nil, errors.DataError("invalid base64 for privateFor - got %v", err)
			}
			privateForByteSlice = append(privateForByteSlice, b)
		}
		privateRecipientEncoded = privateForByteSlice
	}

	return privateFromEncoded, privateRecipientEncoded, nil
}

func encodePrivateTx(tx *ethtypes.Transaction, privateArgs *types.PrivateArgs) []byte {
	v, r, s := tx.RawSignatureValues()
	privateFromEncoded, privateRecipientEncoded, _ := privateArgsEncoded(privateArgs)

	rplEncoding, _ := rlpEncode([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		v,
		r,
		s,
		privateFromEncoded,
		privateRecipientEncoded,
		privateArgs.PrivateTxType,
	})
	return rplEncoding
}

func calculatePrivateV(v *big.Int) *big.Int {
	if v.Int64() == 27 {
		return big.NewInt(37)
	}
	return big.NewInt(38)
}

func rlpHash(object interface{}) (hash ethcommon.Hash, err error) {
	hashAlgo := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hashAlgo, object)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	hashAlgo.Sum(hash[:0])
	return hash, nil
}

func rlpEncode(object interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(object)
}
