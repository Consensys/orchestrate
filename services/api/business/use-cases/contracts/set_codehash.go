package contracts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	models2 "github.com/consensys/orchestrate/services/api/store/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const setCodeHashComponent = "use-cases.set-codehash"

type setCodeHashUseCase struct {
	agent  store.CodeHashAgent
	logger *log.Logger
}

func NewSetCodeHashUseCase(agent store.CodeHashAgent) usecases.SetContractCodeHashUseCase {
	return &setCodeHashUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(setCodeHashComponent),
	}
}

func (uc *setCodeHashUseCase) Execute(ctx context.Context, chainID string, address ethcommon.Address, codeHash hexutil.Bytes) error {
	ctx = log.WithFields(ctx, log.Field("chain_id", chainID), log.Field("address", chainID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("setting code-hash is starting ...")

	codehash := &models2.CodehashModel{
		ChainID:  chainID,
		Address:  address.Hex(),
		Codehash: codeHash.String(),
	}

	err := uc.agent.Insert(ctx, codehash)
	if err != nil {
		return errors.FromError(err).ExtendComponent(setCodeHashComponent)
	}

	logger.Debug("code-hash updated successfully")
	return nil
}
