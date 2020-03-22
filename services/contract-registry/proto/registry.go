package proto

//go:generate mockgen -destination=../client/mock/mock.go -package=mock . ContractRegistryClient
//go:generate mockgen -destination=../service/mock/mock.go -package=mock . ContractRegistryServer
