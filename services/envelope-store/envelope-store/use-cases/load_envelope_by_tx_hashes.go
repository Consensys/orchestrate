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

type LoadEnvelopeByTxHashes interface {
	Execute(ctx context.Context, tenantID, chainID string, txHashes []string) ([]*models.EnvelopeModel, error)
}

type loadEnvelopeByTxHashes struct {
	envelopeAgent store.EnvelopeAgent
}

func NewLoadEnvelopeByTxHashes(envelopeAgent store.EnvelopeAgent) LoadEnvelopeByTxHashes {
	return &loadEnvelopeByTxHashes{
		envelopeAgent: envelopeAgent,
	}
}

func (se *loadEnvelopeByTxHashes) Execute(ctx context.Context, tenantID, chainID string, txHashes []string) ([]*models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	// TODO: Filter also by tenantID
	envelopes, err := se.envelopeAgent.FindByTxHashes(ctx, txHashes)

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"chain.id":  chainID,
				"tx.hashes": txHashes,
				"tenant":    tenantID,
			}).
			Debugf("could not load envelopes")
		if err == pg.ErrNoRows {
			return nil, errors.NotFoundError("no envelope with hashes %v", txHashes)
		}
		return nil, errors.StorageError(err.Error())
	}

	return envelopes, nil
}
