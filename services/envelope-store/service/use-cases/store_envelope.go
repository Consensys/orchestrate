package usecases

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

//go:generate mockgen -source=store_envelope.go -destination=mocks/store_envelope.go -package=mocks

type StoreEnvelope interface {
	Execute(ctx context.Context, tenantID string, envelopeItem *tx.TxEnvelope) (models.EnvelopeModel, error)
}

// RegisterContract is a use case to register a new contract
type storeEnvelope struct {
	envelopeAgent store.EnvelopeAgent
}

// NewGetCatalog creates a new GetCatalog
func NewStoreEnvelope(envelopeAgent store.EnvelopeAgent) StoreEnvelope {
	return &storeEnvelope{
		envelopeAgent: envelopeAgent,
	}
}

func (se *storeEnvelope) Execute(ctx context.Context, tenantID string, envelopeTx *tx.TxEnvelope) (models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	// create envelopeItem from envelope
	envelope, err := models.NewEnvelopeFromTx(tenantID, envelopeTx)
	if err != nil {
		logger.WithError(err).Errorf("invalid envelope")
		return models.EnvelopeModel{}, err
	}

	err = se.envelopeAgent.InsertDoUpdateOnEnvelopeIDKey(ctx, &envelope)
	if err != nil {
		err = se.envelopeAgent.InsertDoUpdateOnUniTx(ctx, &envelope)
		if err != nil {
			logger.WithError(err).Errorf("could not store envelope")
			return models.EnvelopeModel{}, err
		}
	}

	log.FromContext(ctx).
		WithFields(logrus.Fields{
			"chain.id":    envelope.ChainID,
			"tx.hash":     envelope.TxHash,
			"tenant":      envelope.TenantID,
			"envelope.id": envelope.EnvelopeID,
		}).
		Infof("envelope stored")

	return envelope, nil
}
