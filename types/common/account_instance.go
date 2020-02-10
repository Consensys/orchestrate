package common

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Short returns a string representation of the account instance
func (i *AccountInstance) Short() string {
	var addr, id string
	if i.GetAccount() == "nil" {
		addr = ""
	} else {
		addr = i.GetAccount()
	}
	if i.GetChain() == nil {
		id = ""
	} else {
		id = i.GetChain().GetChainId()
	}
	return fmt.Sprintf("%v@%v", addr, id)
}

func (i *AccountInstance) GetAccountAddress() ethcommon.Address {
	return ethcommon.HexToAddress(i.GetAccount())
}
