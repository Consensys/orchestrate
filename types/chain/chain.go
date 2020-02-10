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
		ChainId: big.NewInt(chainID).String(),
	}
}

// FromBigInt create a new chain id from a big integer value
func FromBigInt(chainID *big.Int) *Chain {
	return &Chain{
		ChainId: chainID.String(),
	}
}

// UUID return chain UUID in big.Int format
func (c *Chain) GetBigChainID() *big.Int {
	chainID, _ := big.NewInt(0).SetString(c.GetChainId(), 10)
	return chainID
}

// SetChainID set chain UUID
func (c *Chain) SetChainID(chainID *big.Int) *Chain {
	c.ChainId = chainID.String()
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
