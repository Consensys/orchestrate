package common

import (
	"fmt"
	"math/big"
)

// ID return chain ID in big.Int format
func (chain *Chain) ID() *big.Int {
	if chain.Id == "" {
		return big.NewInt(0)
	}

	chainID, ok := big.NewInt(0).SetString(chain.Id, 10)
	if !ok {
		panic(fmt.Sprintf("invalid decimal chain ID %q", chain.Id))
	}

	return chainID
}

// SetID set chain ID
func (chain *Chain) SetID(id *big.Int) *Chain {
	chain.Id = id.Text(10)
	return chain
}
