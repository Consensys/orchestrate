package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeLog() *entities.Log {
	return &entities.Log{
		Status:  utils.StatusCreated,
		Message: "job message",
	}
}
