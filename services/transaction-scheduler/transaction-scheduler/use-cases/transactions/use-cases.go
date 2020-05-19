package transactions

type UseCases interface {
	SendTransaction() SendTxUseCase
}
