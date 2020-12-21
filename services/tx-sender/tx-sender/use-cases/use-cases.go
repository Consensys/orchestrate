package usecases

//go:generate mockgen -source=use-cases.go -destination=mocks/use-cases.go -package=mocks

type UseCases interface {
	SendETHRawTx() SendETHRawTxUseCase
	SendETHTx() SendETHTxUseCase
	SendEEAPrivateTx() SendEEAPrivateTxUseCase
	SendTesseraPrivateTx() SendTesseraPrivateTxUseCase
	SendTesseraMarkingTx() SendTesseraMarkingTxUseCase
}
