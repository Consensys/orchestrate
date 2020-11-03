package builder

import (
	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer/tx-signer/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer/tx-signer/use-cases/ethereum"
)

type useCases struct {
	signTransaction              usecases.SignTransactionUseCase
	signEEATransaction           usecases.SignEEATransactionUseCase
	signQuorumPrivateTransaction usecases.SignQuorumPrivateTransactionUseCase
	sendEnvelope                 usecases.SendEnvelopeUseCase
}

func NewUseCases(keyManagerClient client.KeyManagerClient, producer sarama.SyncProducer) usecases.UseCases {
	return &useCases{
		signTransaction:              ethereum.NewSignTransactionUseCase(keyManagerClient),
		signEEATransaction:           ethereum.NewSignEEATransactionUseCase(keyManagerClient),
		signQuorumPrivateTransaction: ethereum.NewSignQuorumPrivateTransactionUseCase(keyManagerClient),
		sendEnvelope:                 ethereum.NewSendEnvelopeUseCase(producer),
	}
}

func (u *useCases) SignTransaction() usecases.SignTransactionUseCase {
	return u.signTransaction
}

func (u *useCases) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return u.signEEATransaction
}

func (u *useCases) SignQuorumPrivateTransaction() usecases.SignQuorumPrivateTransactionUseCase {
	return u.signQuorumPrivateTransaction
}

func (u *useCases) SendEnvelope() usecases.SendEnvelopeUseCase {
	return u.sendEnvelope
}
