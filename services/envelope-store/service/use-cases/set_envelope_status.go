package usecases

import (
	"context"
	"strings"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

//go:generate mockgen -source=set_envelope_status.go -destination=mocks/set_envelope_status.go -package=mocks

type SetEnvelopeStatus interface {
	Execute(ctx context.Context, tenantID string, envelopeID string, nextStatus string) (models.EnvelopeModel, error)
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

func (se *setEnvelopeStatus) Execute(ctx context.Context, tenantID, envelopeID, nextStatus string) (models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelope, err := se.envelopeAgent.FindByFieldSet(ctx, map[string]string{
		"envelope_id": envelopeID,
		"tenant_id":   tenantID,
	})

	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"envelope.id": envelopeID,
				"tenant":      tenantID,
			}).
			Debugf("could not load envelope")
		return models.EnvelopeModel{}, errors.StorageError(err.Error())
	}

	envelope.Status = strings.ToLower(nextStatus)

	err = se.envelopeAgent.UpdateStatus(ctx, &envelope)
	if err != nil {
		logger.WithError(err).Errorf("could not update envelope")
		return models.EnvelopeModel{}, err
	}

	log.FromContext(ctx).
		WithFields(logrus.Fields{
			"status":      envelope.Status,
			"tenant":      envelope.TenantID,
			"envelope.id": envelope.EnvelopeID,
		}).
		Infof("envelope updated")

	return envelope, nil
}
