package identity

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
)

const createIdentityComponent = "use-cases.create-identity"

type createIdentityUseCase struct {
	db                store.DB
	searchUC          usecases.SearchIdentitiesUseCase
	fundingIdentityUC usecases.FundingIdentityUseCase
	keyManagerClient  client.KeyManagerClient
}

func NewCreateIdentityUseCase(db store.DB, searchUC usecases.SearchIdentitiesUseCase, fundingIdentityUC usecases.FundingIdentityUseCase,
	keyManagerClient client.KeyManagerClient) usecases.CreateIdentityUseCase {
	return &createIdentityUseCase{
		db:                db,
		searchUC:          searchUC,
		keyManagerClient:  keyManagerClient,
		fundingIdentityUC: fundingIdentityUC,
	}
}

func (uc *createIdentityUseCase) Execute(ctx context.Context, identity *entities.Identity, privateKey, chainName, tenantID string) (*entities.Identity, error) {
	logger := log.WithContext(ctx).
		WithField("alias", identity.Alias)

	logger.Debug("creating new identity...")
	idens, err := uc.searchUC.Execute(ctx, &entities.IdentityFilters{Aliases: []string{identity.Alias}}, []string{tenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createIdentityComponent)
	}

	if len(idens) > 0 {
		return nil, errors.InvalidParameterError("alias %s already exists", identity.Alias)
	}

	// REMINDER: For now, Identity API only support ETH accounts
	var resp *types.ETHAccountResponse
	if privateKey != "" {
		resp, err = uc.keyManagerClient.ImportETHAccount(ctx, &types.ImportETHAccountRequest{
			Namespace:  tenantID,
			PrivateKey: privateKey,
		})
	} else {
		resp, err = uc.keyManagerClient.CreateETHAccount(ctx, &types.CreateETHAccountRequest{
			Namespace: tenantID,
		})
	}

	if err != nil {
		return nil, err
	}

	identity.Address = resp.Address
	identity.PublicKey = resp.PublicKey
	identity.CompressedPublicKey = resp.CompressedPublicKey
	identity.Active = true

	identityModel := parsers.NewIdentityModelFromEntities(identity)
	identityModel.TenantID = tenantID
	err = uc.db.Identity().Insert(ctx, identityModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(createIdentityComponent)
	}

	if chainName != "" {
		err = uc.fundingIdentityUC.Execute(ctx, identity, chainName)
		if err != nil {
			logger.WithError(err).Error("cannot trigger funding identity")
		}
	}

	logger.WithField("address", identity.Address).Info("identity was created successfully")

	return parsers.NewIdentityEntityFromModels(identityModel), nil
}
