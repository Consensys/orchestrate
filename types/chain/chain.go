package chain

import (
	"math/big"
)

// CreateChainInt creates a new chain id from an integer value
func CreateChainInt(chainID int64) *Chain {
	return &Chain{
		Id: big.NewInt(chainID).Bytes(),
	}
}

// CreateChainBigInt create a new chain id from a big integer value
func CreateChainBigInt(chainID *big.Int) *Chain {
	return &Chain{
		Id: chainID.Bytes(),
	}
}

// ID return chain ID in big.Int format
func (c *Chain) ID() *big.Int {
	return big.NewInt(0).SetBytes(c.GetId())
}

// SetID set chain ID
func (c *Chain) SetID(id *big.Int) *Chain {
	c.Id = id.Bytes()
	return c
}
