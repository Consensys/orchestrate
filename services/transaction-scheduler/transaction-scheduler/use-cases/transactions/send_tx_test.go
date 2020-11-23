// +build unit

package transactions

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client/mock"
	models2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/models/testutils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/use-cases/mocks"
	mocks6 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/use-cases/mocks"
	mocks4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/use-cases/mocks"
	mocks5 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/use-cases/mocks"
	mocks3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/transaction-scheduler/validators/mocks"
)

type sendTxSuite struct {
	suite.Suite
	usecase             usecases.SendTxUseCase
	DB                  *mocks2.MockDB
	DBTX                *mocks2.MockTx
	Validators          *mocks3.MockTransactionValidator
	ChainRegistryClient *mock.MockChainRegistryClient
	TxRequestDA         *mocks2.MockTransactionRequestAgent
	ScheduleDA          *mocks2.MockScheduleAgent
	StartJobUC          *mocks.MockStartJobUseCase
	CreateJobUC         *mocks6.MockCreateJobUseCase
	CreateScheduleUC    *mocks4.MockCreateScheduleUseCase
	GetTxUC             *mocks5.MockGetTxUseCase
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
	s.Validators = mocks3.NewMockTransactionValidator(ctrl)
	s.TxRequestDA = mocks2.NewMockTransactionRequestAgent(ctrl)
	s.ScheduleDA = mocks2.NewMockScheduleAgent(ctrl)
	s.StartJobUC = mocks.NewMockStartJobUseCase(ctrl)
	s.CreateJobUC = mocks6.NewMockCreateJobUseCase(ctrl)
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
		chain := &models2.Chain{UUID: "32da9731-0fb8-4235-a4cc-35070ffe5bf0"}
		jobUUID := txRequest.Schedule.Jobs[0].UUID
		txData := ""
		txRequestModel := testutils2.FakeTxRequest(0)
		txRequestModel.RequestHash = "ea2d3e36db863014fdfaf49a88c31f1d"

		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(txRequestModel, nil)
		s.GetTxUC.EXPECT().Execute(ctx, txRequestModel.Schedule.UUID, tenants).Return(txRequest, nil)
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(ctx, gomock.Any(), chain.UUID).Return(nil, faucetNotFoundErr)
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
		txRequest.Schedule.Jobs[0].Logs = append(txRequest.Schedule.Jobs[0].Logs, &entities.Log{
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
	chain := &models2.Chain{UUID: uuid.Must(uuid.NewV4()).String()}
	jobUUID := uuid.Must(uuid.NewV4()).String()
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	txData := ""

	s.T().Run("should fail with same error if chain registry client fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(nil, expectedErr)
	
		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	
	s.T().Run("should fail with InvalidParameterError if chain registry client fails with NotFoundError", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(nil, expectedErr)
	
		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
		assert.Nil(t, response)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
	
	s.T().Run("should fail with same error if FindOne fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), common.HexToAddress(txRequest.Params.From), chain.UUID).Return(nil, faucetNotFoundErr)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], expectedErr)
	
		response, err := s.usecase.Execute(ctx, txRequest, txData, tenantID)
	
		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(sendTxComponent), err)
	})
	
	s.T().Run("should fail with same error if getFaucetCandidate request fails", func(t *testing.T) {
		expectedErr := errors.PostgresConnectionError("error")
		chain := &models2.Chain{UUID: uuid.Must(uuid.NewV4()).String()}
		txRequest := testutils3.FakeTxRequest()
		txRequest.Schedule.UUID = scheduleUUID
		txRequest.Schedule.Jobs[0].UUID = jobUUID
		txData = ""
		scheduleModel := testutils2.FakeSchedule(tenants[0])
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), common.HexToAddress(txRequest.Params.From), chain.UUID).Return(nil, expectedErr)
	
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), common.HexToAddress(txRequest.Params.From), chain.UUID).
			Return(nil, faucetNotFoundErr)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(expectedErr)
	
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
	
		s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
		s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
		s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
		s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
		s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(gomock.Any(), common.HexToAddress(txRequest.Params.From), chain.UUID).
			Return(nil, faucetNotFoundErr)
		s.CreateJobUC.EXPECT().Execute(gomock.Any(), gomock.Any(), tenants).Return(txRequest.Schedule.Jobs[0], nil)
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
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
	chain := &models2.Chain{UUID: "32da9731-0fb8-4235-a4cc-35070ffe5bf0"}
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	txData := ""
	scheduleModel := testutils2.FakeSchedule(tenants[0])
	jobIdx := 0

	s.ChainRegistryClient.EXPECT().GetChainByName(ctx, txRequest.ChainName).Return(chain, nil)
	s.TxRequestDA.EXPECT().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID).Return(nil, errors.NotFoundError(""))
	s.CreateScheduleUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(txRequest.Schedule, nil)
	s.ScheduleDA.EXPECT().FindOneByUUID(ctx, txRequest.Schedule.UUID, tenants).Return(scheduleModel, nil)
	s.TxRequestDA.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
	s.CreateJobUC.EXPECT().
		Execute(gomock.Any(), gomock.Any(), tenants).
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
		creditor := common.HexToAddress("0xacd")
		fctJobUUID := "fctJobUUID"
		fct := &chainregistry.Faucet{
			Creditor: creditor,
			Amount:   big.NewInt(10),
		}
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(multitenancy.WithTenantID(ctx, tenantID), common.HexToAddress(txRequest.Params.From), chain.UUID).
			Return(fct, nil)

		s.CreateJobUC.EXPECT().Execute(gomock.Any(), generateFaucetJob(fct, txRequest.Schedule.UUID, chain.UUID, txRequest.Params.From), tenants).
			DoAndReturn(func(ctx context.Context, jobEntity *entities.Job, tenants []string) (*entities.Job, error) {
				if jobEntity.Transaction.From != creditor.Hex() {
					return nil, fmt.Errorf("invalid from account. Got %s, expected %s", jobEntity.Transaction.From, creditor.Hex())
				}

				jobEntity.UUID = fctJobUUID
				return jobEntity, nil
			})
		s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
	} else {
		s.ChainRegistryClient.EXPECT().GetFaucetCandidate(multitenancy.WithTenantID(ctx, tenantID), common.HexToAddress(txRequest.Params.From), chain.UUID).
			Return(nil, faucetNotFoundErr)
	}

	s.StartJobUC.EXPECT().Execute(ctx, jobUUID, tenants).Return(nil)
	s.GetTxUC.EXPECT().Execute(ctx, txRequest.Schedule.UUID, tenants).Return(txRequest, nil)

	return s.usecase.Execute(ctx, txRequest, txData, tenantID)
}
