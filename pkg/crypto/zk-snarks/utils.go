package zksnarks

import (
	"crypto/sha256"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
	"github.com/consensys/quorum/common/hexutil"
)

func VerifyZKSMessage(publicKey, signature string, message []byte) (bool, error) {
	pubKey := eddsa.PublicKey{}
	pubKeyB, err := hexutil.Decode(publicKey)
	if err != nil {
		return false, errors.EncodingError("failed to decode public key").AppendReason(err.Error())
	}

	_, err = pubKey.SetBytes(pubKeyB)
	if err != nil {
		return false, errors.InvalidParameterError("invalid public key value").AppendReason(err.Error())
	}

	signB, err := hexutil.Decode(signature)
	if err != nil {
		return false, errors.EncodingError("failed to decode signature").AppendReason(err.Error())
	}

	verified, err := pubKey.Verify(signB, message, sha256.New())
	if err != nil {
		return false, errors.InvalidParameterError("invalid verification values").AppendReason(err.Error())
	}

	return verified, nil
}
