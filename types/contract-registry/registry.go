package contract_registry //nolint:golint,stylecheck // reason

//go:generate mockgen -source=registry.pb.go -destination=client/mocks/mock_client.go -package=mocks
