package contracts

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	models2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const setCodeHashComponent = "use-cases.set-codehash"

type setCodeHashUseCase struct {
	agent store.CodeHashAgent
}

func NewSetCodeHashUseCase(agent store.CodeHashAgent) usecases.SetContractCodeHashUseCase {
	return &setCodeHashUseCase{
		agent: agent,
	}
}

func (usecase *setCodeHashUseCase) Execute(ctx context.Context, chainID, address, codeHash string) error {
	logger := log.WithContext(ctx).WithField("chainID", chainID).WithField("address", address).
		WithField("code_hash", utils.ShortString(codeHash, 10))
	logger.Debug("setting CodeHash is starting ...")

	codehash := &models2.CodehashModel{
		ChainID:  chainID,
		Address:  address,
		Codehash: codeHash,
	}

	err := usecase.agent.Insert(ctx, codehash)
	if err != nil {
		return errors.FromError(err).ExtendComponent(setCodeHashComponent)
	}

	logger.Debug("code hash was set successfully")
	return nil
}
