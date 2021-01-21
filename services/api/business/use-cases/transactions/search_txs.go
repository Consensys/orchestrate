package transactions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const searchTxsComponent = "use-cases.search-txs"

// searchTransactionsUseCase is a use case to get transaction requests by filter (or all)
type searchTransactionsUseCase struct {
	db           store.DB
	getTxUseCase usecases.GetTxUseCase
	logger       *log.Logger
}

// NewSearchTransactionsUseCase creates a new SearchTransactionsUseCase
func NewSearchTransactionsUseCase(db store.DB, getTxUseCase usecases.GetTxUseCase) usecases.SearchTransactionsUseCase {
	return &searchTransactionsUseCase{
		db:           db,
		getTxUseCase: getTxUseCase,
		logger:       log.NewLogger().SetComponent(searchTxsComponent),
	}
}

// Execute gets a transaction requests by filter (or all)
func (uc *searchTransactionsUseCase) Execute(ctx context.Context, filters *entities.TransactionRequestFilters, tenants []string) ([]*entities.TxRequest, error) {
	txRequestModels, err := uc.db.TransactionRequest().Search(ctx, filters, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(searchTxsComponent)
	}

	var txRequests []*entities.TxRequest
	for _, txRequestModel := range txRequestModels {
		txRequest, err := uc.getTxUseCase.Execute(ctx, txRequestModel.Schedule.UUID, tenants)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(searchTxsComponent)
		}

		txRequests = append(txRequests, txRequest)
	}

	uc.logger.Info("transaction requests found successfully")

	return txRequests, nil
}
