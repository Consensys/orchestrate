package zksnarks

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/consensys/quorum/common/hexutil"
)

func VerifyZKSMessage(publicKey, signature, message string) (bool, error) {
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

	messageB, err := hexutil.Decode(message)
	if err != nil {
		return false, errors.EncodingError("failed to decode message").AppendReason(err.Error())
	}

	verified, err := pubKey.Verify(signB, messageB, hash.MIMC_BN254.New("seed"))
	if err != nil {
		return false, errors.InvalidParameterError("invalid verification values").AppendReason(err.Error())
	}

	return verified, nil
}
