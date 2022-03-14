// +build unit

package transactions

import (
	"context"
	"fmt"
	"testing"

	testutils3 "github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
	mocks2 "github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models"
	testutils2 "github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	GetTxUC            *mocks.MockGetTxUseCase
	GetFaucetCandidate *mocks.MockGetFaucetCandidateUseCase
	userInfo           *multitenancy.UserInfo
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
	s.GetTxUC = mocks.NewMockGetTxUseCase(ctrl)
	s.GetFaucetCandidate = mocks.NewMockGetFaucetCandidateUseCase(ctrl)

	s.DB.EXPECT().Begin().Return(s.DBTX, nil).AnyTimes()
	s.DB.EXPECT().TransactionRequest().Return(s.TxRequestDA).AnyTimes()
	s.DB.EXPECT().Schedule().Return(s.ScheduleDA).AnyTimes()
	s.DBTX.EXPECT().Schedule().Return(s.ScheduleDA).AnyTimes()
	s.DBTX.EXPECT().Commit().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Rollback().Return(nil).AnyTimes()
	s.DBTX.EXPECT().Close().Return(nil).AnyTimes()
	s.DBTX.EXPECT().TransactionRequest().Return(s.TxRequestDA).AnyTimes()
	s.CreateJobUC.EXPECT().WithDBTransaction(s.DBTX).Return(s.CreateJobUC).AnyTimes()
	s.userInfo = multitenancy.NewUserInfo("tenantOne", "username")

	s.usecase = NewSendTxUseCase(
		s.DB,
		s.SearchChainsUC,
		s.StartJobUC,
		s.CreateJobUC,
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

		response, err := successfulTestExecution(s, txRequest, false, entities.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a public tx with faucet", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, true, entities.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a EEA tx", func(t *testing.T) {
		txRequest := testutils3.FakeEEATxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = entities.EEAChainType

		response, err := successfulTestExecution(s, txRequest, false, entities.EEAPrivateTransaction,
			entities.EEAMarkingTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a tessera tx", func(t *testing.T) {
		txRequest := testutils3.FakeTesseraTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.Protocol = entities.TesseraChainType

		response, err := successfulTestExecution(s, txRequest, false, entities.TesseraPrivateTransaction,
			entities.TesseraMarkingTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a raw tx", func(t *testing.T) {
		txRequest := testutils3.FakeRawTxRequest()
		txRequest.Params.Raw = hexutil.MustDecode("0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e")
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		response, err := successfulTestExecution(s, txRequest, false, entities.EthereumRawTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should not insert and start job in DB if TxRequest already exists and send if status is CREATED", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		ctx := context.Background()
		chains := []*entities.Chain{testutils3.FakeChain()}
		chains[0].UUID = "myChainUUID"
		jobUUID := txRequest.Schedule.Jobs[0].UUID
		txData := (hexutil.Bytes)(hexutil.MustDecode("0x"))
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "64ab842d4e824ac64e0cb5585164db7f"

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(gomock.Any(), txRequestModel.Schedule.UUID, s.userInfo).Return(txRequest, nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), gomock.Any(), chains[0], s.userInfo).
			Return(nil, faucetNotFoundErr)
		s.StartJobUC.EXPECT().Execute(gomock.Any(), jobUUID, s.userInfo).Return(nil)
		s.GetTxUC.EXPECT().Execute(gomock.Any(), txRequest.Schedule.UUID, s.userInfo).Return(txRequest, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)

		require.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should not insert and not start job if TxRequest already exists and not send if status is not CREATED", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Schedule.Jobs[0].Status = entities.StatusStarted

		ctx := context.Background()
		txData := (hexutil.Bytes)(hexutil.MustDecode("0x"))
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "64ab842d4e824ac64e0cb5585164db7f"
		chains := []*entities.Chain{testutils3.FakeChain()}
		chains[0].UUID = "myChainUUID"

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(gomock.Any(), txRequestModel.Schedule.UUID, s.userInfo).Return(txRequest, nil)
		s.GetTxUC.EXPECT().Execute(gomock.Any(), txRequest.Schedule.UUID, s.userInfo).Return(txRequest, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
		require.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
	})

	s.T().Run("should execute send successfully a oneTimeKey tx", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txRequest.Params.From = nil
		txRequest.InternalData.OneTimeKey = true

		response, err := successfulTestExecution(s, txRequest, false, entities.EthereumTransaction)
		assert.NoError(t, err)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.Equal(t, txRequest.IdempotencyKey, response.IdempotencyKey)
		assert.Equal(t, txRequest.Schedule.UUID, response.Schedule.UUID)
		assert.True(t, response.Schedule.Jobs[0].InternalData.OneTimeKey)
	})
}

func (s *sendTxSuite) TestSendTx_ExpectedErrors() {
	ctx := context.Background()

	chains := []*entities.Chain{testutils3.FakeChain()}
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	txData := (hexutil.Bytes)(hexutil.MustDecode("0x"))

	s.T().Run("should fail with same error if chain agent fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).
			Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with InvalidParameterError if no chain is found", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).
			Return([]*entities.Chain{}, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
		assert.Nil(t, response)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with AlreadyExistsError if request found has different request hash", func(t *testing.T) {
		expectedErr := errors.AlreadyExistsError("transaction request with the same idempotency key and different params already exists")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(&models.TransactionRequest{
			RequestHash: "differentRequestHash",
		}, nil)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createSchedule UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).
			Return(nil, errors.NotFoundError(""))
		s.ScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if select or insert txRequest fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(s.userInfo.TenantID, s.userInfo.Username)

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(nil, errors.NotFoundError(""))
		s.ScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), txRequest.Schedule.UUID, s.userInfo.AllowedTenants, s.userInfo.Username).
			Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if createJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData := (hexutil.Bytes)(hexutil.MustDecode("0x"))
		scheduleModel := testutils2.FakeSchedule(s.userInfo.TenantID, s.userInfo.Username)

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}}, s.userInfo).
			Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(nil, errors.NotFoundError(""))
		s.ScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), txRequest.Schedule.UUID, s.userInfo.AllowedTenants, s.userInfo.Username).
			Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), *txRequest.Schedule.Jobs[0].Transaction.From, chains[0], s.userInfo).
			Return(nil, faucetNotFoundErr)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).
			Return(txRequest.Schedule.Jobs[0], expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if getFaucetCandidate request fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		chains := []*entities.Chain{testutils3.FakeChain()}
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData := (hexutil.Bytes)(hexutil.MustDecode("0x"))
		scheduleModel := testutils2.FakeSchedule(s.userInfo.TenantID, s.userInfo.Username)

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}}, s.userInfo).
			Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(nil, errors.NotFoundError(""))
		s.ScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), txRequest.Schedule.UUID, s.userInfo.AllowedTenants, s.userInfo.Username).
			Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).Return(txRequest.Schedule.Jobs[0], nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), *txRequest.Schedule.Jobs[0].Transaction.From, gomock.Any(), s.userInfo).
			Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if startJob UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		scheduleModel := testutils2.FakeSchedule(s.userInfo.TenantID, s.userInfo.Username)

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(nil, errors.NotFoundError(""))
		s.ScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), txRequest.Schedule.UUID, s.userInfo.AllowedTenants, s.userInfo.Username).
			Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), s.userInfo).Return(txRequest.Schedule.Jobs[0], nil)
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), *txRequest.Schedule.Jobs[0].Transaction.From, gomock.Any(), s.userInfo).Return(nil, faucetNotFoundErr)
		s.StartJobUC.EXPECT().Execute(gomock.Any(), jobUUID, s.userInfo).Return(expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})

	s.T().Run("should fail with same error if getTx UseCase fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID

		s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
			s.userInfo).Return(chains, nil)

		requestHash, _ := generateRequestHash(chains[0].UUID, txRequest.Params)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
			s.userInfo.Username).Return(&models.TransactionRequest{
			RequestHash: requestHash,
			Schedule:    parsers.NewScheduleModelFromEntities(txRequest.Schedule),
		}, nil)
		s.GetTxUC.EXPECT().Execute(gomock.Any(), txRequest.Schedule.UUID, s.userInfo).Return(nil, expectedErr)

		response, err := s.usecase.Execute(ctx, txRequest, txData, s.userInfo)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
}

func successfulTestExecution(s *sendTxSuite, txRequest *entities.TxRequest, withFaucet bool, jobTypes ...entities.JobType) (*entities.TxRequest, error) {
	ctx := context.Background()
	chains := []*entities.Chain{testutils3.FakeChain()}
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	txData := (hexutil.Bytes)(hexutil.MustDecode("0x"))
	scheduleModel := testutils2.FakeSchedule(s.userInfo.TenantID, s.userInfo.Username)
	jobIdx := 0
	from := new(ethcommon.Address)
	if txRequest.Params.From != nil {
		from = txRequest.Params.From
	}

	s.SearchChainsUC.EXPECT().Execute(gomock.Any(), &entities.ChainFilters{Names: []string{txRequest.ChainName}},
		s.userInfo).Return(chains, nil)
	s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey, s.userInfo.TenantID,
		s.userInfo.Username).Return(nil, errors.NotFoundError(""))
	s.ScheduleDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	s.ScheduleDA.EXPECT().FindOneByUUID(gomock.Any(), txRequest.Schedule.UUID, s.userInfo.AllowedTenants, s.userInfo.Username).
		Return(scheduleModel, nil)
	s.TxRequestDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	s.CreateJobUC.EXPECT().
		Execute(gomock.Any(), gomock.Any(), s.userInfo).
		DoAndReturn(func(ctx context.Context, jobEntity *entities.Job, userInfo *multitenancy.UserInfo) (*entities.Job, error) {
			if jobEntity.Type != jobTypes[jobIdx] {
				return nil, fmt.Errorf("invalid job type. Got %s, expected %s", jobEntity.Type, jobTypes[jobIdx])
			}

			jobEntity.Transaction.From = from
			jobEntity.UUID = jobUUID
			jobEntity.Status = entities.StatusCreated
			jobIdx = jobIdx + 1
			return jobEntity, nil
		}).Times(len(jobTypes))

	// We flag this "special" scenario as faucet funding tx flow
	if withFaucet {
		internalAdminUser := multitenancy.NewInternalAdminUser()
		internalAdminUser.TenantID = s.userInfo.TenantID
		faucet := testutils3.FakeFaucet()
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), *from, chains[0], s.userInfo).Return(faucet, nil)

		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), internalAdminUser).
			DoAndReturn(func(ctx context.Context, jobEntity *entities.Job, userInfo *multitenancy.UserInfo) (*entities.Job, error) {
				if jobEntity.Transaction.From.String() != faucet.CreditorAccount.String() {
					return nil, fmt.Errorf("invalid from account. Got %s, expected %s", jobEntity.Transaction.From, faucet.CreditorAccount)
				}

				jobEntity.UUID = faucet.UUID
				return jobEntity, nil
			})
		s.StartJobUC.EXPECT().Execute(gomock.Any(), faucet.UUID, internalAdminUser).Return(nil)
	} else {
		s.GetFaucetCandidate.EXPECT().Execute(gomock.Any(), *from, chains[0], s.userInfo).Return(nil, faucetNotFoundErr)
	}

	s.StartJobUC.EXPECT().Execute(gomock.Any(), jobUUID, s.userInfo).Return(nil)

	return s.usecase.Execute(ctx, txRequest, txData, s.userInfo)
}
