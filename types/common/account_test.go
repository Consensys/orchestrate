package common

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	acc := &Account{
		Id:   "abcd",
		Addr: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
	}
	address, _ := acc.Address()
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", address.Hex(), "Address should match")

	addr := common.HexToAddress("0xc99a171AA7365FA16E52e737c24CD78E4aA8c7F5")
	acc = acc.SetAddress(addr)
	assert.Equal(t, "0xc99a171AA7365FA16E52e737c24CD78E4aA8c7F5", acc.Addr, "SetAddress should set the address")
}
