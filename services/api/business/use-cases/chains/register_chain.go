package chains

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const registerChainComponent = "use-cases.register-chain"

// registerChainUseCase is a use case to register a new chain
type registerChainUseCase struct {
	db             store.DB
	searchChainsUC usecases.SearchChainsUseCase
	ethClient      ethclient.Client
}

// NewRegisterChainUseCase creates a new RegisterChainUseCase
func NewRegisterChainUseCase(db store.DB, searchChainsUC usecases.SearchChainsUseCase, ec ethclient.Client) usecases.RegisterChainUseCase {
	return &registerChainUseCase{
		db:             db,
		searchChainsUC: searchChainsUC,
		ethClient:      ec,
	}
}

// Execute registers a new chain
func (uc *registerChainUseCase) Execute(ctx context.Context, chain *entities.Chain, fromLatest bool) (*entities.Chain, error) {
	logger := log.WithContext(ctx).
		WithField("name", chain.Name).
		WithField("tenant", chain.TenantID)
	logger.Debug("registering new chain")

	chains, err := uc.searchChainsUC.Execute(ctx, &entities.ChainFilters{Names: []string{chain.Name}}, []string{chain.TenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(registerChainComponent)
	}

	if len(chains) > 0 {
		errMessage := "a chain with the same name already exists"
		log.WithContext(ctx).Error(errMessage)
		return nil, errors.AlreadyExistsError(errMessage).ExtendComponent(registerChainComponent)
	}

	chainID, err := uc.getChainID(ctx, chain.URLs)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(registerChainComponent)
	}
	chain.ChainID = chainID

	if fromLatest {
		chainTip, der := uc.getChainTip(ctx, chain.URLs)
		if der != nil {
			return nil, errors.FromError(der).ExtendComponent(registerChainComponent)
		}
		chain.ListenerStartingBlock = chainTip
		chain.ListenerCurrentBlock = chainTip
	}

	chainModel := parsers.NewChainModelFromEntity(chain)
	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		der := tx.(store.Tx).Chain().Insert(ctx, chainModel)
		if der != nil {
			return der
		}

		for _, privateTxManager := range chainModel.PrivateTxManagers {
			privateTxManager.ChainUUID = chainModel.UUID

			der = tx.(store.Tx).PrivateTxManager().Insert(ctx, privateTxManager)
			if der != nil {
				return der
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(registerChainComponent)
	}

	logger.
		WithField("chain_uuid", chainModel.UUID).
		WithField("chain_id", chainModel.ChainID).
		Info("chain registered successfully")
	return parsers.NewChainFromModel(chainModel), nil
}

func (uc *registerChainUseCase) getChainID(ctx context.Context, uris []string) (string, error) {
	var prevChainID string
	for i, uri := range uris {
		chainID, err := uc.ethClient.Network(ctx, uri)
		if err != nil {
			errMessage := "failed to fetch chain id"
			log.WithContext(ctx).WithField("url", uri).WithError(err).Error(errMessage)
			return "", errors.InvalidParameterError(errMessage)
		}

		if i > 0 && chainID.String() != prevChainID {
			errMessage := "URLs in the list point to different networks"
			log.WithContext(ctx).
				WithField("url", uri).
				WithField("previous_chain_id", prevChainID).
				WithField("chain_id", chainID.String()).
				Error(errMessage)
			return "", errors.InvalidParameterError(errMessage)
		}

		prevChainID = chainID.String()
	}

	return prevChainID, nil
}

func (uc *registerChainUseCase) getChainTip(ctx context.Context, uris []string) (uint64, error) {
	for _, uri := range uris {
		header, err := uc.ethClient.HeaderByNumber(ctx, uri, nil)
		if err != nil {
			errMessage := "failed to fetch chain tip"
			log.WithContext(ctx).WithField("url", uri).WithError(err).Warning(errMessage)
			continue
		}

		return header.Number.Uint64(), nil
	}

	errMessage := "failed to fetch chain tip for all urls"
	log.WithContext(ctx).WithField("uris", uris).Error(errMessage)
	return 0, errors.InvalidParameterError(errMessage)
}
