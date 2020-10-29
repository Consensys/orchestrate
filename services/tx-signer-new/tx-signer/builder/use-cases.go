package builder

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer-new/tx-signer/use-cases/ethereum"
)

type useCases struct {
	signTransaction    usecases.SignTransactionUseCase
	signEEATransaction usecases.SignEEATransactionUseCase
	sendEnvelope       usecases.SendEnvelopeUseCase
}

func NewUseCases(keyManagerClient client.KeyManagerClient, producer sarama.SyncProducer) usecases.UseCases {
	return &useCases{
		signTransaction:    ethereum.NewSignTransactionUseCase(keyManagerClient),
		signEEATransaction: ethereum.NewSignEEATransactionUseCase(keyManagerClient),
		sendEnvelope:       ethereum.NewSendEnvelopeUseCase(producer),
	}
}

func (u *useCases) SignTransaction() usecases.SignTransactionUseCase {
	return u.signTransaction
}

func (u *useCases) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return u.signEEATransaction
}

func (u *useCases) SendEnvelope() usecases.SendEnvelopeUseCase {
	return u.sendEnvelope
}
