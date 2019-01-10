package services

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EthCrediter is an interface for crediting an account with ether
type EthCrediter interface {
	Credit(chainID *big.Int, a common.Address, value *big.Int) error
}

// EthCreditController is an interface to control if a credit should append
type EthCreditController interface {
	ShouldCredit(chainID *big.Int, a common.Address, value *big.Int) (*big.Int, bool)
}
