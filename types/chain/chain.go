package chain

import (
	"math/big"
)

// FromString creates a new chain id from an string value
func FromString(chainID string) *Chain {
	id, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		panic("Cannot parse string to chain")
	}
	return FromBigInt(id)
}

// FromInt creates a new chain id from an integer value
func FromInt(chainID int64) *Chain {
	return &Chain{
		Id: big.NewInt(chainID).Bytes(),
	}
}

// FromBigInt create a new chain id from a big integer value
func FromBigInt(chainID *big.Int) *Chain {
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
