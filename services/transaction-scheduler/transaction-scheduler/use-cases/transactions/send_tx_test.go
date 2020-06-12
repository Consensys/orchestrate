// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
	mocks4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules/mocks"
	mocks3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
)

type sendTxSuite struct {
	suite.Suite
	usecase          SendTxUseCase
	DB               *mocks2.MockDB
	DBTX             *mocks2.MockTx
	Validators       *mocks3.MockTransactionValidator
	TxRequestDA      *mocks2.MockTransactionRequestAgent
	StartJobUC       *mocks.MockStartJobUseCase
	CreateJobUC      *mocks.MockCreateJobUseCase
	CreateScheduleUC *mocks4.MockCreateScheduleUseCase
	GetScheduleUC    *mocks4.MockGetScheduleUseCase
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
	s.Validators = mocks3.NewMockTransactionValidator(ctrl)
	s.TxRequestDA = mocks2.NewMockTransactionRequestAgent(ctrl)
	s.StartJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.CreateJobUC = mocks.NewMockCreateJobUseCase(ctrl)
	s.CreateScheduleUC = mocks4.NewMockCreateScheduleUseCase(ctrl)
	s.GetScheduleUC = mocks4.NewMockGetScheduleUseCase(ctrl)
	s.usecase = NewSendTxUseCase(s.Validators, s.DB, s.StartJobUC, s.CreateJobUC, s.CreateScheduleUC, s.GetScheduleUC)
}

func (s *sendTxSuite) TestSendTx_Success() {
	jobUUID := uuid.NewV4().String()
	scheduleUUID := uuid.NewV4().String()

	s.T().Run("should execute send successfully a public tx", func(t *testing.T) {
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, types.EthereumTransaction)
		assert.Nil(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a orion tx", func(t *testing.T) {
		txRequest := testutils.FakeOrionTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = utils.OrionChainType

		response, err := successfulTestExecution(s, txRequest, types.OrionEEATransaction)
		assert.Nil(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a tessera tx", func(t *testing.T) {
		txRequest := testutils.FakeTesseraTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = utils.TesseraChainType

		response, err := successfulTestExecution(s, txRequest, types.TesseraPrivateTransaction)
		assert.Nil(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})
	
	s.T().Run("should execute send successfully a raw tx", func(t *testing.T) {
		txRequest := testutils.FakeRawTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
	
		response, err := successfulTestExecution(s, txRequest, types.EthereumRawTransaction)
		assert.Nil(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})
}

func (s *sendTxSuite) TestSendTx_ExpectedErrors() {
	ctx := context.Background()

	tenantID := "tenantID"
	requestHash := "requestHash"
	chainUUID := uuid.NewV4().String()
	jobUUID := uuid.NewV4().String()
	scheduleUUID := uuid.NewV4().String()
	txData := ""

	s.T().Run("should fail with same error if validator fails to validate fields", func(t *testing.T) {
		expectedErr := errors.InvalidParameterError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if validator fails to validate request hash", func(t *testing.T) {
		expectedErr := errors.InvalidParameterError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)
		
		s.Validators.EXPECT().
			ValidateRequestHash(gomock.Any(), chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return(requestHash, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if select or insert txRequest fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)

		s.Validators.EXPECT().
			ValidateRequestHash(gomock.Any(), chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return(requestHash, nil)

		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Close().Return(nil)
		s.DBTX.EXPECT().Rollback().Return(nil)

		s.DBTX.EXPECT().
			TransactionRequest().
			Return(s.TxRequestDA)

		s.TxRequestDA.EXPECT().
			SelectOrInsert(ctx, gomock.Any()).
			Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createSchedule UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)

		s.Validators.EXPECT().
			ValidateRequestHash(gomock.Any(), chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return(requestHash, nil)

		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Rollback().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)

		s.DBTX.EXPECT().
			TransactionRequest().
			Return(s.TxRequestDA)

		s.TxRequestDA.EXPECT().
			SelectOrInsert(ctx, gomock.Any()).
			Return(nil)

		s.CreateScheduleUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateScheduleUC)

		s.CreateScheduleUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData = ""

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)

		s.Validators.EXPECT().
			ValidateRequestHash(gomock.Any(), chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return(requestHash, nil)

		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Rollback().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)

		s.DBTX.EXPECT().TransactionRequest().
			Return(s.TxRequestDA).Times(1)

		s.TxRequestDA.EXPECT().
			SelectOrInsert(ctx, gomock.Any()).
			Return(nil)

		s.CreateScheduleUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateScheduleUC)

		s.CreateScheduleUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule, nil)

		s.CreateJobUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateJobUC)

		s.CreateJobUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule.Jobs[0], expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if startJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)

		s.Validators.EXPECT().
			ValidateRequestHash(gomock.Any(), chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return(requestHash, nil)

		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Commit().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)

		s.DBTX.EXPECT().TransactionRequest().
			Return(s.TxRequestDA).Times(1)

		s.TxRequestDA.EXPECT().
			SelectOrInsert(ctx, gomock.Any()).
			Return(nil)

		s.CreateScheduleUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateScheduleUC)

		s.CreateScheduleUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule, nil)

		s.CreateJobUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateJobUC)

		s.CreateJobUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule.Jobs[0], nil)

		s.StartJobUC.EXPECT().
			Execute(ctx, jobUUID, tenantID).
			Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if getSchedule UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)
		
		s.Validators.EXPECT().
			ValidateRequestHash(gomock.Any(), chainUUID, txRequest.Params, txRequest.IdempotencyKey).
			Return(requestHash, nil)
		
		s.DB.EXPECT().Begin().Return(s.DBTX, nil)
		s.DBTX.EXPECT().Commit().Return(nil)
		s.DBTX.EXPECT().Close().Return(nil)

		s.DBTX.EXPECT().TransactionRequest().
			Return(s.TxRequestDA).Times(1)

		s.TxRequestDA.EXPECT().
			SelectOrInsert(ctx, gomock.Any()).
			Return(nil)

		s.CreateScheduleUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateScheduleUC)

		s.CreateScheduleUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule, nil)

		s.CreateJobUC.EXPECT().
			WithDBTransaction(s.DBTX).
			Return(s.CreateJobUC)

		s.CreateJobUC.EXPECT().
			Execute(gomock.Any(), gomock.Any(), tenantID).
			Return(txRequest.Schedule.Jobs[0], nil)

		s.StartJobUC.EXPECT().
			Execute(ctx, jobUUID, tenantID).
			Return(nil)

		s.GetScheduleUC.EXPECT().
			Execute(ctx, scheduleUUID, tenantID).
			Return(txRequest.Schedule, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}

func successfulTestExecution(s *sendTxSuite, txRequest *entities.TxRequest, jobType string) (*entities.TxRequest, error) {
	ctx := context.Background()
	tenantID := "tenantID"
	requestHash := "requestHash"
	chainUUID := uuid.NewV4().String()
	scheduleUUID := txRequest.Schedule.UUID
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	txData := ""

	s.Validators.EXPECT().
			ValidateFields(gomock.Any(), txRequest).
			Return(nil)

	s.Validators.EXPECT().
		ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey).
		Return(requestHash, nil)

	s.DB.EXPECT().Begin().Return(s.DBTX, nil).Times(1)
	s.DBTX.EXPECT().Commit().Return(nil).Times(1)
	s.DBTX.EXPECT().Close().Return(nil).Times(1)

	s.DBTX.EXPECT().TransactionRequest().
		Return(s.TxRequestDA).Times(1)

	s.TxRequestDA.EXPECT().
		SelectOrInsert(ctx, gomock.Any()).
		Return(nil)

	s.CreateScheduleUC.EXPECT().
		WithDBTransaction(s.DBTX).
		Return(s.CreateScheduleUC)

	s.CreateScheduleUC.EXPECT().
		Execute(gomock.Any(), gomock.Any(), tenantID).
		Return(txRequest.Schedule, nil)

	s.CreateJobUC.EXPECT().
		WithDBTransaction(s.DBTX).
		Return(s.CreateJobUC)

	s.CreateJobUC.EXPECT().
		Execute(gomock.Any(), gomock.Any(), tenantID).
		DoAndReturn(func(ctx context.Context, jobEntity *types.Job, tenantID string) (*types.Job, error) {
			if jobEntity.Type != jobType {
				return nil, fmt.Errorf("invalid job type")
			}
			return txRequest.Schedule.Jobs[0], nil
		})

	s.StartJobUC.EXPECT().
		Execute(ctx, jobUUID, tenantID).
		Return(nil)

	s.GetScheduleUC.EXPECT().
		Execute(ctx, scheduleUUID, tenantID).
		Return(txRequest.Schedule, nil)

	return s.usecase.Execute(ctx, txRequest, txData, chainUUID, tenantID)
}
