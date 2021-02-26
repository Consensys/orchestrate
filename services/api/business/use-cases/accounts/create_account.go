package accounts

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/ConsenSys/orchestrate/services/key-manager/client"
)

const createAccountComponent = "use-cases.create-account"

type createAccountUseCase struct {
	db               store.DB
	searchUC         usecases.SearchAccountsUseCase
	fundAccountUC    usecases.FundAccountUseCase
	keyManagerClient client.KeyManagerClient
	logger           *log.Logger
}

func NewCreateAccountUseCase(db store.DB, searchUC usecases.SearchAccountsUseCase, fundAccountUC usecases.FundAccountUseCase,
	keyManagerClient client.KeyManagerClient) usecases.CreateAccountUseCase {
	return &createAccountUseCase{
		db:               db,
		searchUC:         searchUC,
		keyManagerClient: keyManagerClient,
		fundAccountUC:    fundAccountUC,
		logger:           log.NewLogger().SetComponent(createAccountComponent),
	}
}

func (uc *createAccountUseCase) Execute(ctx context.Context, account *entities.Account, privateKey, chainName, tenantID string) (*entities.Account, error) {
	ctx = log.WithFields(ctx, log.Field("alias", account.Alias), log.Field("address", account.Address))
	logger := uc.logger.WithContext(ctx)

	logger.Debug("creating new ethereum account")

	accounts, err := uc.searchUC.Execute(ctx, &entities.AccountFilters{Aliases: []string{account.Alias}}, []string{tenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if len(accounts) > 0 {
		errMsg := "alias already exists"
		logger.Error(errMsg)
		return nil, errors.AlreadyExistsError(errMsg).ExtendComponent(createAccountComponent)
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
		logger.WithError(err).Error("failed to import/create ethereum account")
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	account.Address = resp.Address
	account.PublicKey = resp.PublicKey
	account.CompressedPublicKey = resp.CompressedPublicKey
	account.TenantID = tenantID

	// IMPORTANT: Addresses are unique across every tenant
	_, err = uc.db.Account().FindOneByAddress(ctx, account.Address, []string{})
	if err == nil {
		errMsg := "account already exists"
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

	logger.WithField("address", account.Address).Info("ethereum account created successfully")
	return parsers.NewAccountEntityFromModels(accountModel), nil
}
