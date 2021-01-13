// +build unit

package transactions

import (
	"context"
	"fmt"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models/testutils"
)

type sendTxSuite struct {
	suite.Suite
	usecase            usecases.SendTxUseCase
	DB                 *mocks2.MockDB
	DBTX               *mocks2.MockTx
	SearchChainsUC     *mocks.MockSearchChainsUseCase
	TxRequestDA        *mocks2.MockTransactionRequestAgent
	ScheduleDA         *mocks2.MockScheduleAgent
	StartJobUC         *mocks.MockStartJobUseCase
	CreateJobUC        *mocks.MockCreateJobUseCase
	CreateScheduleUC   *mocks.MockCreateScheduleUseCase
	GetTxUC            *mocks.MockGetTxUseCase
	GetFaucetCandidate *mocks.MockGetFaucetCandidateUseCase
}

var (
	faucetNotFoundErr = errors.NotFoundError("not found faucet candidate")
)

func TestSendTx(t *testing.T) {
	s := new(sendTxSuite)
	suite.Run(t, s)
}

func (s *sendTxSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.DB = mocks2.NewMockDB(ctrl)
	s.DBTX = mocks2.NewMockTx(ctrl)
	s.SearchChainsUC = mocks.NewMockSearchChainsUseCase(ctrl)
	s.TxRequestDA = mocks2.NewMockTransactionRequestAgent(ctrl)
	s.ScheduleDA = mocks2.NewMockScheduleAgent(ctrl)
	s.StartJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.CreateJobUC = mocks.NewMockCreateJobUseCase(ctrl)
	s.CreateScheduleUC = mocks.NewMockCreateScheduleUseCase(ctrl)
	s.GetTxUC = mocks.NewMockGetTxUseCase(ctrl)
	s.GetFaucetCandidate = mocks.NewMockGetFaucetCandidateUseCase(ctrl)

	s.DB.EXPECT().Begin().Return(s.DBTX, nil).AnyTimes()
	s.DB.EXPECT().TransactionRequest().Return(s.TxRequestDA).AnyTimes()
	s.DBTX.EXPECT().Schedule().Return(s.ScheduleDA).AnyTimes()
	s.DBTX.EXPECT().Commit().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Close().Return(nil).AnyTimes()
	s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA).AnyTimes()
	s.CreateScheduleUC.EXPECT().WithDBTransaction(s.DBTX).Return(s.CreateScheduleUC).AnyTimes()
	s.CreateJobUC.EXPECT().WithDBTransaction(s.DBTX).Return(s.CreateJobUC).AnyTimes()

	s.usecase = NewSendTxUseCase(
		s.DB,
		s.SearchChainsUC,
		s.StartJobUC,
		s.CreateJobUC,
		s.CreateScheduleUC,
		s.GetTxUC,
		s.GetFaucetCandidate,
	)
}

func (s *sendTxSuite) TestSendTx_Success() {
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()

	s.T().Run("should execute send successfully a public tx", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, false, utils.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a public tx with faucet", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, true, utils.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a orion tx", func(t *testing.T) {
		txRequest := testutils3.FakeOrionTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = utils.OrionChainType

		response, err := successfulTestExecution(s, txRequest, false, utils.OrionEEATransaction, utils.OrionMarkingTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a tessera tx", func(t *testing.T) {
		txRequest := testutils3.FakeTesseraTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = utils.TesseraChainType

		response, err := successfulTestExecution(s, txRequest, false, utils.TesseraPrivateTransaction, utils.TesseraMarkingTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a raw tx", func(t *testing.T) {
		txRequest := testutils3.FakeRawTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, false, utils.EthereumRawTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should not insert and start job in DB if TxRequest already exists and send if status is CREATED", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		ctx := context.Background()
		tenantID := "tenantID"
		tenants := []string{tenantID}
		allowedTenants := []string{tenantID, multitenancy.DefaultTenant}
		chains := []*entities.Chain{testutils3.FakeChain()}
		chains[0].UUID = "myChainUUID"
		jobUUID := txRequest.Schedule.Jobs[0].UUID
		txData := ""
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "8ba6ffc20366e5326fc7d4a3f4833306"

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequestModel.Schedule.UUID, tenants).Return(txRequest, nil)
		s.GetFaucetCandidate.EXPECT().Execute(ctx, gomock.Any(), chains[0], allowedTenants).Return(nil, faucetNotFoundErr)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(txRequest, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should not insert and not start job if TxRequest already exists and not send if status is not CREATED", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Schedule.Jobs[0].Status = utils.StatusStarted

		ctx := context.Background()
		tenantID := "tenantID"
		tenants := []string{tenantID}
		txData := ""
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "8ba6ffc20366e5326fc7d4a3f4833306"
		chains := []*entities.Chain{testutils3.FakeChain()}
		chains[0].UUID = "myChainUUID"

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequestModel.Schedule.UUID, tenants).Return(txRequest, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(txRequest, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a oneTimeKey tx", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.From = ""
		txRequest.InternalData.OneTimeKey = true

		response, err := successfulTestExecution(s, txRequest, false, utils.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.True(t, response.Schedule.Jobs[0].InternalData.OneTimeKey)
	})
}

func (s *sendTxSuite) TestSendTx_ExpectedErrors() {
	ctx := context.Background()

	tenantID := "tenantID"
	tenants := []string{tenantID}
	allowedTenants := []string{tenantID, multitenancy.DefaultTenant}
	chains := []*entities.Chain{testutils3.FakeChain()}
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	txData := ""

	s.T().Run("should fail with same error if chain agent fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(ctx, gomock.Any(), []string{tenantID, multitenancy.DefaultTenant}).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with InvalidParameterError if no chain is found", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(ctx, gomock.Any(), []string{tenantID, multitenancy.DefaultTenant}).Return([]*entities.Chain{}, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with AlreadyExistsError if request found has different request hash", func(t *testing.T) {
		expectedErr := errors.AlreadyExistsError("a transaction request with the same idempotency key and different params already exists")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(&models.TransactionRequest{
			RequestHash: "differentRequestHash",
		}, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createSchedule UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if find schedule fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if select or insert txRequest fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
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
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData = ""
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, allowedTenants).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), txRequest.Params.From, chains[0], allowedTenants).Return(nil, faucetNotFoundErr)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), allowedTenants).Return(txRequest.Schedule.Jobs[0], expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if getFaucetCandidate request fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		chains := []*entities.Chain{testutils3.FakeChain()}
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData = ""
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), allowedTenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), txRequest.Params.From, gomock.Any(), allowedTenants).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if startJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), allowedTenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), txRequest.Params.From, gomock.Any(), allowedTenants).Return(nil, faucetNotFoundErr)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, allowedTenants).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if getTx UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(tenants[0])

		s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.GetFaucetCandidate.EXPECT().
			Execute(gomock.Any(), txRequest.Params.From, gomock.Any(), allowedTenants).
			Return(nil, faucetNotFoundErr)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), allowedTenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, allowedTenants).Return(nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}

func successfulTestExecution(s *sendTxSuite, txRequest *entities.TxRequest, withFaucet bool, jobTypes ...string) (*entities.TxRequest, error) {
	ctx := context.Background()
	tenantID := "tenantID"
	tenants := []string{"tenantID"}
	allowedTenants := []string{tenantID, multitenancy.DefaultTenant}
	chains := []*entities.Chain{testutils3.FakeChain()}
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	txData := ""
	scheduleModel := testutils2.FakeSchedule(tenants[0])
	jobIdx := 0

	s.SearchChainsUC.EXPECT().Execute(ctx, &entities.ChainFilters{Names: []string{txRequest.ChainName}}, []string{tenantID, multitenancy.DefaultTenant}).Return(chains, nil)
	s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
	s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
	s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
	s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
	s.CreateJobUC.EXPECT().
		Execute(gomock.Any(), gomock.Any(), allowedTenants).
		DoAndReturn(func(ctx context.Context, jobEntity *entities.Job, tenants []string) (*entities.Job, error) {
			if jobEntity.Type != jobTypes[jobIdx] {
				return nil, fmt.Errorf("invalid job type. Got %s, expected %s", jobEntity.Type, jobTypes[jobIdx])
			}

			jobEntity.Transaction.From = txRequest.Params.From
			jobEntity.UUID = jobUUID
			jobIdx = jobIdx + 1
			return jobEntity, nil
		}).Times(len(jobTypes))

	// We flag this "special" scenario as faucet funding tx flow
	if withFaucet {
		faucet := testutils3.FakeFaucet()
		s.GetFaucetCandidate.EXPECT().Execute(ctx, txRequest.Params.From, chains[0], tenants).Return(faucet, nil)

		expectedFaucetJob := &entities.Job{
			ScheduleUUID: txRequest.Schedule.UUID,
			ChainUUID:    chains[0].UUID,
			Type:         utils.EthereumTransaction,
			Labels: map[string]string{
				"faucetUUID": faucet.UUID,
			},
			InternalData: &entities.InternalData{},
			Transaction: &entities.ETHTransaction{
				From:  faucet.CreditorAccount,
				To:    txRequest.Params.From,
				Value: faucet.Amount,
			},
		}
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), expectedFaucetJob, tenants).
			DoAndReturn(func(ctx context.Context, jobEntity *entities.Job, tenants []string) (*entities.Job, error) {
				if jobEntity.Transaction.From != faucet.CreditorAccount {
					return nil, fmt.Errorf("invalid from account. Got %s, expected %s", jobEntity.Transaction.From, faucet.CreditorAccount)
				}

				jobEntity.UUID = faucet.UUID
				return jobEntity, nil
			})
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, allowedTenants).Return(nil)
	} else {
		s.GetFaucetCandidate.EXPECT().Execute(ctx, txRequest.Params.From, chains[0], allowedTenants).Return(nil, faucetNotFoundErr)
	}

	s.StartJobUC.EXPECT().Execute(ctx, jobUUID, allowedTenants).Return(nil)
	s.GetTxUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(txRequest, nil)

	return s.usecase.Execute(ctx, txRequest, txData, tenantID)
}
