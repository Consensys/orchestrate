package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

type Builder struct {
	postgres postgres.Manager
}

func NewBuilder(mngr postgres.Manager) *Builder {
	return &Builder{postgres: mngr}
}

func (b *Builder) Build(ctx context.Context, conf *Config) (*PG, error) {
	return New(b.postgres.Connect(ctx, conf.PG)), nil
}

// PG is an Envelope Store store based on PostgreSQL
type PG struct {
	db *pg.DB
}

// NewEnvelopeStore creates a new envelope store
func New(db *pg.DB) *PG {
	return &PG{db: db}
}

// Store context envelope
func (s *PG) Store(ctx context.Context, req *svc.StoreRequest) (*svc.StoreResponse, error) {
	logger := log.FromContext(ctx)

	// create model from envelope
	model, err := FromEnvelope(ctx, req.GetEnvelope())
	if err != nil {
		logger.WithError(err).Errorf("invalid envelope")
		return &svc.StoreResponse{}, errors.FromError(err)
	}

	// Execute ORM query
	// If uniqueness constraint is broken then it update the former value
	_, err = s.db.ModelContext(ctx, model).
		OnConflict("ON CONSTRAINT envelopes_envelope_id_key DO UPDATE").
		Set("envelope = ?envelope").
		Set("chain_id = ?chain_id").
		Set("tx_hash = ?tx_hash").
		Returning("*").
		Insert()
	if err != nil {
		// Possibly we got an error due to unique contraint on tx,chain_id so we try again
		_, err = s.db.ModelContext(ctx, model).
			OnConflict("ON CONSTRAINT uni_tx DO UPDATE").
			Set("envelope = ?envelope").
			Set("envelope_id = ?envelope_id").
			Returning("*").
			Insert()
		if err != nil {
			logger.WithError(err).Errorf("could not store envelope")
			return &svc.StoreResponse{}, errors.StorageError("%v", err)
		}
	}

	log.FromContext(ctx).
		WithFields(logrus.Fields{
			"chain.id":    model.ChainID,
			"tx.hash":     model.TxHash,
			"tenant":      model.TenantID,
			"envelope.id": model.EnvelopeID,
		}).
		Infof("envelope stored")

	return model.ToStoreResponse()
}

// LoadByTxHash load envelope by transaction hash
func (s *PG) LoadByTxHash(ctx context.Context, req *svc.LoadByTxHashRequest) (*svc.StoreResponse, error) { //nolint:interfacer // reason
	tenantID := multitenancy.TenantIDFromContext(ctx)
	model := &EnvelopeModel{
		ChainID:  req.GetChainId(),
		TenantID: tenantID,
		TxHash:   req.GetTxHash(),
	}

	err := s.db.ModelContext(ctx, model).
		Where("chain_id = ?", model.ChainID).
		Where("tx_hash = ?", model.TxHash).
		Where("tenant_id = ?", model.TenantID).
		Select()
	if err != nil {
		log.FromContext(ctx).
			WithError(err).
			WithFields(logrus.Fields{
				"chain.id": model.ChainID,
				"tx.hash":  model.TxHash,
				"tenant":   model.TenantID,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return &svc.StoreResponse{}, errors.NotFoundError("no envelope with hash %v", model.TxHash)
		}
		return &svc.StoreResponse{}, errors.StorageError(err.Error())
	}

	return model.ToStoreResponse()
}

// LoadByID context envelope by envelope UUID
func (s *PG) LoadByID(ctx context.Context, req *svc.LoadByIDRequest) (*svc.StoreResponse, error) { //nolint:interfacer // reason
	tenantID := multitenancy.TenantIDFromContext(ctx)
	model := &EnvelopeModel{
		EnvelopeID: req.GetId(),
		TenantID:   tenantID,
	}

	err := s.db.ModelContext(ctx, model).
		Where("envelope_id = ?", model.EnvelopeID).
		Where("tenant_id = ?", model.TenantID).
		Select()
	if err != nil {
		log.FromContext(ctx).WithError(err).
			WithFields(logrus.Fields{
				"envelope.id": model.EnvelopeID,
				"tenant":      model.TenantID,
			}).
			Debugf("could not load envelope")
		if err == pg.ErrNoRows {
			return &svc.StoreResponse{}, errors.NotFoundError("envelope with id %v does not exist", model.EnvelopeID)
		}
		return &svc.StoreResponse{}, errors.StorageError(err.Error())
	}

	return model.ToStoreResponse()
}

// SetStatus set a context status
func (s *PG) SetStatus(ctx context.Context, req *svc.SetStatusRequest) (*svc.StatusResponse, error) {
	tenantID := multitenancy.TenantIDFromContext(ctx)
	// Define model
	model := &EnvelopeModel{
		EnvelopeID: req.GetId(),
		TenantID:   tenantID,
		Status:     strings.ToLower(req.GetStatus().String()),
	}

	// Update status value
	_, err := s.db.ModelContext(ctx, model).
		Set("status = ?status").
		Where("envelope_id = ?", model.EnvelopeID).
		Where("tenant_id = ?", model.TenantID).
		Returning("*").
		Update()
	if err != nil {
		log.FromContext(ctx).WithError(err).
			WithFields(logrus.Fields{
				"envelope.id": model.ChainID,
				"tenant":      model.TenantID,
				"status":      model.Status,
			}).
			Errorf("could not set envelope status")
		return &svc.StatusResponse{}, errors.NotFoundError("envelope not found")
	}

	return model.ToStatusResponse()
}

// LoadPending loads pending envelopes
func (s *PG) LoadPending(ctx context.Context, req *svc.LoadPendingRequest) (*svc.LoadPendingResponse, error) {
	var models []*EnvelopeModel

	err := s.db.ModelContext(ctx, &models).
		Where("status = 'pending'").
		Where("sent_at < ?", time.Now().Add(-utils.PDurationToDuration(req.GetDuration()))).
		Select()
	if err != nil {
		return nil, errors.NotFoundError("envelope not found")
	}

	var resps []*svc.StoreResponse
	for _, model := range models {
		resp, err := model.ToStoreResponse()
		if err != nil {
			return &svc.LoadPendingResponse{}, errors.FromError(err)
		}
		resps = append(resps, resp)
	}

	return &svc.LoadPendingResponse{
		Responses: resps,
	}, nil
}
