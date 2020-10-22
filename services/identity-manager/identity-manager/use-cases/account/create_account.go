package account

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

const createAccountComponent = "use-cases.create-account"

type createAccountUseCase struct {
	db               store.DB
	searchUC         usecases.SearchAccountsUseCase
	fundingAccountUC usecases.FundingAccountUseCase
	keyManagerClient client.KeyManagerClient
}

func NewCreateAccountUseCase(db store.DB, searchUC usecases.SearchAccountsUseCase, fundingAccountUC usecases.FundingAccountUseCase,
	keyManagerClient client.KeyManagerClient) usecases.CreateAccountUseCase {
	return &createAccountUseCase{
		db:               db,
		searchUC:         searchUC,
		keyManagerClient: keyManagerClient,
		fundingAccountUC: fundingAccountUC,
	}
}

func (uc *createAccountUseCase) Execute(ctx context.Context, account *entities.Account, privateKey, chainName, tenantID string) (*entities.Account, error) {
	logger := log.WithContext(ctx).WithField("alias", account.Alias).WithField("chain", "chainName")

	logger.Debug("creating new account...")
	idens, err := uc.searchUC.Execute(ctx, &entities.AccountFilters{Aliases: []string{account.Alias}}, []string{tenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if len(idens) > 0 {
		errMsg := fmt.Sprintf("alias %s already exists", account.Alias)
		logger.Error(errMsg)
		return nil, errors.InvalidParameterError(errMsg)
	}

	// REMINDER: For now, Account API only support ETH accounts
	var resp *types.ETHAccountResponse
	if privateKey != "" {
		resp, err = uc.keyManagerClient.ETHImportAccount(ctx, &types.ImportETHAccountRequest{
			Namespace:  tenantID,
			PrivateKey: privateKey,
		})
	} else {
		resp, err = uc.keyManagerClient.ETHCreateAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		})
	}

	if err != nil {
		return nil, err
	}

	account.Address = resp.Address
	account.PublicKey = resp.PublicKey
	account.CompressedPublicKey = resp.CompressedPublicKey

	accountModel := parsers.NewAccountModelFromEntities(account)
	accountModel.TenantID = tenantID
	err = uc.db.Account().Insert(ctx, accountModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if chainName != "" {
		err = uc.fundingAccountUC.Execute(ctx, account, chainName)
		if err != nil {
			logger.WithError(err).Error("cannot trigger funding account")
			return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
		}
	}

	logger.WithField("address", account.Address).Info("account was created successfully")

	return parsers.NewAccountEntityFromModels(accountModel), nil
}
