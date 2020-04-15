package proto

//go:generate mockgen -destination=../client/mock/grpc.go -package=mock . ContractRegistryClient
//go:generate mockgen -destination=../service/mock/grpc.go -package=mock . ContractRegistryServer
