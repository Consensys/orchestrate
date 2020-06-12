package mocks

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"google.golang.org/grpc"
)
type ContractRegistryFaker struct {
	GetContract func() (*proto.GetContractResponse, error)
}
type contractRegistryClientMock struct {
	Faker *ContractRegistryFaker
}

func NewContractRegistryClientMock(faker *ContractRegistryFaker) proto.ContractRegistryClient {
	return &contractRegistryClientMock{
		Faker: faker,
	}
}

func (c contractRegistryClientMock) RegisterContract(ctx context.Context, in *proto.RegisterContractRequest, opts ...grpc.CallOption) (*proto.RegisterContractResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) DeregisterContract(ctx context.Context, in *proto.DeregisterContractRequest, opts ...grpc.CallOption) (*proto.DeregisterContractResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) DeleteArtifact(ctx context.Context, in *proto.DeleteArtifactRequest, opts ...grpc.CallOption) (*proto.DeleteArtifactResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetContract(ctx context.Context, in *proto.GetContractRequest, opts ...grpc.CallOption) (*proto.GetContractResponse, error) {
	return c.Faker.GetContract()
	// return &proto.GetContractResponse{
	// 	Contract: testutils.FakeContract(),
	// }, nil
}

func (c contractRegistryClientMock) GetContractABI(ctx context.Context, in *proto.GetContractRequest, opts ...grpc.CallOption) (*proto.GetContractABIResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetContractBytecode(ctx context.Context, in *proto.GetContractRequest, opts ...grpc.CallOption) (*proto.GetContractBytecodeResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetContractDeployedBytecode(ctx context.Context, in *proto.GetContractRequest, opts ...grpc.CallOption) (*proto.GetContractDeployedBytecodeResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetCatalog(ctx context.Context, in *proto.GetCatalogRequest, opts ...grpc.CallOption) (*proto.GetCatalogResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetTags(ctx context.Context, in *proto.GetTagsRequest, opts ...grpc.CallOption) (*proto.GetTagsResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetMethodSignatures(ctx context.Context, in *proto.GetMethodSignaturesRequest, opts ...grpc.CallOption) (*proto.GetMethodSignaturesResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetMethodsBySelector(ctx context.Context, in *proto.GetMethodsBySelectorRequest, opts ...grpc.CallOption) (*proto.GetMethodsBySelectorResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) GetEventsBySigHash(ctx context.Context, in *proto.GetEventsBySigHashRequest, opts ...grpc.CallOption) (*proto.GetEventsBySigHashResponse, error) {
	panic("implement me")
}

func (c contractRegistryClientMock) SetAccountCodeHash(ctx context.Context, in *proto.SetAccountCodeHashRequest, opts ...grpc.CallOption) (*proto.SetAccountCodeHashResponse, error) {
	panic("implement me")
}
