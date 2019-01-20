package faucet

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestComputeKey(t *testing.T) {
	key := computeKey(big.NewInt(1234), common.HexToAddress("0xabcdef"))
	expected := "4d2-0x0000000000000000000000000000000000abcDeF"
	if key != expected {
		t.Errorf("computeKey expected %v but got %v", expected, key)
	}
}
