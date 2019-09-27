package types

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Request holds information for a Faucet Credit Request
type Request struct {
	ChainID     *big.Int
	Creditor    ethcommon.Address
	Beneficiary ethcommon.Address
	Amount      *big.Int
}
