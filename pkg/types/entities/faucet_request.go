package entities

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type FaucetRequest struct {
	Chain       *Chain
	Beneficiary ethcommon.Address
	Candidates  map[string]*Faucet
}
