package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EthCrediter is an interface for crediting an account with ether
type Faucet interface {
	// Credit should credit an account based on its own set of security rules
	// If credit is successful it should return amount credited
	// Credit should respond synchronously (not wait for a credit transaction to be mined)
	Credit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool, error)
}
