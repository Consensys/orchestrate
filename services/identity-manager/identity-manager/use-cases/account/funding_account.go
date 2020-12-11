package account

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases"
	client3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
)

const fundingAccountComponent = "use-cases.funding-account"

type fundingAccountUseCase struct {
	registryClient    client2.ChainRegistryClient
	txSchedulerClient client3.TransactionSchedulerClient
}

func NewFundingAccountUseCase(registryClient client2.ChainRegistryClient, txSchedulerClient client3.TransactionSchedulerClient) usecases.FundingAccountUseCase {
	return &fundingAccountUseCase{
		registryClient:    registryClient,
		txSchedulerClient: txSchedulerClient,
	}
}

func (uc *fundingAccountUseCase) Execute(ctx context.Context, account *entities.Account, chainName string) error {
	logger := log.WithContext(ctx).
		WithField("alias", account.Alias)

	logger.Debug("creating new account...")

	chain, err := uc.registryClient.GetChainByName(ctx, chainName)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundingAccountComponent)
	}

	fct, err := uc.registryClient.GetFaucetCandidate(ctx, ethcommon.HexToAddress(account.Address), chain.UUID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil
		}

		return errors.FromError(err).ExtendComponent(fundingAccountComponent)
	}

	_, err = uc.txSchedulerClient.SendTransferTransaction(ctx,
		&txscheduler.TransferRequest{
			ChainName: chain.Name,
			Params: txscheduler.TransferParams{
				From:  fct.Creditor.Hex(),
				To:    account.Address,
				Value: fct.Amount.String(),
			},
			Labels: chainregistry.FaucetToJobLabels(fct),
		})
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundingAccountComponent)
	}

	logger.WithField("faucet", fct.UUID).WithField("address", account.Address).
		WithField("value", fct.Amount.String()).Info("funding transaction was sent")

	return nil
}
