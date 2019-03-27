package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Address return Address
func (acc *Account) Address() common.Address {
	if acc.GetAddr() == "" {
		return common.HexToAddress("")
	}

	if !common.IsHexAddress(acc.GetAddr()) {
		panic(fmt.Sprintf("%q is an invalid Ethereum address", acc.GetAddr()))
	}

	return common.HexToAddress(acc.GetAddr())
}

// SetAddress sets account address
func (acc *Account) SetAddress(addr common.Address) *Account {
	acc.Addr = addr.Hex()
	return acc
}
