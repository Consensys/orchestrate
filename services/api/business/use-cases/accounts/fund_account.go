package accounts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"

	"github.com/consensys/orchestrate/pkg/utils"

	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

const fundAccountComponent = "use-cases.fund-account"

type fundAccountUseCase struct {
	searchChainsUC     usecases.SearchChainsUseCase
	sendTxUseCase      usecases.SendTxUseCase
	getFaucetCandidate usecases.GetFaucetCandidateUseCase
	logger             *log.Logger
}

func NewFundAccountUseCase(
	searchChainsUC usecases.SearchChainsUseCase,
	sendTxUseCase usecases.SendTxUseCase,
	getFaucetCandidate usecases.GetFaucetCandidateUseCase,
) usecases.FundAccountUseCase {
	return &fundAccountUseCase{
		searchChainsUC:     searchChainsUC,
		sendTxUseCase:      sendTxUseCase,
		getFaucetCandidate: getFaucetCandidate,
		logger:             log.NewLogger().SetComponent(fundAccountComponent),
	}
}

func (uc *fundAccountUseCase) Execute(ctx context.Context, account *entities.Account, chainName, tenantID string) error {
	ctx = log.WithFields(ctx, log.Field("address", account.Address))
	logger := uc.logger.WithContext(ctx)

	tenants := []string{tenantID, multitenancy.DefaultTenant}

	chains, err := uc.searchChainsUC.Execute(ctx, &entities.ChainFilters{Names: []string{chainName}}, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	if len(chains) == 0 {
		errMsg := "chain does not exist"
		logger.Warn(errMsg)
		return errors.InvalidParameterError(errMsg).ExtendComponent(fundAccountComponent)
	}

	faucet, err := uc.getFaucetCandidate.Execute(ctx, account.Address, chains[0], tenants)
	if err != nil {
		if errors.IsNotFoundError(err) {
			logger.Debug("unnecessary funding, skipping top-up")
			return nil
		}

		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	txRequest := &entities.TxRequest{
		IdempotencyKey: utils.RandString(16),
		ChainName:      chains[0].Name,
		Params: &entities.ETHTransactionParams{
			From:  faucet.CreditorAccount,
			To:    account.Address,
			Value: faucet.Amount,
		},
		Labels: map[string]string{
			"faucetUUID": faucet.UUID,
		},
		InternalData: &entities.InternalData{},
	}
	_, err = uc.sendTxUseCase.Execute(ctx, txRequest, "", tenantID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	logger.WithField("faucet", faucet.UUID).WithField("value", faucet.Amount).
		Debug("account was topped successfully (funding transaction sent)")

	return nil
}
