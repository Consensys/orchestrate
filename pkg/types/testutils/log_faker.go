package testutils

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
)

func FakeLog() *entities.Log {
	return &entities.Log{
		Status:  entities.StatusCreated,
		Message: "job message",
	}
}
