package controllers

import (
	"context"
	"time"

	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/envelope-store/use-cases"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/models"
)

type GRPCService struct {
	storeEnvelopeUseCase        usecases.StoreEnvelope
	loadEnvelopeByTxHashUseCase usecases.LoadEnvelopeByTxHash
	loadEnvelopeByIDUseCase     usecases.LoadEnvelopeByID
	loadPendingEnvelopesUseCase usecases.LoadPendingEnvelopes
	setEnvelopesStatusUseCase   usecases.SetEnvelopeStatus
}

func NewGRPCService(
	storeda store.DataAgents,
) (*GRPCService, error) {
	return &GRPCService{
		storeEnvelopeUseCase:        usecases.NewStoreEnvelope(storeda.Envelope),
		loadEnvelopeByTxHashUseCase: usecases.NewLoadEnvelopeByTxHash(storeda.Envelope),
		loadEnvelopeByIDUseCase:     usecases.NewLoadEnvelopeByID(storeda.Envelope),
		loadPendingEnvelopesUseCase: usecases.NewLoadPendingEnvelopes(storeda.Envelope),
		setEnvelopesStatusUseCase:   usecases.NewSetEnvelopeStatus(storeda.Envelope),
	}, nil
}

func (s *GRPCService) Store(ctx context.Context, req *svc.StoreRequest) (*svc.StoreResponse, error) {
	envelope, err := s.storeEnvelopeUseCase.Execute(ctx, multitenancy.TenantIDFromContext(ctx), req.GetEnvelope())
	if err != nil {
		return &svc.StoreResponse{}, err
	}

	resp, err := envelopeModelToStoreResponse(&envelope)
	if err != nil {
		return &svc.StoreResponse{}, err
	}

	return resp, nil
}

func (s *GRPCService) LoadByID(ctx context.Context, req *svc.LoadByIDRequest) (*svc.StoreResponse, error) {
	envelope, err := s.loadEnvelopeByIDUseCase.Execute(
		ctx,
		multitenancy.TenantIDFromContext(ctx),
		req.GetId(),
	)
	if err != nil {
		return &svc.StoreResponse{}, err
	}

	resp, err := envelopeModelToStoreResponse(envelope)
	if err != nil {
		return &svc.StoreResponse{}, err
	}

	return resp, nil
}

func (s *GRPCService) LoadByTxHash(ctx context.Context, req *svc.LoadByTxHashRequest) (*svc.StoreResponse, error) {
	envelope, err := s.loadEnvelopeByTxHashUseCase.Execute(
		ctx,
		multitenancy.TenantIDFromContext(ctx),
		req.GetChainId(),
		req.GetTxHash(),
	)
	if err != nil {
		return &svc.StoreResponse{}, err
	}

	resp, err := envelopeModelToStoreResponse(envelope)
	if err != nil {
		return &svc.StoreResponse{}, err
	}

	return resp, nil
}

func (s *GRPCService) SetStatus(ctx context.Context, req *svc.SetStatusRequest) (*svc.StatusResponse, error) {
	envelope, err := s.setEnvelopesStatusUseCase.Execute(
		ctx,
		multitenancy.TenantIDFromContext(ctx),
		req.GetId(),
		req.GetStatus().String(),
	)
	if err != nil {
		return &svc.StatusResponse{}, err
	}

	return &svc.StatusResponse{
		StatusInfo: envelope.StatusInfo(),
	}, nil
}

func (s *GRPCService) LoadPending(ctx context.Context, req *svc.LoadPendingRequest) (*svc.LoadPendingResponse, error) {
	envelopes, err := s.loadPendingEnvelopesUseCase.Execute(
		ctx,
		time.Now().Add(-utils.PDurationToDuration(req.GetDuration())),
	)

	if err != nil {
		return &svc.LoadPendingResponse{}, err
	}

	var resps []*svc.StoreResponse
	for _, envelope := range envelopes {
		resp, err := envelopeModelToStoreResponse(envelope)
		if err != nil {
			return &svc.LoadPendingResponse{}, errors.FromError(err)
		}
		resps = append(resps, resp)
	}

	return &svc.LoadPendingResponse{
		Responses: resps,
	}, nil
}

func envelopeModelToStoreResponse(envelope *models.EnvelopeModel) (*svc.StoreResponse, error) {
	resp := &svc.StoreResponse{
		StatusInfo: envelope.StatusInfo(),
		Envelope:   &tx.TxEnvelope{},
	}

	// Unmarshal envelope
	err := encoding.Unmarshal(envelope.Envelope, resp.GetEnvelope())
	return resp, err
}
