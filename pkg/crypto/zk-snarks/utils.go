package zksnarks

import (
	"crypto/sha256"

	eddsa "github.com/consensys/gnark/crypto/signature/eddsa/bn256"
	"github.com/consensys/quorum/common/hexutil"
)

func VerifyZKSMessage(publicKey, signature string, message []byte) (bool, error) {
	pubKey := eddsa.PublicKey{}
	pubKeyB, err := hexutil.Decode(publicKey)
	if err != nil {
		return false, err
	}

	_, err = pubKey.SetBytes(pubKeyB)
	if err != nil {
		return false, err
	}

	signB, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	return pubKey.Verify(signB, message, sha256.New())
}
