// +build unit

package service

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer/tx-signer/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-signer/tx-signer/use-cases/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type messageListenerCtrlTestSuite struct {
	suite.Suite
	listener              *MessageListener
	signTxUC              *mocks.MockSignTransactionUseCase
	signEEATxUC           *mocks.MockSignEEATransactionUseCase
	signQuorumPrivateTxUC *mocks.MockSignQuorumPrivateTransactionUseCase
	sendEnvelopeUC        *mocks.MockSendEnvelopeUseCase
	txSchedulerClient     *mock2.MockTransactionSchedulerClient
	tenants               []string
	senderTopic           string
	recoverTopic          string
}

var _ usecases.EthereumUseCases = &messageListenerCtrlTestSuite{}

func (s messageListenerCtrlTestSuite) SignTransaction() usecases.SignTransactionUseCase {
	return s.signTxUC
}

func (s messageListenerCtrlTestSuite) SignEEATransaction() usecases.SignEEATransactionUseCase {
	return s.signEEATxUC
}

func (s messageListenerCtrlTestSuite) SignQuorumPrivateTransaction() usecases.SignQuorumPrivateTransactionUseCase {
	return s.signQuorumPrivateTxUC
}

func (s messageListenerCtrlTestSuite) SendEnvelope() usecases.SendEnvelopeUseCase {
	return s.sendEnvelopeUC
}

func TestMessageListener(t *testing.T) {
	s := new(messageListenerCtrlTestSuite)
	suite.Run(t, s)
}

func (s *messageListenerCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenants = []string{"tenantID"}
	s.signTxUC = mocks.NewMockSignTransactionUseCase(ctrl)
	s.signEEATxUC = mocks.NewMockSignEEATransactionUseCase(ctrl)
	s.signQuorumPrivateTxUC = mocks.NewMockSignQuorumPrivateTransactionUseCase(ctrl)
	s.sendEnvelopeUC = mocks.NewMockSendEnvelopeUseCase(ctrl)
	s.txSchedulerClient = mock2.NewMockTransactionSchedulerClient(ctrl)
	s.senderTopic = "sender-topic"
	s.recoverTopic = "recover-topic"

	s.listener = NewMessageListener(s, s.senderTopic, s.recoverTopic, s.txSchedulerClient)
}

func (s *messageListenerCtrlTestSuite) TestMessageListener_PublicEthereum() {
	s.T().Run("should execute use case for multiple public ethereum transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope()
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())
		raw := "0xraw"
		txHash := "0xhash"

		s.signTxUC.EXPECT().Execute(ctx, gomock.Any()).Return(raw, txHash, nil).Times(2)
		s.sendEnvelopeUC.EXPECT().Execute(ctx, gomock.Any(), s.senderTopic, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5@1").Return(nil).Times(2)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)
		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})

	s.T().Run("should execute use case for multiple eea transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())
		raw := "0xraw"
		txHash := "0xhash"
		expectedPartitionKey := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5@orion-d24bc1408ff6fb773a0d30588f37553e@1"

		s.signEEATxUC.EXPECT().Execute(ctx, gomock.Any()).Return(raw, txHash, nil).Times(2)
		s.sendEnvelopeUC.EXPECT().Execute(ctx, gomock.Any(), s.senderTopic, expectedPartitionKey).Return(nil).Times(2)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)
		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})

	s.T().Run("should execute use case for multiple quorum private transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())
		raw := "0xraw"
		txHash := "0xhash"
		expectedPartitionKey := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5@1"

		s.signQuorumPrivateTxUC.EXPECT().Execute(ctx, gomock.Any()).Return(raw, txHash, nil).Times(2)
		s.sendEnvelopeUC.EXPECT().Execute(ctx, gomock.Any(), s.senderTopic, expectedPartitionKey).Return(nil).Times(2)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)
		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})

	s.T().Run("should update transaction and send message to tx-recover if signing fails", func(t *testing.T) {
		var claims map[string][]int32
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope()
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.signTxUC.EXPECT().Execute(ctx, gomock.Any()).Return("", "", errors.InternalError("error"))
		s.txSchedulerClient.EXPECT().UpdateJob(ctx, envelope.GetJobUUID(), &txschedulertypes.UpdateJobRequest{
			Status:  utils.StatusFailed,
			Message: "FF000@: error",
		}).Return(&txschedulertypes.JobResponse{}, nil)
		s.sendEnvelopeUC.EXPECT().
			Execute(ctx, gomock.Any(), s.recoverTopic, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5@1").
			Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})

	s.T().Run("should not fail tx and send envelope to tx-recover on ConnectionError", func(t *testing.T) {
		var claims map[string][]int32
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope()
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.signTxUC.EXPECT().Execute(ctx, gomock.Any()).Return("", "", errors.KafkaConnectionError("error"))

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})

	s.T().Run("should skip signing for raw transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		_ = envelope.SetRawString("0xf851018227108252088082c35080820713a09a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d52a065bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b")
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.sendEnvelopeUC.EXPECT().
			Execute(ctx, gomock.Any(), s.senderTopic, "").
			Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})
}

func fakeEnvelope() *tx.Envelope {
	jobUUID := uuid.Must(uuid.NewV4()).String()

	envelope := tx.NewEnvelope()
	_ = envelope.SetID(jobUUID)
	_ = envelope.SetJobUUID(jobUUID)
	_ = envelope.SetNonce(0)
	_ = envelope.SetFromString("0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5")
	_ = envelope.SetGas(21000)
	_ = envelope.SetGasPriceString("10000000")
	_ = envelope.SetValueString("10000000")
	_ = envelope.SetDataString("0x")
	_ = envelope.SetChainIDString("1")
	_ = envelope.SetHeadersValue(multitenancy.TenantIDMetadata, "tenantID")
	_ = envelope.SetPrivateFrom("A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=")
	_ = envelope.SetPrivateFor([]string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="})

	return envelope
}
