package faucet

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func computeKey(chainID *big.Int, a common.Address) string {
	return fmt.Sprintf("%v-%v", chainID.Text(16), a.Hex())
}
