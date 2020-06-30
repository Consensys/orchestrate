package transactions

type UseCases interface {
	SendContractTransaction() SendContractTxUseCase
	SendDeployTransaction() SendDeployTxUseCase
	SendTransaction() SendTxUseCase
	GetTransaction() GetTxUseCase
	SearchTransactions() SearchTransactionsUseCase
}
