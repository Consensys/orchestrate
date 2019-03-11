package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	acc := &Account{
		Id:   "abcd",
		Addr: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
	}

	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", acc.Address().Hex(), "Address should match")
}
