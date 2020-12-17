package accounts

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
)

const fundAccountComponent = "use-cases.fund-account"

type fundAccountUseCase struct {
	registryClient client.ChainRegistryClient
	sendTxUseCase  usecases.SendTxUseCase
}

func NewFundAccountUseCase(registryClient client.ChainRegistryClient, sendTxUseCase usecases.SendTxUseCase) usecases.FundAccountUseCase {
	return &fundAccountUseCase{registryClient: registryClient, sendTxUseCase: sendTxUseCase}
}

func (uc *fundAccountUseCase) Execute(ctx context.Context, account *entities.Account, chainName, tenantID string) error {
	logger := log.WithContext(ctx).WithField("address", account.Address)
	logger.Debug("funding account...")

	chain, err := uc.registryClient.GetChainByName(ctx, chainName)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	fct, err := uc.registryClient.GetFaucetCandidate(ctx, ethcommon.HexToAddress(account.Address), chain.UUID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil
		}

		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	txRequest := &entities.TxRequest{
		IdempotencyKey: utils.RandomString(16),
		ChainName:      chain.Name,
		Params: &entities.ETHTransactionParams{
			From:  fct.Creditor.Hex(),
			To:    account.Address,
			Value: fct.Amount.String(),
		},
		Labels:       chainregistry.FaucetToJobLabels(fct),
		InternalData: &entities.InternalData{},
	}
	_, err = uc.sendTxUseCase.Execute(ctx, txRequest, "", tenantID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(fundAccountComponent)
	}

	logger.WithField("faucet_uuid", fct.UUID).
		WithField("value", fct.Amount.String()).
		Info("account was topped successfully (funding transaction sent)")

	return nil
}
