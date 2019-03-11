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
