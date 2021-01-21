package faucets

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getFaucetComponent = "use-cases.get-faucet"

// getFaucetUseCase is a use case to get a faucet
type getFaucetUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewGetFaucetUseCase creates a new GetFaucetUseCase
func NewGetFaucetUseCase(db store.DB) usecases.GetFaucetUseCase {
	return &getFaucetUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(getFaucetComponent),
	}
}

// Execute gets a faucet
func (uc *getFaucetUseCase) Execute(ctx context.Context, uuid string, tenants []string) (*entities.Faucet, error) {
	ctx = log.WithFields(ctx, log.Field("faucet", uuid))
	logger := uc.logger.WithContext(ctx)

	faucetModel, err := uc.db.Faucet().FindOneByUUID(ctx, uuid, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getFaucetComponent)
	}

	logger.Debug("faucet found successfully")
	return parsers.NewFaucetFromModel(faucetModel), nil
}
