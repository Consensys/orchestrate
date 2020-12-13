package builder

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	nonce2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/nonce"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases/sender"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-signer/tx-signer/use-cases/signer"
)

type useCases struct {
	sendETHTx            usecases.SendETHTxUseCase
	sendETHRawTx         usecases.SendETHRawTxUseCase
	sendEEAPrivateTx     usecases.SendEEAPrivateTxUseCase
	sendTesseraPrivateTx usecases.SendTesseraPrivateTxUseCase
	sendTesseraMarkingTx usecases.SendTesseraMarkingTxUseCase
}

func NewUseCases(txSchedulerClient client2.TransactionSchedulerClient, keyManagerClient client.KeyManagerClient,
	ec ethclient.MultiClient, nonceManager nonce2.Sender, chainRegistryURL string, checkerMaxRecovery uint64) usecases.UseCases {
	signETHTransactionUC := signer.NewSignETHTransactionUseCase(keyManagerClient)
	signEEATransactionUC := signer.NewSignEEATransactionUseCase(keyManagerClient)
	signQuorumTransactionUC := signer.NewSignQuorumPrivateTransactionUseCase(keyManagerClient)

	checker := nonce.NewNonceChecker(ec, nonceManager, nonce.NewRecoveryTracker(), chainRegistryURL, checkerMaxRecovery)
	return &useCases{
		sendETHTx:            sender.NewSendEthTxUseCase(signETHTransactionUC, ec, txSchedulerClient, chainRegistryURL, checker),
		sendETHRawTx:         sender.NewSendETHRawTxUseCase(ec, txSchedulerClient, chainRegistryURL),
		sendEEAPrivateTx:     sender.NewSendEEAPrivateTxUseCase(signEEATransactionUC, ec, txSchedulerClient, chainRegistryURL, checker),
		sendTesseraPrivateTx: sender.NewSendTesseraPrivateTxUseCase(ec, txSchedulerClient, chainRegistryURL),
		sendTesseraMarkingTx: sender.NewSendTesseraMarkingTxUseCase(signQuorumTransactionUC, ec, txSchedulerClient, chainRegistryURL, checker),
	}
}

func (u *useCases) SendETHTx() usecases.SendETHTxUseCase {
	return u.sendETHTx
}

func (u *useCases) SendETHRawTx() usecases.SendETHRawTxUseCase {
	return u.sendETHRawTx
}

func (u *useCases) SendEEAPrivateTx() usecases.SendEEAPrivateTxUseCase {
	return u.sendEEAPrivateTx
}

func (u *useCases) SendTesseraPrivateTx() usecases.SendTesseraPrivateTxUseCase {
	return u.sendTesseraPrivateTx
}

func (u *useCases) SendTesseraMarkingTx() usecases.SendTesseraMarkingTxUseCase {
	return u.sendTesseraMarkingTx
}
