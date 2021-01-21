// +build unit

package jobs

import (
	"context"
	"fmt"
	"testing"

	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/metrics/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics/mock"

	mocks2 "github.com/Shopify/sarama/mocks"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
)

func TestStartJob_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobDA := mocks.NewMockJobAgent(ctrl)
	mockLogDA := mocks.NewMockLogAgent(ctrl)
	mockDBTX := mocks.NewMockTx(ctrl)
	mockKafkaProducer := mocks2.NewSyncProducer(t, nil)
	tenants := []string{"tenantID"}
	mockMetrics := mock.NewMockTransactionSchedulerMetrics(ctrl)

	jobsLatencyHistogram := mock2.NewMockHistogram(ctrl)
	jobsLatencyHistogram.EXPECT().With(gomock.Any()).AnyTimes().Return(jobsLatencyHistogram)
	jobsLatencyHistogram.EXPECT().Observe(gomock.Any()).AnyTimes()
	mockMetrics.EXPECT().JobsLatencyHistogram().AnyTimes().Return(jobsLatencyHistogram)

	mockDB := mocks.NewMockDB(ctrl)
	mockDB.EXPECT().Begin().Return(mockDBTX, nil).AnyTimes()

	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDB.EXPECT().Job().Return(mockJobDA).AnyTimes()
	mockDBTX.EXPECT().Log().Return(mockLogDA).AnyTimes()

	usecase := NewStartJobUseCase(mockDB, mockKafkaProducer, sarama.NewKafkaTopicConfig(viper.GetViper()), mockMetrics)

	t.Run("should execute use case successfully", func(t *testing.T) {
		job := testutils.FakeJobModel(1)
		job.ID = 1
		job.UUID = "6380e2b6-b828-43ee-abdc-de0f8d57dc5f"
		job.Transaction.Sender = "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"
		job.Schedule = testutils.FakeSchedule("")

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), job.UUID, tenants).Return(job, nil)
		mockKafkaProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
			txEnvelope := &tx.TxEnvelope{}
			err := encoding.Unmarshal(val, txEnvelope)
			if err != nil {
				return err
			}
			envelope, err := txEnvelope.Envelope()
			if err != nil {
				return err
			}

			assert.Equal(t, envelope.GetJobUUID(), job.UUID)
			assert.False(t, envelope.IsOneTimeKeySignature())
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockDBTX.EXPECT().Commit().Return(nil)
		err := usecase.Execute(ctx, job.UUID, tenants)

		assert.NoError(t, err)
	})

	t.Run("should execute use case with one-time-key successfully", func(t *testing.T) {
		job := testutils.FakeJobModel(1)
		job.ID = 1
		job.UUID = "6380e2b6-b828-43ee-abdc-de0f8d57dc5f"
		job.Transaction.Sender = "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"
		job.Schedule = testutils.FakeSchedule("")
		job.InternalData = &entities.InternalData{
			OneTimeKey: true,
		}

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), job.UUID, tenants).Return(job, nil)
		mockKafkaProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {
			txEnvelope := &tx.TxEnvelope{}
			err := encoding.Unmarshal(val, txEnvelope)
			if err != nil {
				return err
			}
			envelope, err := txEnvelope.Envelope()
			if err != nil {
				return err
			}

			assert.Equal(t, envelope.GetJobUUID(), job.UUID)
			assert.True(t, envelope.IsOneTimeKeySignature())
			return nil
		})
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockDBTX.EXPECT().Commit().Return(nil)
		err := usecase.Execute(ctx, job.UUID, tenants)

		assert.NoError(t, err)
	})

	t.Run("should fail with same error if FindOne fails", func(t *testing.T) {
		job := testutils.FakeJobModel(1)
		job.UUID = "6380e2b6-b828-43ee-abdc-de0f8d57dc5f"
		expectedErr := errors.NotFoundError("error")

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), job.UUID, tenants).Return(nil, expectedErr)

		err := usecase.Execute(ctx, job.UUID, tenants)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(startJobComponent), err)
	})

	t.Run("should fail with same error if Insert log fails", func(t *testing.T) {
		job := testutils.FakeJobModel(1)
		job.ID = 1
		job.UUID = "6380e2b6-b828-43ee-abdc-de0f8d57dc5f"
		job.Transaction.Sender = "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"
		job.Transaction.ID = 1
		job.Schedule = testutils.FakeSchedule("")
		job.Schedule.ID = 1
		expectedErr := errors.PostgresConnectionError("error")

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), job.UUID, tenants).Return(job, nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		err := usecase.Execute(ctx, job.UUID, tenants)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(startJobComponent), err)
	})

	t.Run("should fail with KafkaConnectionError if Produce fails", func(t *testing.T) {
		job := testutils.FakeJobModel(1)
		job.UUID = "6380e2b6-b828-43ee-abdc-de0f8d57dc5f"
		job.Transaction.Sender = "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"
		job.Schedule = testutils.FakeSchedule("")

		mockJobDA.EXPECT().FindOneByUUID(gomock.Any(), job.UUID, tenants).Return(job, nil)
		mockLogDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		mockKafkaProducer.ExpectSendMessageAndFail(fmt.Errorf("error"))
		mockDBTX.EXPECT().Rollback().Return(nil)
		err := usecase.Execute(ctx, job.UUID, tenants)
		assert.True(t, errors.IsKafkaConnectionError(err))
	})
}
