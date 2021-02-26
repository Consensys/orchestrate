package chains

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

const registerChainComponent = "use-cases.register-chain"

// registerChainUseCase is a use case to register a new chain
type registerChainUseCase struct {
	db             store.DB
	searchChainsUC usecases.SearchChainsUseCase
	ethClient      ethclient.Client
	logger         *log.Logger
}

// NewRegisterChainUseCase creates a new RegisterChainUseCase
func NewRegisterChainUseCase(db store.DB, searchChainsUC usecases.SearchChainsUseCase, ec ethclient.Client) usecases.RegisterChainUseCase {
	return &registerChainUseCase{
		db:             db,
		searchChainsUC: searchChainsUC,
		ethClient:      ec,
		logger:         log.NewLogger().SetComponent(registerChainComponent),
	}
}

// Execute registers a new chain
func (uc *registerChainUseCase) Execute(ctx context.Context, chain *entities.Chain, fromLatest bool) (*entities.Chain, error) {
	ctx = log.WithFields(ctx, log.Field("chain_name", chain.Name))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("registering new chain")

	chains, err := uc.searchChainsUC.Execute(ctx, &entities.ChainFilters{Names: []string{chain.Name}}, []string{chain.TenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(registerChainComponent)
	}

	if len(chains) > 0 {
		errMessage := "a chain with the same name already exists"
		logger.Error(errMessage)
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

	logger.WithField("chain_uuid", chainModel.UUID).Info("chain registered successfully")
	return parsers.NewChainFromModel(chainModel), nil
}

func (uc *registerChainUseCase) getChainID(ctx context.Context, uris []string) (string, error) {
	var prevChainID string
	for i, uri := range uris {
		chainID, err := uc.ethClient.Network(ctx, uri)
		if err != nil {
			errMessage := "failed to fetch chain id"
			uc.logger.WithContext(ctx).WithField("url", uri).WithError(err).Error(errMessage)
			return "", errors.InvalidParameterError(errMessage)
		}

		if i > 0 && chainID.String() != prevChainID {
			errMessage := "URLs in the list point to different networks"
			uc.logger.WithContext(ctx).
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
			uc.logger.WithContext(ctx).WithField("url", uri).WithError(err).Warn(errMessage)
			continue
		}

		return header.Number.Uint64(), nil
	}

	errMessage := "failed to fetch chain tip for all urls"
	uc.logger.WithContext(ctx).WithField("uris", uris).Error(errMessage)
	return 0, errors.InvalidParameterError(errMessage)
}
