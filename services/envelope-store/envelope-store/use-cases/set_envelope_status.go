package usecases

import (
	"context"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

type SetEnvelopeStatus interface {
	Execute(ctx context.Context, tenants []string, envelopeID string, nextStatus string) (*models.EnvelopeModel, error)
}

// RegisterContract is a use case to register a new contract
type setEnvelopeStatus struct {
	envelopeAgent store.EnvelopeAgent
}

// NewGetCatalog creates a new GetCatalog
func NewSetEnvelopeStatus(envelopeAgent store.EnvelopeAgent) SetEnvelopeStatus {
	return &setEnvelopeStatus{
		envelopeAgent: envelopeAgent,
	}
}

func (se *setEnvelopeStatus) Execute(ctx context.Context, tenants []string, envelopeID, nextStatus string) (*models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelope, err := se.envelopeAgent.FindByFieldSet(ctx, map[string]string{
		"envelope_id": envelopeID,
	}, tenants)

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"id":      envelopeID,
				"tenants": tenants,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return nil, errors.NotFoundError("no envelope with envelope id %v", envelopeID)
		}
		return nil, errors.StorageError(err.Error())
	}

	envelope.Status = strings.ToLower(nextStatus)

	err = se.envelopeAgent.UpdateStatus(ctx, envelope, tenants)
	if err != nil {
		logger.WithError(err).Errorf("could not update envelope")
		return nil, err
	}

	log.FromContext(ctx).
		WithFields(logrus.Fields{
			"status":   envelope.Status,
			"tenantID": envelope.TenantID,
			"id":       envelope.EnvelopeID,
		}).
		Infof("envelope updated")

	return envelope, nil
}
