package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
)

// FakeAccount returns a new fake account
func FakeAccount() *common.AccountInstance {
	return &common.AccountInstance{
		Account: "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		ChainId: "chainId",
	}
}
