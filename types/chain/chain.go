package chain

import (
	"log"
	"math/big"
)

// CreateChainString creates a new chain id from an string value
func CreateChainString(chainID string) *Chain {
	id := new(big.Int)
	id, ok := id.SetString(chainID, 10)
	if !ok {
		log.Fatalf("Cannot cast string to big.Int: %v", chainID)
	}
	return CreateChainBigInt(id)
}

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
