package mocks

import (
	"context"

	"google.golang.org/grpc"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/memory"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

// ContractRegistryClient is a client that wraps an RegistryServer into an ContractRegistryClient
type ContractRegistryClient struct {
	srv svc.ContractRegistryServer
}

func New() *ContractRegistryClient {
	return &ContractRegistryClient{
		srv: memory.NewContractRegistry(),
	}
}

// RegisterContract register a contract including ABI, bytecode and deployed bytecode
func (client *ContractRegistryClient) RegisterContract(ctx context.Context, req *svc.RegisterContractRequest, opts ...grpc.CallOption) (*svc.RegisterContractResponse, error) {
	return client.srv.RegisterContract(ctx, req)
}

// DeregisterContract remove the name + tag association to a contract artifact (abi, bytecode, deployedBytecode). Artifacts are not deleted.
func (client *ContractRegistryClient) DeregisterContract(ctx context.Context, req *svc.DeregisterContractRequest, opts ...grpc.CallOption) (*svc.DeregisterContractResponse, error) {
	return client.srv.DeregisterContract(ctx, req)
}

// DeleteArtifact remove an artifacts based on its BytecodeHash.
func (client *ContractRegistryClient) DeleteArtifact(ctx context.Context, req *svc.DeleteArtifactRequest, opts ...grpc.CallOption) (*svc.DeleteArtifactResponse, error) {
	return client.srv.DeleteArtifact(ctx, req)
}

// GetContractABI loads a contract
func (client *ContractRegistryClient) GetContract(ctx context.Context, req *svc.GetContractRequest, opts ...grpc.CallOption) (*svc.GetContractResponse, error) {
	return client.srv.GetContract(ctx, req)
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

// GetCatalog returns a list of all registered contracts.
func (client *ContractRegistryClient) GetCatalog(ctx context.Context, req *svc.GetCatalogRequest, opts ...grpc.CallOption) (*svc.GetCatalogResponse, error) {
	return client.srv.GetCatalog(ctx, req)
}

// Returns a list of all tags available for a contract name.
func (client *ContractRegistryClient) GetTags(ctx context.Context, req *svc.GetTagsRequest, opts ...grpc.CallOption) (*svc.GetTagsResponse, error) {
	return client.srv.GetTags(ctx, req)
}

// SetAccountCodeHash set the codehash of a contract address for a given chain
func (client *ContractRegistryClient) SetAccountCodeHash(ctx context.Context, req *svc.SetAccountCodeHashRequest, opts ...grpc.CallOption) (*svc.SetAccountCodeHashResponse, error) {
	return client.srv.SetAccountCodeHash(ctx, req)
}
