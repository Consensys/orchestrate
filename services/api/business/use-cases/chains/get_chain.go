package chains

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getChainComponent = "use-cases.get-chain"

// getChainUseCase is a use case to get a faucet
type getChainUseCase struct {
	db store.DB
}

// NewGetChainUseCase creates a new GetChainUseCase
func NewGetChainUseCase(db store.DB) usecases.GetChainUseCase {
	return &getChainUseCase{
		db: db,
	}
}

// Execute gets a chain
func (uc *getChainUseCase) Execute(ctx context.Context, uuid string, tenants []string) (*entities.Chain, error) {
	logger := log.WithContext(ctx).WithField("chain_uuid", uuid).WithField("tenants", tenants)
	logger.Debug("getting chain")

	chainModel, err := uc.db.Chain().FindOneByUUID(ctx, uuid, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getChainComponent)
	}

	logger.Debug("chain found successfully")
	return parsers.NewChainFromModel(chainModel), nil
}
