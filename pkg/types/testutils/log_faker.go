package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

func FakeLog() *entities.Log {
	return &entities.Log{
		Status:  entities.StatusCreated,
		Message: "job message",
	}
}
