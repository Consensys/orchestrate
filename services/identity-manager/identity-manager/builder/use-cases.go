package builder

import (
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
)

func NewUseCases(db store.DB) usecases.UseCases {
	return &useCases{}
}

type useCases struct{}
