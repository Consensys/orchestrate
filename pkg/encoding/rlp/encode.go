package rlp

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func Hash(object interface{}) (hash ethcommon.Hash, err error) {
	hashAlgo := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hashAlgo, object)
	if err != nil {
		return ethcommon.Hash{}, err
	}
	hashAlgo.Sum(hash[:0])
	return hash, nil
}

func Encode(object interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(object)
}
