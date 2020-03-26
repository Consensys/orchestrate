# Domain layer
mockgen -source=../contract-registry/use-cases/register_contract.go -destination=../contract-registry/use-cases/mocks/register_contract_mock.go -package=mocks
mockgen -source=../contract-registry/use-cases/get_catalog.go -destination=../contract-registry/use-cases/mocks/get_catalog_mock.go -package=mocks
mockgen -source=../contract-registry/use-cases/get_contract.go -destination=../contract-registry/use-cases/mocks/get_contract_mock.go -package=mocks
mockgen -source=../contract-registry/use-cases/get_events.go -destination=../contract-registry/use-cases/mocks/get_events_mock.go -package=mocks
mockgen -source=../contract-registry/use-cases/get_methods.go -destination=../contract-registry/use-cases/mocks/get_methods_mock.go -package=mocks
mockgen -source=../contract-registry/use-cases/set_codehash.go -destination=../contract-registry/use-cases/mocks/set_codehash_mock.go -package=mocks
mockgen -source=../contract-registry/use-cases/get_tags.go -destination=../contract-registry/use-cases/mocks/get_tags_mock.go -package=mocks

# Data layer
mockgen -source=../store/store.go -destination=../store/mocks/store_mock.go -package=mocks
