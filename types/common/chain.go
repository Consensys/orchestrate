package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ID return chain ID in big.Int format
func (chain *Chain) ID() *big.Int {
	if chain.Id == "" {
		return big.NewInt(0)
	}
	return hexutil.MustDecodeBig(chain.Id)
}

// SetID set chain ID
func (chain *Chain) SetID(id *big.Int) *Chain {
	chain.Id = hexutil.EncodeBig(id)
	return chain
}
