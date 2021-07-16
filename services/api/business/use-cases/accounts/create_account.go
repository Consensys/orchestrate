package accounts

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/consensys/quorum-key-manager/pkg/client"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	createAccountComponent = "use-cases.create-account"
)

type createAccountUseCase struct {
	db               store.DB
	searchUC         usecases.SearchAccountsUseCase
	fundAccountUC    usecases.FundAccountUseCase
	keyManagerClient client.Eth1Client
	storeName        string
	logger           *log.Logger
}

func NewCreateAccountUseCase(db store.DB, searchUC usecases.SearchAccountsUseCase, fundAccountUC usecases.FundAccountUseCase,
	keyManagerClient client.Eth1Client) usecases.CreateAccountUseCase {
	return &createAccountUseCase{
		db:               db,
		searchUC:         searchUC,
		keyManagerClient: keyManagerClient,
		fundAccountUC:    fundAccountUC,
		logger:           log.NewLogger().SetComponent(createAccountComponent),
		storeName:        qkm.GlobalStoreName(),
	}
}

func (uc *createAccountUseCase) Execute(ctx context.Context, account *entities.Account, privateKey hexutil.Bytes, chainName, tenantID string) (*entities.Account, error) {
	ctx = log.WithFields(ctx, log.Field("alias", account.Alias))
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

	var accountID = generateKeyID(tenantID, account.Alias)
	var resp *qkmtypes.Eth1AccountResponse
	if privateKey != nil {
		resp, err = uc.keyManagerClient.ImportEth1Account(ctx, uc.storeName, &qkmtypes.ImportEth1AccountRequest{
			KeyID:      accountID,
			PrivateKey: privateKey,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		})

		// In case key already exists we need to append the allowed tenants
		if err != nil && isAccountAlreadyExistErr(err) {
			logger.WithError(err).Debug("duplicated account has been imported")
			privKey, _ := crypto.HexToECDSA(privateKey.String()[2:])
			address := crypto.PubkeyToAddress(privKey.PublicKey).Hex()
			resp, err = uc.keyManagerClient.GetEth1Account(ctx, uc.storeName, address)
			if err == nil {
				logger.WithField("address", address).Debug("updating account to amend allowed tenants")
				// @TODO Prevent duplicated tenantIds
				curTags := resp.Tags
				curTags[qkm.TagIDAllowedTenants] += qkm.TagSeparatorAllowedTenants + tenantID
				_, err = uc.keyManagerClient.UpdateEth1Account(ctx, uc.storeName, address, &qkmtypes.UpdateEth1AccountRequest{
					Tags: curTags,
				})
			} else {
				logger.WithError(err).WithField("address", address).Debug("failed to find account")
			}
		}
	} else {
		resp, err = uc.keyManagerClient.CreateEth1Account(ctx, uc.storeName, &qkmtypes.CreateEth1AccountRequest{
			KeyID: accountID,
			Tags: map[string]string{
				qkm.TagIDAllowedTenants: tenantID,
			},
		})
	}

	if err != nil {
		errMsg := "failed to import/create ethereum account"
		logger.WithError(err).Error(errMsg)
		return nil, errors.FromError(err).ExtendComponent(createAccountComponent)
	}

	account.Address = resp.Address.String()
	account.PublicKey = resp.PublicKey.String()
	account.CompressedPublicKey = resp.CompressedPublicKey.String()
	account.TenantID = tenantID

	// TODO Discuss decision made on allowing same account imported over different tenants
	_, err = uc.db.Account().FindOneByAddress(ctx, account.Address, []string{tenantID})
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

func isAccountAlreadyExistErr(err interface{}) bool {
	if err == nil {
		return false
	}
	if qerr, ok := err.(*client.ResponseError); ok {
		return qerr.ErrorCode == qkm.AlreadyExists
	}
	return false
}

func generateKeyID(tenantID, alias string) string {
	if alias == "" {
		return utils.RandString(14)
	}

	return fmt.Sprintf("%s%s", tenantID, alias)
}
