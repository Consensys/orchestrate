package identity

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	client3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
)

const fundingIdentityComponent = "use-cases.funding-identity"

type fundingIdentityUseCase struct {
	registryClient    client2.ChainRegistryClient
	txSchedulerClient client3.TransactionSchedulerClient
}

func NewFundingIdentityUseCase(registryClient client2.ChainRegistryClient, txSchedulerClient client3.TransactionSchedulerClient) usecases.FundingIdentityUseCase {
	return &fundingIdentityUseCase{
		registryClient:    registryClient,
		txSchedulerClient: txSchedulerClient,
	}
}

func (uc *fundingIdentityUseCase) Execute(ctx context.Context, identity *entities.Identity, chainName string) error {
	logger := log.WithContext(ctx).
		WithField("alias", identity.Alias)

	logger.Debug("creating new identity...")

	chain, err := uc.registryClient.GetChainByName(ctx, chainName)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundingIdentityComponent)
	}

	fct, err := uc.registryClient.GetFaucetCandidate(ctx, ethcommon.HexToAddress(identity.Address), chain.UUID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil
		}

		return errors.FromError(err).ExtendComponent(fundingIdentityComponent)
	}

	_, err = uc.txSchedulerClient.SendTransferTransaction(ctx,
		&txscheduler.TransferRequest{
			ChainName: chain.Name,
			Params: txscheduler.TransferParams{
				From:  fct.Creditor.Hex(),
				To:    identity.PublicKey,
				Value: fct.Amount.String(),
			},
			Labels: chainregistry.FaucetToJobLabels(fct),
		})

	if err != nil {
		return errors.FromError(err).ExtendComponent(fundingIdentityComponent)
	}

	logger.WithField("faucet", fct.UUID).WithField("value", fct.Amount.String()).
		Info("funding transaction was sent")

	return nil
}
