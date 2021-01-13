package faucets

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const registerFaucetComponent = "use-cases.register-faucet"

// registerFaucetUseCase is a use case to register a new faucet
type registerFaucetUseCase struct {
	db             store.DB
	searchFaucetUC usecases.SearchFaucetsUseCase
}

// NewRegisterFaucetUseCase creates a new RegisterFaucetUseCase
func NewRegisterFaucetUseCase(db store.DB, searchFaucetUC usecases.SearchFaucetsUseCase) usecases.RegisterFaucetUseCase {
	return &registerFaucetUseCase{
		db:             db,
		searchFaucetUC: searchFaucetUC,
	}
}

// Execute registers a new faucet
func (uc *registerFaucetUseCase) Execute(ctx context.Context, faucet *entities.Faucet) (*entities.Faucet, error) {
	logger := log.WithContext(ctx).
		WithField("name", faucet.Name).
		WithField("chain_rule", faucet.ChainRule).
		WithField("tenant", faucet.TenantID).
		WithField("creditor_account", faucet.CreditorAccount)
	logger.Debug("registering new faucet")

	faucetsRetrieved, err := uc.searchFaucetUC.Execute(ctx, &entities.FaucetFilters{
		Names: []string{faucet.Name},
	}, []string{faucet.TenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(registerFaucetComponent)
	}

	if len(faucetsRetrieved) > 0 {
		errMessage := "faucet with same name already exists"
		logger.Error(errMessage)
		return nil, errors.AlreadyExistsError(errMessage).ExtendComponent(registerFaucetComponent)
	}

	faucetModel := parsers.NewFaucetModelFromEntity(faucet)
	err = uc.db.Faucet().Insert(ctx, faucetModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(registerFaucetComponent)
	}

	logger.WithField("faucet_uuid", faucetModel.UUID).Info("faucet registered successfully")
	return parsers.NewFaucetFromModel(faucetModel), nil
}
