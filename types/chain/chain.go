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
		ChainId: big.NewInt(chainID).Bytes(),
	}
}

// FromBigInt create a new chain id from a big integer value
func FromBigInt(chainID *big.Int) *Chain {
	return &Chain{
		ChainId: chainID.Bytes(),
	}
}

// UUID return chain UUID in big.Int format
func (c *Chain) GetBigChainID() *big.Int {
	return big.NewInt(0).SetBytes(c.GetChainId())
}

// SetChainID set chain UUID
func (c *Chain) SetChainID(id *big.Int) *Chain {
	c.ChainId = id.Bytes()
	return c
}

// SetUUID set chain uuid of the chain registry
func (c *Chain) SetUUID(chainUUID string) *Chain {
	c.Uuid = chainUUID
	return c
}

// SetName set chain name of the chain registry
func (c *Chain) SetName(chainName string) *Chain {
	c.Name = chainName
	return c
}
