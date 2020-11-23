// +build unit

package ethereum

import (
	"context"
	"github.com/Shopify/sarama/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

func TestSendEnvelope_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewSyncProducer(t, nil)
	senderTopic := "tx-sender-topic"
	ctx := context.Background()
	fakeEnvelope := tx.NewEnvelope()

	usecase := NewSendEnvelopeUseCase(mockProducer)

	t.Run("should execute use case successfully", func(t *testing.T) {
		mockProducer.ExpectSendMessageAndSucceed()

		err := usecase.Execute(ctx, fakeEnvelope.TxEnvelopeAsRequest(), senderTopic, fakeEnvelope.PartitionKey())

		assert.NoError(t, err)
	})

	t.Run("should fail with KafkaConnectionError error if SendMessage fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		mockProducer.ExpectSendMessageAndFail(expectedErr)

		err := usecase.Execute(ctx, fakeEnvelope.TxResponse(), senderTopic, fakeEnvelope.PartitionKey())

		assert.True(t, errors.IsKafkaConnectionError(err))
	})
}
