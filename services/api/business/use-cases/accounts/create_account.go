package accounts

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/ethereum/account"
	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/quorum-key-manager/pkg/client"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const createAccountComponent = "use-cases.create-account"

type createAccountUseCase struct {
	db               store.DB
	searchUC         usecases.SearchAccountsUseCase
	fundAccountUC    usecases.FundAccountUseCase
	keyManagerClient client.EthClient
	logger           *log.Logger
}

func NewCreateAccountUseCase(
	db store.DB,
	searchUC usecases.SearchAccountsUseCase,
	fundAccountUC usecases.FundAccountUseCase,
	keyManagerClient client.EthClient,
) usecases.CreateAccountUseCase {
	return &createAccountUseCase{
		db:               db,
		searchUC:         searchUC,
		keyManagerClient: keyManagerClient,
		fundAccountUC:    fundAccountUC,
		logger:           log.NewLogger().SetComponent(createAccountComponent),
	}
}

func (uc *createAccountUseCase) Execute(ctx context.Context, acc *entities.Account, privateKey hexutil.Bytes, chainName string,
	userInfo *multitenancy.UserInfo) (*entities.Account, error) {
	ctx = log.WithFields(ctx, log.Field("alias", acc.Alias))
	logger := uc.logger.WithContext(ctx)

	logger.Debug("creating new ethereum account")

	accounts, err := uc.searchUC.Execute(ctx,
		&entities.AccountFilters{Aliases: []string{acc.Alias}, TenantID: userInfo.TenantID},
		userInfo)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if len(accounts) > 0 {
		errMsg := "alias already exists"
		logger.Error(errMsg)
		return nil, errors.AlreadyExistsError(errMsg).ExtendComponent(createAccountComponent)
	}

	var accountID = generateKeyID(userInfo.TenantID, acc.Alias)
	var resp *qkmtypes.EthAccountResponse
	if privateKey != nil {
		importedAccount, der := account.NewAccountFromPrivateKey(privateKey.String())
		if der != nil {
			logger.WithError(err).Error("invalid private key")
			return nil, errors.InvalidParameterError(der.Error()).ExtendComponent(createAccountComponent)
		}

		existingAcc, der := uc.db.Account().FindOneByAddress(ctx, importedAccount.Address.Hex(), userInfo.AllowedTenants, userInfo.Username)
		if existingAcc != nil {
			errMsg := "account already exists"
			logger.Error(errMsg)
			return nil, errors.AlreadyExistsError(errMsg).ExtendComponent(createAccountComponent)
		}

		if der != nil && !errors.IsNotFoundError(der) {
			errMsg := "failed to get account"
			logger.WithError(der).Error(errMsg)
			return nil, errors.FromError(der).ExtendComponent(createAccountComponent)
		}

		resp, err = uc.keyManagerClient.ImportEthAccount(ctx, acc.StoreID, &qkmtypes.ImportEthAccountRequest{
			KeyID:      accountID,
			PrivateKey: privateKey,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		})
	} else {
		resp, err = uc.keyManagerClient.CreateEthAccount(ctx, acc.StoreID, &qkmtypes.CreateEthAccountRequest{
			KeyID: accountID,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants:  userInfo.TenantID,
				qkm.TagIDAllowedUsername: userInfo.Username,
			},
		})
	}
	if err != nil {
		errMsg := "failed to import/create ethereum account"
		logger.WithError(err).Error(errMsg)
		return nil, errors.DependencyFailureError(errMsg).ExtendComponent(createAccountComponent)
	}

	acc.Address = resp.Address
	acc.PublicKey = resp.PublicKey
	acc.CompressedPublicKey = resp.CompressedPublicKey
	acc.TenantID = userInfo.TenantID
	acc.OwnerID = userInfo.Username

	accountModel := parsers.NewAccountModelFromEntities(acc)
	err = uc.db.Account().Insert(ctx, accountModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	if chainName != "" {
		err = uc.fundAccountUC.Execute(ctx, acc, chainName, userInfo)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
		}
	}

	logger.WithField("address", acc.Address).Info("ethereum account created successfully")
	return parsers.NewAccountEntityFromModels(accountModel), nil
}

func generateKeyID(tenantID, alias string) string {
	if alias == "" {
		fmt.Println("test!")
		return utils.RandString(20)
	}

	// The goal is to generate an unique ID to prevent duplicated aliases using md5 it generates values compliant
	// with AKV and AWS which requires regex [a-zA-z]+$
	return fmt.Sprintf("%x", md5.Sum([]byte(tenantID+alias)))
}
