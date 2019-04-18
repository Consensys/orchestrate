package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Address return Address
func (acc *Account) Address() (common.Address, error) {
	if acc.GetAddr() == "" {
		return common.HexToAddress(""), nil
	}

	if !common.IsHexAddress(acc.GetAddr()) {
		return common.HexToAddress(""), fmt.Errorf("%q is an invalid Ethereum address", acc.GetAddr())
	}

	return common.HexToAddress(acc.GetAddr()), nil
}

// SetAddress sets account address
func (acc *Account) SetAddress(addr common.Address) *Account {
	acc.Addr = addr.Hex()
	return acc
}
