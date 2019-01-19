package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// FaucetRequest holds information for a faucet request
type FaucetRequest struct {
	ChainID *big.Int
	Address common.Address
	Value   *big.Int
}

// Faucet is an interface for crediting an account with ether
type Faucet interface {
	// Credit should credit an account based on its own set of security rules
	// If credit is successful it should return amount credited and true
	// Credit should respond synchronously (not wait for a credit transaction to be mined)
	Credit(r *FaucetRequest) (*big.Int, bool, error)
}
