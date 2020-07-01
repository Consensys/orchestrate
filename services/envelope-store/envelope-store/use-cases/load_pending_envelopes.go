package usecases

import (
	"context"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

type LoadPendingEnvelopes interface {
	Execute(ctx context.Context, sentBeforeAt time.Time, tenants []string) ([]*models.EnvelopeModel, error)
}

// RegisterContract is a use case to register a new contract
type loadPendingEnvelopes struct {
	envelopeAgent store.EnvelopeAgent
}

// NewGetCatalog creates a new GetCatalog
func NewLoadPendingEnvelopes(envelopeAgent store.EnvelopeAgent) LoadPendingEnvelopes {
	return &loadPendingEnvelopes{
		envelopeAgent: envelopeAgent,
	}
}

func (se *loadPendingEnvelopes) Execute(ctx context.Context, sentBeforeAt time.Time, tenants []string) ([]*models.EnvelopeModel, error) {
	logger := log.FromContext(ctx)

	envelopes, err := se.envelopeAgent.FindPending(ctx, sentBeforeAt, tenants)
	if err != nil {
		logger.
			WithError(err).
			WithFields(logrus.Fields{
				"status":  "pending",
				"sent_at": sentBeforeAt,
			}).
			Debugf("could not load envelope")
	}

	return envelopes, nil
}
