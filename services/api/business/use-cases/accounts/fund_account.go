package accounts

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
)

const fundAccountComponent = "use-cases.fund-account"

type fundAccountUseCase struct {
	registryClient     client.ChainRegistryClient
	sendTxUseCase      usecases.SendTxUseCase
	getFaucetCandidate usecases.GetFaucetCandidateUseCase
}

func NewFundAccountUseCase(
	registryClient client.ChainRegistryClient,
	sendTxUseCase usecases.SendTxUseCase,
	getFaucetCandidate usecases.GetFaucetCandidateUseCase,
) usecases.FundAccountUseCase {
	return &fundAccountUseCase{
		registryClient:     registryClient,
		sendTxUseCase:      sendTxUseCase,
		getFaucetCandidate: getFaucetCandidate,
	}
}

func (uc *fundAccountUseCase) Execute(ctx context.Context, account *entities.Account, chainName, tenantID string) error {
	logger := log.WithContext(ctx).WithField("address", account.Address)
	logger.Debug("funding account...")

	chain, err := uc.registryClient.GetChainByName(ctx, chainName)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	faucet, err := uc.getFaucetCandidate.Execute(ctx, account.Address, chain, []string{tenantID, multitenancy.DefaultTenant})
	if err != nil {
		if errors.IsNotFoundError(err) {
			logger.Debug("unnecessary funding, skipping top-up")
			return nil
		}

		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	txRequest := &entities.TxRequest{
		IdempotencyKey: utils.RandomString(16),
		ChainName:      chain.Name,
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

	logger.WithField("faucet_uuid", faucet.UUID).
		WithField("value", faucet.Amount).
		Info("account was topped successfully (funding transaction sent)")

	return nil
}
