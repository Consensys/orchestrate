package types

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Envelope holds information for a Faucet Credit Envelope
type Request struct {
	ChainID     *big.Int
	ChainURL    string
	ChainName   string
	ChainUUID   string
	Creditor    ethcommon.Address
	Beneficiary ethcommon.Address
	Amount      *big.Int
}
