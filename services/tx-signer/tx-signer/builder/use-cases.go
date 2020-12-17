package builder

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	keymanagerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"
	nonce2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
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

func NewUseCases(
	client client.OrchestrateClient,
	keyManagerClient keymanagerclient.KeyManagerClient,
	ec ethclient.MultiClient,
	nonceManager nonce2.Sender,
	chainRegistryURL string,
	checkerMaxRecovery uint64,
) usecases.UseCases {
	signETHTransactionUC := signer.NewSignETHTransactionUseCase(keyManagerClient)
	signEEATransactionUC := signer.NewSignEEATransactionUseCase(keyManagerClient)
	signQuorumTransactionUC := signer.NewSignQuorumPrivateTransactionUseCase(keyManagerClient)

	checker := nonce.NewNonceChecker(ec, nonceManager, nonce.NewRecoveryTracker(), chainRegistryURL, checkerMaxRecovery)
	return &useCases{
		sendETHTx:            sender.NewSendEthTxUseCase(signETHTransactionUC, ec, client, chainRegistryURL, checker),
		sendETHRawTx:         sender.NewSendETHRawTxUseCase(ec, client, chainRegistryURL),
		sendEEAPrivateTx:     sender.NewSendEEAPrivateTxUseCase(signEEATransactionUC, ec, client, chainRegistryURL, checker),
		sendTesseraPrivateTx: sender.NewSendTesseraPrivateTxUseCase(ec, client, chainRegistryURL),
		sendTesseraMarkingTx: sender.NewSendTesseraMarkingTxUseCase(signQuorumTransactionUC, ec, client, chainRegistryURL, checker),
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
