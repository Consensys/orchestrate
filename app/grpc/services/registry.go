package services

import (
	"context"
	"encoding/json"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	types "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry"
)

// RegistryService is the service dealing with registering contracts
type RegistryService struct {
	registry registry.Registry
}

// NewRegistryService creates a RegistryService
func NewRegistryService(r registry.Registry) *RegistryService {
	return &RegistryService{registry: r}
}

// RegisterContract register a contract (abi, bytecode & deployedBytecode) on the registry
func (r RegistryService) RegisterContract(ctx context.Context, req *types.RegisterContractRequest) (*types.RegisterContractResponse, error) {
	err := r.registry.RegisterContract(req.Contract)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not register contract %v %v", err, req)
	}
	return &types.RegisterContractResponse{}, nil
}

// LoadByTxHash load a envelope by transaction hash
func (r RegistryService) GetContractABI(ctx context.Context, req *types.GetContractRequest) (*types.GetContractABIResponse, error) {
	abi, err := r.registry.GetContractABI(req.Contract)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not get contract abi %v %v", err, req)
	}
	return &types.GetContractABIResponse{
		Abi: abi,
	}, nil
}

// Returns the bytecode
func (r RegistryService) GetContractBytecode(ctx context.Context, req *types.GetContractRequest) (*types.GetContractBytecodeResponse, error) {
	bytecode, err := r.registry.GetContractBytecode(req.Contract)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not get contract bytecode %v %v", err, req)
	}
	return &types.GetContractBytecodeResponse{
		Bytecode: bytecode,
	}, nil
}

// Returns the deployed bytecode
func (r RegistryService) GetContractDeployedBytecode(ctx context.Context, req *types.GetContractRequest) (*types.GetContractDeployedBytecodeResponse, error) {
	deployedBytecode, err := r.registry.GetContractDeployedBytecode(req.Contract)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not get contract deployedBytecode %v %v", err, req)
	}
	return &types.GetContractDeployedBytecodeResponse{
		DeployedBytecode: deployedBytecode,
	}, nil
}

// Retrieve methods using 4 bytes unique selector
func (r RegistryService) GetMethodsBySelector(ctx context.Context, req *types.GetMethodsBySelectorRequest) (*types.GetMethodsBySelectorResponse, error) {
	var sel [4]byte
	copy(sel[:], req.Selector)
	method, defaultMethods, err := r.registry.GetMethodsBySelector(sel, *req.AccountInstance)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not get methods %v %v", err, req)
	}

	pMethod, err := json.Marshal(method)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not unmarshal method %v %v", err, req)
	}
	pDefaultMethods := make([][]byte, len(defaultMethods))
	for i := range defaultMethods {
		pDefaultMethods[i], err = json.Marshal(defaultMethods[i])
		if err != nil {
			return nil, grpcStatus.Errorf(codes.Internal, "Could not unmarshal defaultMethods %v %v", err, req)
		}
	}

	return &types.GetMethodsBySelectorResponse{
		Method:         pMethod,
		DefaultMethods: pDefaultMethods,
	}, nil
}

// Retrieve events using 4 bytes unique selector
func (r RegistryService) GetEventsBySigHash(ctx context.Context, req *types.GetEventsBySigHashRequest) (*types.GetEventsBySigHashResponse, error) {
	sigHash := ethCommon.BytesToHash(req.SigHash)
	event, defaultEvents, err := r.registry.GetEventsBySigHash(sigHash, *req.AccountInstance, uint(req.IndexedInputCount))
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not get events %v %v", err, req)
	}

	pEvent, err := json.Marshal(event)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not unmarshal event %v %v", err, req)
	}
	pDefaultEvents := make([][]byte, len(defaultEvents))
	for i := range defaultEvents {
		pDefaultEvents[i], err = json.Marshal(defaultEvents[i])
		if err != nil {
			return nil, grpcStatus.Errorf(codes.Internal, "Could not unmarshal defaultEvents %v %v", err, req)
		}
	}

	return &types.GetEventsBySigHashResponse{
		Event:         pEvent,
		DefaultEvents: pDefaultEvents,
	}, nil
}

// Request an update of the codehash of the contract address
func (r RegistryService) RequestAddressUpdate(ctx context.Context, req *types.AddressUpdateRequest) (*types.AddressUpdateResponse, error) {
	err := r.registry.RequestAddressUpdate(*req.AccountInstance)
	if err != nil {
		return nil, grpcStatus.Errorf(codes.Internal, "Could not register contract %v %v", err, req)
	}
	return &types.AddressUpdateResponse{}, nil
}
