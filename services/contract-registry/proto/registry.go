package proto

//go:generate mockgen -source=registry.pb.go -destination=../client/mock/grpc.go -package=mock ContractRegistryClient
//go:generate mockgen -source=registry.pb.go -destination=../service/mock/grpc.go -package=mock ContractRegistryServer
