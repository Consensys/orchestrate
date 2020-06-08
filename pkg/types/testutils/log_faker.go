package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
)

func FakeLog() *types.Log {
	return &types.Log{
		Status:  types.StatusCreated,
		Message: "job message",
	}
}
