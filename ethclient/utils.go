package ethclient

import (
	"math/big"
)

// ChainIDToString transform a chain ID from big.Int to string
// TODO: to be moved in pkg
func ChainIDToString(chainID *big.Int) string {
	return chainID.Text(16)
}
