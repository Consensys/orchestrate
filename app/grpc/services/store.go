package services

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/context-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// StoreService is the service dealing with storing
type StoreService struct {
	store infra.TraceStore
}

// NewStoreService creates a StoreService
func NewStoreService(store infra.TraceStore) *StoreService {
	return &StoreService{store: store}
}

// Store store a trace
func (s StoreService) Store(ctx context.Context, req *store.StoreRequest) (*store.StoreResponse, error) {
	status, last, err := s.store.Store(ctx, req.GetTrace())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not store %v %v", err)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not store %v %v", err, req)
	}

	return &store.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
	}, nil
}

// LoadByTxHash load a trace by transaction hash
func (s StoreService) LoadByTxHash(ctx context.Context, req *store.TxHashRequest) (*store.StoreResponse, error) {
	tr := &trace.Trace{}
	status, last, err := s.store.LoadByTxHash(ctx, req.GetChainId(), req.GetTxHash(), tr)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not load by TxHash %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not load by TxHash  %v %v", err, req)
	}

	return &store.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
		Trace:       tr,
	}, nil
}

// LoadByTraceID load a trace by identifier
func (s StoreService) LoadByTraceID(ctx context.Context, req *store.TraceIDRequest) (*store.StoreResponse, error) {
	tr := &trace.Trace{}
	status, last, err := s.store.LoadByTraceID(ctx, req.GetTraceId(), tr)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not store load by TraceID %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not store load by TraceID  %v %v", err, req)
	}

	return &store.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
		Trace:       tr,
	}, nil
}

// SetStatus set a trace status
func (s StoreService) SetStatus(ctx context.Context, req *store.SetStatusRequest) (*common.Error, error) {
	err := s.store.SetStatus(ctx, req.GetTraceId(), req.GetStatus())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not set status %v %v", err, req)
	}

	return &common.Error{}, nil
}

// GetStatus get a trace status
func (s StoreService) GetStatus(ctx context.Context, req *store.TraceIDRequest) (*store.StoreResponse, error) {
	status, last, err := s.store.GetStatus(ctx, req.GetTraceId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not set status %v %v", err, req)
	}

	lastUpdated, err := ptypes.TimestampProto(last)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not store %v %v", err, req)
	}

	return &store.StoreResponse{
		Status:      status,
		LastUpdated: lastUpdated,
	}, nil
}

// LoadPendingTraces load pending traces
func (s StoreService) LoadPendingTraces(ctx context.Context, req *store.PendingTracesRequest) (*store.PendingTracesResponse, error) {
	traces, err := s.store.LoadPendingTraces(ctx, time.Duration(req.GetDuration()))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Could not load pending tracess %v %v", err, req)
	}

	return &store.PendingTracesResponse{
		Traces: traces,
	}, nil
}
