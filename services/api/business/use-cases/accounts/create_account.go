package accounts

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
)

const createAccountComponent = "use-cases.create-account"

type createAccountUseCase struct {
	db               store.DB
	searchUC         usecases.SearchAccountsUseCase
	fundAccountUC    usecases.FundAccountUseCase
	keyManagerClient client.KeyManagerClient
}

func NewCreateAccountUseCase(db store.DB, searchUC usecases.SearchAccountsUseCase, fundAccountUC usecases.FundAccountUseCase,
	keyManagerClient client.KeyManagerClient) usecases.CreateAccountUseCase {
	return &createAccountUseCase{
		db:               db,
		searchUC:         searchUC,
		keyManagerClient: keyManagerClient,
		fundAccountUC:    fundAccountUC,
	}
}

func (uc *createAccountUseCase) Execute(ctx context.Context, account *entities.Account, privateKey, chainName, tenantID string) (*entities.Account, error) {
	logger := log.WithContext(ctx).WithField("alias", account.Alias).WithField("chain", chainName)
	logger.Debug("creating new Ethereum account")

	accounts, err := uc.searchUC.Execute(ctx, &entities.AccountFilters{Aliases: []string{account.Alias}}, []string{tenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if len(accounts) > 0 {
		errMsg := fmt.Sprintf("alias %s already exists", account.Alias)
		logger.Error(errMsg)
		return nil, errors.AlreadyExistsError(errMsg)
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
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	account.Address = resp.Address
	account.PublicKey = resp.PublicKey
	account.CompressedPublicKey = resp.CompressedPublicKey
	account.TenantID = tenantID

	// IMPORTANT: Addresses are unique across every tenant
	_, err = uc.db.Account().FindOneByAddress(ctx, account.Address, []string{})
	if err == nil {
		errMsg := fmt.Sprintf("account %s already exists", account.Address)
		logger.Error(errMsg)
		return nil, errors.AlreadyExistsError(errMsg).ExtendComponent(createAccountComponent)
	} else if !errors.IsNotFoundError(err) {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	accountModel := parsers.NewAccountModelFromEntities(account)
	err = uc.db.Account().Insert(ctx, accountModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if chainName != "" {
		err = uc.fundAccountUC.Execute(ctx, account, chainName, tenantID)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
		}
	}

	logger.WithField("address", account.Address).Info("account created successfully")
	return parsers.NewAccountEntityFromModels(accountModel), nil
}
