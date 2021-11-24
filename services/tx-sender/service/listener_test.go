// +build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/broker/sarama/mock"
	"github.com/consensys/orchestrate/pkg/encoding/proto"
	"github.com/consensys/orchestrate/pkg/errors"
	mock3 "github.com/consensys/orchestrate/pkg/sdk/client/mock"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/tx"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases/mocks"
	"github.com/gofrs/uuid"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type messageListenerCtrlTestSuite struct {
	suite.Suite
	listener           *MessageListener
	producer           *mock.MockSyncProducer
	sendETHUC          *mocks.MockSendETHTxUseCase
	sendETHRawUC       *mocks.MockSendETHRawTxUseCase
	sendEEAPrivateUC   *mocks.MockSendEEAPrivateTxUseCase
	sendTesseraMarking *mocks.MockSendTesseraMarkingTxUseCase
	sendTesseraPrivate *mocks.MockSendTesseraPrivateTxUseCase
	apiClient          *mock3.MockOrchestrateClient
	tenantID           string
	allowedTenants     []string
	senderTopic        string
	recoverTopic       string
}

var _ usecases.UseCases = &messageListenerCtrlTestSuite{}

func (s *messageListenerCtrlTestSuite) SendETHRawTx() usecases.SendETHRawTxUseCase {
	return s.sendETHRawUC
}

func (s *messageListenerCtrlTestSuite) SendETHTx() usecases.SendETHTxUseCase {
	return s.sendETHUC
}

func (s *messageListenerCtrlTestSuite) SendEEAPrivateTx() usecases.SendEEAPrivateTxUseCase {
	return s.sendEEAPrivateUC
}

func (s *messageListenerCtrlTestSuite) SendTesseraPrivateTx() usecases.SendTesseraPrivateTxUseCase {
	return s.sendTesseraPrivate
}

func (s *messageListenerCtrlTestSuite) SendTesseraMarkingTx() usecases.SendTesseraMarkingTxUseCase {
	return s.sendTesseraMarking
}

func TestMessageListener(t *testing.T) {
	s := new(messageListenerCtrlTestSuite)
	suite.Run(t, s)
}

func (s *messageListenerCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.tenantID = "tenantID"
	s.allowedTenants = []string{s.tenantID, "_"}
	s.sendETHRawUC = mocks.NewMockSendETHRawTxUseCase(ctrl)
	s.sendETHUC = mocks.NewMockSendETHTxUseCase(ctrl)
	s.sendEEAPrivateUC = mocks.NewMockSendEEAPrivateTxUseCase(ctrl)
	s.sendTesseraPrivate = mocks.NewMockSendTesseraPrivateTxUseCase(ctrl)
	s.sendTesseraMarking = mocks.NewMockSendTesseraMarkingTxUseCase(ctrl)
	s.apiClient = mock3.NewMockOrchestrateClient(ctrl)
	s.senderTopic = "sender-topic"
	s.recoverTopic = "recover-topic"
	s.producer = mock.NewMockSyncProducer()

	bckoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Millisecond*100), 2)
	s.listener = NewMessageListener(s, s.apiClient, s.producer, s.recoverTopic, s.senderTopic, bckoff)
}

func (s *messageListenerCtrlTestSuite) TestMessageListener_PublicEthereum() {
	s.T().Run("should execute use case for multiple public ethereum transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope(s.tenantID)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.sendETHUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(2).Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)
		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
		assert.Nil(t, s.producer.LastMessage())
	})

	s.T().Run("should execute use case for public raw ethereum transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope(s.tenantID)
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.sendETHRawUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
		assert.Nil(t, s.producer.LastMessage())
	})

	s.T().Run("should execute use case for eea transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope(s.tenantID)
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.sendEEAPrivateUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
		assert.Nil(t, s.producer.LastMessage())
	})

	s.T().Run("should execute use case for tessera marking transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope(s.tenantID)
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.sendTesseraMarking.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
		assert.Nil(t, s.producer.LastMessage())
	})

	s.T().Run("should execute use case for multiple tessera private transactions", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		envelope := fakeEnvelope(s.tenantID)
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(envelope.TxEnvelopeAsRequest())

		s.sendTesseraPrivate.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(2).Return(nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)
		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
		assert.Nil(t, s.producer.LastMessage())
	})
}

func (s *messageListenerCtrlTestSuite) TestMessageListener_PublicEthereum_Errors() {
	s.T().Run("should update transaction and send message to tx-recover if sending fails", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		s.producer.Clean()

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		evlp := fakeEnvelope(s.tenantID)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(evlp.TxEnvelopeAsRequest())

		err := errors.InternalError("error")
		s.sendETHUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(err)
		s.apiClient.EXPECT().UpdateJob(gomock.Any(), evlp.GetJobUUID(), &api.UpdateJobRequest{
			Status:      entities.StatusFailed,
			Message:     err.Error(),
			Transaction: nil,
		}).Return(&api.JobResponse{}, nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
		assert.NotNil(t, s.producer.LastMessage())
		assert.Equal(t, s.recoverTopic, s.producer.LastMessage().Topic)
	})

	s.T().Run("should update transaction and retry job if sending fails by nonce error", func(t *testing.T) {
		var claims map[string][]int32
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)

		mockSession := mock.NewConsumerGroupSession(ctx, "kafka-consumer-group", claims)
		mockClaim := mock.NewConsumerGroupClaim("topic", 0, 0)
		evlp := fakeEnvelope(s.tenantID)
		msg := &sarama.ConsumerMessage{}
		msg.Value, _ = proto.Marshal(evlp.TxEnvelopeAsRequest())

		err := errors.InvalidNonceWarning("nonce too low")
		gomock.InOrder(
			s.sendETHUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(err),
			s.sendETHUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil),
		)

		s.apiClient.EXPECT().UpdateJob(gomock.Any(), evlp.GetJobUUID(), &api.UpdateJobRequest{
			Status:      entities.StatusRecovering,
			Message:     err.Error(),
			Transaction: nil,
		}).Return(&api.JobResponse{}, nil)

		cerr := make(chan error)
		go func() {
			cerr <- s.listener.ConsumeClaim(mockSession, mockClaim)
		}()

		mockClaim.ExpectMessage(msg)

		assert.NoError(t, <-cerr)
	})
}

func fakeEnvelope(tenantID string) *tx.Envelope {
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()

	envelope := tx.NewEnvelope()
	_ = envelope.SetID(jobUUID)
	_ = envelope.SetJobUUID(jobUUID)
	_ = envelope.SetScheduleUUID(scheduleUUID)
	_ = envelope.SetNonce(0)
	_ = envelope.SetFromString("0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5")
	_ = envelope.SetGas(21000)
	_ = envelope.SetGasPriceString("10000000")
	_ = envelope.SetValueString("10000000")
	_ = envelope.SetDataString("0x")
	_ = envelope.SetChainIDString("1")
	_ = envelope.SetHeadersValue(utils.TenantIDHeader, tenantID)
	_ = envelope.SetPrivateFrom("A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=")
	_ = envelope.SetPrivateFor([]string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="})

	return envelope
}
