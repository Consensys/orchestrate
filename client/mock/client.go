package clientmock

import (
	"context"

	"google.golang.org/grpc"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/registry/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

// ContractRegistryClient is a client that wraps an RegistryServer into an ContractRegistryClient
type ContractRegistryClient struct {
	srv svc.RegistryServer
}

func New(client ethclient.ChainStateReader) *ContractRegistryClient {
	return &ContractRegistryClient{
		srv: mock.NewRegistry(client),
	}
}

// RegisterContract register a contract including ABI, bytecode and deployed bytecode
func (client *ContractRegistryClient) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest, opts ...grpc.CallOption) (*svc.RegisterContractResponse, error) {
	return client.srv.RegisterContract(ctx, req)
}

// GetContractABI loads contract ABI
func (client *ContractRegistryClient) GetContractABI(ctx context.Context, req *svc.GetContractRequest, opts ...grpc.CallOption) (*svc.GetContractABIResponse, error) {
	return client.srv.GetContractABI(ctx, req)
}

// GetContractBytecode loads contract bytecode
func (client *ContractRegistryClient) GetContractBytecode(ctx context.Context, req *svc.GetContractRequest, opts ...grpc.CallOption) (*svc.GetContractBytecodeResponse, error) {
	return client.srv.GetContractBytecode(ctx, req)
}

// GetContractDeployedBytecode loads contract deployed bytecode
func (client *ContractRegistryClient) GetContractDeployedBytecode(ctx context.Context, req *svc.GetContractRequest, opts ...grpc.CallOption) (*svc.GetContractDeployedBytecodeResponse, error) {
	return client.srv.GetContractDeployedBytecode(ctx, req)
}

// GetMethodsBySelector load method using 4 bytes unique selector and the address of the contract
func (client *ContractRegistryClient) GetMethodsBySelector(ctx context.Context, req *svc.GetMethodsBySelectorRequest, opts ...grpc.CallOption) (*svc.GetMethodsBySelectorResponse, error) {
	return client.srv.GetMethodsBySelector(ctx, req)
}

// GetEventsBySigHash load event using event signature hash
func (client *ContractRegistryClient) GetEventsBySigHash(ctx context.Context, req *svc.GetEventsBySigHashRequest, opts ...grpc.CallOption) (*svc.GetEventsBySigHashResponse, error) {
	return client.srv.GetEventsBySigHash(ctx, req)
}

// RequestAddressUpdate Request an update of the codehash of the contract address
func (client *ContractRegistryClient) RequestAddressUpdate(ctx context.Context, req *svc.AddressUpdateRequest, opts ...grpc.CallOption) (*svc.AddressUpdateResponse, error) {
	return client.srv.RequestAddressUpdate(ctx, req)
}
