package use_cases

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
	Execute(ctx context.Context, tenantId string, chainId string, txHash string) (models.EnvelopeModel, error)
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

func (se *loadEnvelopeByTxHash) Execute(ctx context.Context, tenantId string, chainId string, txHash string) (models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelope, err := se.envelopeAgent.FindByFieldSet(ctx, map[string]string{
		"chain_id": chainId,
		"tx_hash": txHash,
		"tenant_id": tenantId,
	})

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"chain.id": chainId,
				"tx.hash":  txHash,
				"tenant":   tenantId,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return models.EnvelopeModel{}, errors.NotFoundError("no envelope with hash %v", txHash)
		}
		return models.EnvelopeModel{}, errors.StorageError(err.Error())
	}
	
	return envelope, nil
}
