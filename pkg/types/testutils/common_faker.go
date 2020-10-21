package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
)

// FakeAccountInstance returns a new fake account
func FakeAccountInstance() *common.AccountInstance {
	return &common.AccountInstance{
		Account: "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		ChainId: "chainId",
	}
}
