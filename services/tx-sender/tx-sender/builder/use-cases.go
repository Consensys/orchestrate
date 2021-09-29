package builder

import (
	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	keymanager "github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/nonce"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases/crafter"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases/sender"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases/signer"
)

type useCases struct {
	sendETHTx            usecases.SendETHTxUseCase
	sendETHRawTx         usecases.SendETHRawTxUseCase
	sendEEAPrivateTx     usecases.SendEEAPrivateTxUseCase
	sendTesseraPrivateTx usecases.SendTesseraPrivateTxUseCase
	sendTesseraMarkingTx usecases.SendTesseraMarkingTxUseCase
}

func NewUseCases(jobClient client.JobClient, keyManagerClient keymanager.KeyManagerClient,
	ec ethclient.MultiClient, nonceManager nonce.Manager, chainRegistryURL string, checkerMaxRecovery uint64) usecases.UseCases {
	signETHTransactionUC := signer.NewSignETHTransactionUseCase(keyManagerClient)
	signEEATransactionUC := signer.NewSignEEATransactionUseCase(keyManagerClient)
	signQuorumTransactionUC := signer.NewSignQuorumPrivateTransactionUseCase(keyManagerClient)

	crafterUC := crafter.NewCraftTransactionUseCase(ec, chainRegistryURL, nonceManager)

	return &useCases{
		sendETHTx:            sender.NewSendEthTxUseCase(signETHTransactionUC, crafterUC, ec, jobClient, chainRegistryURL, nonceManager),
		sendETHRawTx:         sender.NewSendETHRawTxUseCase(ec, jobClient, chainRegistryURL),
		sendEEAPrivateTx:     sender.NewSendEEAPrivateTxUseCase(signEEATransactionUC, crafterUC, ec, jobClient, chainRegistryURL, nonceManager),
		sendTesseraPrivateTx: sender.NewSendTesseraPrivateTxUseCase(ec, crafterUC, jobClient, chainRegistryURL),
		sendTesseraMarkingTx: sender.NewSendTesseraMarkingTxUseCase(signQuorumTransactionUC, crafterUC, ec, jobClient, chainRegistryURL, nonceManager),
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
