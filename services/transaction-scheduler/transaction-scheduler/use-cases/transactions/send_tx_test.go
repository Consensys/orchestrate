// +build unit

package transactions

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	mocks4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	mocks3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
)

type sendTxSuite struct {
	suite.Suite
	usecase     SendTxUseCase
	DB          *mocks2.MockDB
	DBTX        *mocks2.MockTx
	ORM         *mocks4.MockORM
	Validators  *mocks3.MockTransactionValidator
	TxRequestDA *mocks2.MockTransactionRequestAgent
	StartJobUC  *mocks.MockStartJobUseCase
}

func TestSendTx(t *testing.T) {
	s := new(sendTxSuite)
	suite.Run(t, s)
}

func (s *sendTxSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.DB = mocks2.NewMockDB(ctrl)
	s.DBTX = mocks2.NewMockTx(ctrl)
	s.ORM = mocks4.NewMockORM(ctrl)
	s.Validators = mocks3.NewMockTransactionValidator(ctrl)
	s.TxRequestDA = mocks2.NewMockTransactionRequestAgent(ctrl)
	s.StartJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.usecase = NewSendTxUseCase(s.Validators, s.DB, s.ORM, s.StartJobUC)
}

func (s *sendTxSuite) TestSendTx_Success() {
	ctx := context.Background()

	// We skip the test of DB transaction as it will be unit-tested at /pkg/database:ExecuteInDbTx()
	s.DB.EXPECT().Begin().Return(s.DBTX, nil).AnyTimes()
	s.DBTX.EXPECT().Begin().Return(s.DBTX, nil).AnyTimes()
	s.DBTX.EXPECT().Commit().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Close().Return(nil).AnyTimes()

	tenantID := "tenantID"
	chainUUID := uuid.NewV4().String()

	s.T().Run("should execute use case successfully", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		jsonParams, _ := utils.ObjectToJSON(txRequest.Params)
		fakeSchedule := testutils2.FakeSchedule("tenantID")
		fakeSchedule.ID = 666
		fakeSchedule.ChainUUID = chainUUID
		txRequestModel := &models.TransactionRequest{
			IdempotencyKey: txRequest.IdempotencyKey,
			ScheduleID:     &fakeSchedule.ID,
			RequestHash:    "requestHash",
			Params:         jsonParams,
		}

		s.Validators.EXPECT().ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		s.ORM.EXPECT().InsertOrUpdateJob(gomock.Any(), s.DBTX, gomock.Any()).DoAndReturn(
			func(ctx context.Context, db store.DB, job *models.Job) error {
				job.Schedule = fakeSchedule
				return nil
			})
		s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA)
		s.TxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Eq(txRequestModel)).Return(nil)
		s.StartJobUC.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(nil)
		s.ORM.EXPECT().FetchScheduleByID(ctx, gomock.Eq(s.DB), gomock.Any()).Return(fakeSchedule, nil)

		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)

		timeNow := time.Now()
		expectedResponse := &types.TransactionResponse{
			IdempotencyKey: txRequest.IdempotencyKey,
			Schedule: &types.ScheduleResponse{
				UUID:      fakeSchedule.UUID,
				ChainUUID: chainUUID,
				CreatedAt: timeNow,
			},
			CreatedAt: timeNow,
		}
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.IdempotencyKey, txResponse.IdempotencyKey)
		assert.Equal(t, expectedResponse.Schedule.UUID, txResponse.Schedule.UUID)
		assert.Equal(t, expectedResponse.Schedule.ChainUUID, txResponse.Schedule.ChainUUID)
		assert.Equal(t, types.JobStatusStarted, types.JobStatusStarted)
		assert.Equal(t, expectedResponse.Schedule.CreatedAt, timeNow)
		assert.Equal(t, expectedResponse.CreatedAt, timeNow)
	})
}

func (s *sendTxSuite) TestSendTx_ExpectedErrors() {
	ctx := context.Background()
	tenantID := "tenantID"
	chainUUID := uuid.NewV4().String()
	
	s.T().Run("should fail with same error if validator fails to validate request hash", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.InvalidParameterError("error")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return("", expectedErr)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if validator fails to validate chain", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.InvalidParameterError("error")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(expectedErr)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	
	s.T().Run("should fail with same error if Begin fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		s.DB.EXPECT().Begin().Return(nil, expectedErr)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	 
	s.T().Run("should fail with same error if Insert Job fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Rollback().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)
		s.ORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(s.DBTX), gomock.Any()).Return(expectedErr)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	
	s.T().Run("should fail with same error if Insert TxRequest fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.PostgresConnectionError("error")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Rollback().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)
		s.ORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(s.DBTX), gomock.Any()).Return(nil)
		s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA)
		s.TxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(expectedErr)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	
	s.T().Run("should fail with same error if start job fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.KafkaConnectionError("error")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Rollback().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)
		s.ORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(s.DBTX), gomock.Any()).Return(nil)
		s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA)
		s.TxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(expectedErr)
		s.StartJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr)
		// s.ORM.EXPECT().FetchScheduleByID(ctx, mockDB, fakeSchedule.ID).Return(fakeSchedule, nil)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	
	s.T().Run("should fail with same error if start job fails", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedErr := errors.KafkaConnectionError("error")
		fakeSchedule := testutils2.FakeSchedule("tenantID")
	
		s.Validators.EXPECT().ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).Return("requestHash", nil)
		s.Validators.EXPECT().ValidateChainExists(ctx, gomock.Any()).Return(nil)
		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Rollback().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)
		s.ORM.EXPECT().InsertOrUpdateJob(gomock.Any(), gomock.Eq(s.DBTX), gomock.Any()).Return(nil)
		s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA)
		s.TxRequestDA.EXPECT().SelectOrInsert(ctx, gomock.Any()).Return(expectedErr)
		s.StartJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.ORM.EXPECT().FetchScheduleByID(ctx, gomock.Any(), gomock.Any()).Return(fakeSchedule, nil)
	
		txResponse, err := s.usecase.Execute(ctx, txRequest, chainUUID, tenantID)
	
		assert.Nil(t, txResponse)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}
