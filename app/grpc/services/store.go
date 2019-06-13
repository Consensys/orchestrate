package services

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	types "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store"
)

// StoreService is the service dealing with storing
type StoreService struct {
	store store.EnvelopeStore
}

// NewStoreService creates a StoreService
func NewStoreService(s store.EnvelopeStore) *StoreService {
	return &StoreService{store: s}
}

// Store store a envelope
func (s StoreService) Store(ctx context.Context, req *types.StoreRequest) (*types.StoreResponse, error) {
	status, last, err := s.store.Store(ctx, req.GetEnvelope())
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not store %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not store %v %v", err, req)
	}

	return &types.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
	}, nil
}

// LoadByTxHash load a envelope by transaction hash
func (s StoreService) LoadByTxHash(ctx context.Context, req *types.TxHashRequest) (*types.StoreResponse, error) {
	en := &envelope.Envelope{}
	status, last, err := s.store.LoadByTxHash(ctx, req.GetChainId(), req.GetTxHash(), en)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not load by TxHash %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not load by TxHash  %v %v", err, req)
	}

	return &types.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
		Envelope:    en,
	}, nil
}

// LoadByID load a envelope by identifier
func (s StoreService) LoadByID(ctx context.Context, req *types.IDRequest) (*types.StoreResponse, error) {
	en := &envelope.Envelope{}

	status, last, err := s.store.LoadByID(ctx, req.GetId(), en)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not store load by EnvelopeID %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not store load by EnvelopeID  %v %v", err, req)
	}

	return &types.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
		Envelope:    en,
	}, nil
}

// SetStatus set a envelope status
func (s StoreService) SetStatus(ctx context.Context, req *types.SetStatusRequest) (*common.Error, error) {
	err := s.store.SetStatus(ctx, req.GetId(), req.GetStatus())
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not set status %v %v", err, req)
	}

	return &common.Error{}, nil
}

// GetStatus get a envelope status
func (s StoreService) GetStatus(ctx context.Context, req *types.IDRequest) (*types.StoreResponse, error) {
	status, last, err := s.store.GetStatus(ctx, req.GetId())
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not set status %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not store %v %v", err, req)
	}

	return &types.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
	}, nil
}

// LoadPending load pending envelopes
func (s StoreService) LoadPending(ctx context.Context, req *types.LoadPendingRequest) (*types.LoadPendingResponse, error) {
	envelopes, err := s.store.LoadPending(ctx, time.Duration(req.GetDuration()))
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not load pending envelopess %v %v", err, req)
	}

	return &types.LoadPendingResponse{
		Envelopes: envelopes,
	}, nil
}
