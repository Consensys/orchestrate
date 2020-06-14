package chains

type UseCases interface {
	GetChainByName() GetChainByNameUseCase
}
