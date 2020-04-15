package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

type LoadEnvelopeByTxHash interface {
	Execute(ctx context.Context, tenantID string, chainID string, txHash string) (*models.EnvelopeModel, error)
}

// RegisterContract is a use case to register a new contract
type loadEnvelopeByTxHash struct {
	envelopeAgent store.EnvelopeAgent
}

// NewGetCatalog creates a new GetCatalog
func NewLoadEnvelopeByTxHash(envelopeAgent store.EnvelopeAgent) LoadEnvelopeByTxHash {
	return &loadEnvelopeByTxHash{
		envelopeAgent: envelopeAgent,
	}
}

func (se *loadEnvelopeByTxHash) Execute(ctx context.Context, tenantID, chainID, txHash string) (*models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelope, err := se.envelopeAgent.FindByFieldSet(ctx, map[string]string{
		"chain_id":  chainID,
		"tx_hash":   txHash,
		"tenant_id": tenantID,
	})

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"chain.id": chainID,
				"tx.hash":  txHash,
				"tenant":   tenantID,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return nil, errors.NotFoundError("no envelope with hash %v", txHash)
		}
		return nil, errors.StorageError(err.Error())
	}

	return envelope, nil
}
