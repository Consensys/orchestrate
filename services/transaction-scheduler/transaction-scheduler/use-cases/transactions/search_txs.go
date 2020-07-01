package transactions

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

//go:generate mockgen -source=search_txs.go -destination=mocks/search_txs.go -package=mocks

const searchTxsComponent = "use-cases.search-txs"

type SearchTransactionsUseCase interface {
	Execute(ctx context.Context, filters *entities.TransactionFilters, tenants []string) ([]*entities.TxRequest, error)
}

// searchTransactionsUseCase is a use case to get transaction requests by filter (or all)
type searchTransactionsUseCase struct {
	db           store.DB
	getTxUseCase GetTxUseCase
}

// NewSearchTransactionsUseCase creates a new SearchTransactionsUseCase
func NewSearchTransactionsUseCase(db store.DB, getTxUseCase GetTxUseCase) SearchTransactionsUseCase {
	return &searchTransactionsUseCase{
		db:           db,
		getTxUseCase: getTxUseCase,
	}
}

// Execute gets a transaction requests by filter (or all)
func (uc *searchTransactionsUseCase) Execute(ctx context.Context, filters *entities.TransactionFilters, tenants []string) ([]*entities.TxRequest, error) {
	log.WithContext(ctx).WithField("filters", filters).Debug("search transaction requests")

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(searchTxsComponent)
	}

	txRequestModels, err := uc.db.TransactionRequest().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchTxsComponent)
	}

	var txRequests []*entities.TxRequest
	for _, txRequestModel := range txRequestModels {
		txRequest, err := uc.getTxUseCase.Execute(ctx, txRequestModel.UUID, tenants)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(searchTxsComponent)
		}

		txRequests = append(txRequests, txRequest)
	}

	log.WithContext(ctx).WithField("filters", filters).Info("transaction requests found successfully")

	return txRequests, nil
}
