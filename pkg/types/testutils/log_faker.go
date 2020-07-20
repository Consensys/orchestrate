package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeLog() *types.Log {
	return &types.Log{
		Status:  utils.StatusCreated,
		Message: "job message",
	}
}
