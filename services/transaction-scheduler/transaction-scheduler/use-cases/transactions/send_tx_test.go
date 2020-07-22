// +build unit

package transactions

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	models2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs/mocks"
	mocks4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules/mocks"
	mocks5 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions/mocks"
	mocks3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators/mocks"
	"testing"
	"time"
)

type sendTxSuite struct {
	suite.Suite
	usecase             SendTxUseCase
	DB                  *mocks2.MockDB
	DBTX                *mocks2.MockTx
	Validators          *mocks3.MockTransactionValidator
	ChainRegistryClient *mock.MockChainRegistryClient
	TxRequestDA         *mocks2.MockTransactionRequestAgent
	ScheduleDA          *mocks2.MockScheduleAgent
	StartJobUC          *mocks.MockStartJobUseCase
	CreateJobUC         *mocks.MockCreateJobUseCase
	CreateScheduleUC    *mocks4.MockCreateScheduleUseCase
	GetTxUC             *mocks5.MockGetTxUseCase
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
	s.ScheduleDA = mocks2.NewMockScheduleAgent(ctrl)
	s.StartJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.CreateJobUC = mocks.NewMockCreateJobUseCase(ctrl)
	s.CreateScheduleUC = mocks4.NewMockCreateScheduleUseCase(ctrl)
	s.GetTxUC = mocks5.NewMockGetTxUseCase(ctrl)
	s.ChainRegistryClient = mock.NewMockChainRegistryClient(ctrl)

	s.DB.EXPECT().Begin().Return(s.DBTX, nil).AnyTimes()
	s.DB.EXPECT().TransactionRequest().Return(s.TxRequestDA).AnyTimes()
	s.DBTX.EXPECT().Schedule().Return(s.ScheduleDA).AnyTimes()
	s.DBTX.EXPECT().Commit().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Close().Return(nil).AnyTimes()
	s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA).AnyTimes()
	s.CreateScheduleUC.EXPECT().WithDBTransaction(s.DBTX).Return(s.CreateScheduleUC).AnyTimes()
	s.CreateJobUC.EXPECT().WithDBTransaction(s.DBTX).Return(s.CreateJobUC).AnyTimes()

	s.usecase = NewSendTxUseCase(s.Validators, s.DB, s.ChainRegistryClient, s.StartJobUC, s.CreateJobUC, s.CreateScheduleUC, s.GetTxUC)
}

func (s *sendTxSuite) TestSendTx_Success() {
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()

	s.T().Run("should execute send successfully a public tx", func(t *testing.T) {
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, utils.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, response.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a orion tx", func(t *testing.T) {
		txRequest := testutils.FakeOrionTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = utils.OrionChainType

		response, err := successfulTestExecution(s, txRequest, utils.OrionEEATransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a tessera tx", func(t *testing.T) {
		txRequest := testutils.FakeTesseraTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = utils.TesseraChainType

		response, err := successfulTestExecution(s, txRequest, utils.TesseraPrivateTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a raw tx", func(t *testing.T) {
		txRequest := testutils.FakeRawTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, utils.EthereumRawTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should not insert and start job in DB if TxRequest already exists and send if status is CREATED", func(t *testing.T) {
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		ctx := context.Background()
		tenantID := "tenantID"
		tenants := []string{tenantID}
		chain := &models2.Chain{UUID: "32da9731-0fb8-4235-a4cc-35070ffe5bf0"}
		jobUUID := txRequest.Schedule.Jobs[0].UUID
		txData := ""
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "ea2d3e36db863014fdfaf49a88c31f1d"

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequestModel.UUID, tenants).Return(txRequest, nil)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequest.UUID, tenants).Return(txRequest, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, response.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should not insert and not start job if TxRequest already exists and not send if status is not CREATED", func(t *testing.T) {
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Schedule.Jobs[0].Logs = append(txRequest.Schedule.Jobs[0].Logs, &types.Log{
			Status:    utils.StatusStarted,
			Message:   "already started, do not resend",
			CreatedAt: time.Now(),
		})
		ctx := context.Background()
		tenantID := "tenantID"
		tenants := []string{tenantID}
		chain := &models2.Chain{UUID: "32da9731-0fb8-4235-a4cc-35070ffe5bf0"}
		txData := ""
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "ea2d3e36db863014fdfaf49a88c31f1d"

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequestModel.UUID, tenants).Return(txRequest, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequest.UUID, tenants).Return(txRequest, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, response.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a oneTimeKey tx", func(t *testing.T) {
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.From = ""
		txRequest.Annotations.OneTimeKey = true

		response, err := successfulTestExecution(s, txRequest, utils.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.UUID, response.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.True(t, response.Schedule.Jobs[0].Annotations.OneTimeKey)
	})
}

func (s *sendTxSuite) TestSendTx_ExpectedErrors() {
	ctx := context.Background()

	tenantID := "tenantID"
	tenants := []string{tenantID}
	chain := &models2.Chain{UUID: "32da9731-0fb8-4235-a4cc-35070ffe5bf0"}
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	txData := ""

	s.T().Run("should fail with same error if chain registry client fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with InvalidParameterError if chain registry client fails with NotFoundError", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with same error if validator fails to validate fields", func(t *testing.T) {
		expectedErr := errors.InvalidParameterError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with InvalidParameter error if from account and OneTimeKey is enabled", func(t *testing.T) {
		expectedErr := errors.InvalidParameterError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Annotations.OneTimeKey = true

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with AlreadyExistsError if request found has different request hash", func(t *testing.T) {
		expectedErr := errors.AlreadyExistsError("a transaction request with the same idempotency key and different params already exists")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(&models.TransactionRequest{
			RequestHash: "differentRequestHash",
		}, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createSchedule UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if find schedule fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if select or insert txRequest fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData = ""
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if startJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if getTx UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils.FakeTxRequestEntity()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequest.UUID, tenants).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}

func successfulTestExecution(s *sendTxSuite, txRequest *entities.TxRequest, jobType string) (*entities.TxRequest, error) {
	ctx := context.Background()
	tenantID := "tenantID"
	tenants := []string{"tenantID"}
	chain := &models2.Chain{UUID: "32da9731-0fb8-4235-a4cc-35070ffe5bf0"}
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	txData := ""
	scheduleModel := testutils2.FakeSchedule(tenants[0])

	s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
	s.Validators.EXPECT().ValidateFields(gomock.Any(), txRequest).Return(nil)
	s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
	s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
	s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
	s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
	s.CreateJobUC.EXPECT().
		Execute(gomock.Any(), gomock.Any(), tenants).
		DoAndReturn(func(ctx context.Context, jobEntity *types.Job, tenants []string) (*types.Job, error) {
			if jobEntity.Type != jobType {
				return nil, fmt.Errorf("invalid job type")
			}

			jobEntity.UUID = txRequest.Schedule.Jobs[0].UUID
			return jobEntity, nil
		})
	s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
	s.GetTxUC.EXPECT().Execute(ctx, txRequest.UUID, tenants).Return(txRequest, nil)

	return s.usecase.Execute(ctx, txRequest, txData, tenantID)
}
