package validators

import (
	"context"
	"fmt"
	"testing"

	abi2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client/mock"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	testutils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/testutils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	chainmodel "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

const (
	chainUUID = "chainUUID"
)

type transactionsTestSuite struct {
	suite.Suite
	validator                  TransactionValidator
	mockTxRequestDA            *mocks.MockTransactionRequestAgent
	mockChainRegistryClient    *mock.MockChainRegistryClient
	mockContractRegistryClient *mock2.MockContractRegistryClient
}

func TestTransactionValidator(t *testing.T) {
	s := new(transactionsTestSuite)
	suite.Run(t, s)
}

func (s *transactionsTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	s.mockTxRequestDA = mocks.NewMockTransactionRequestAgent(ctrl)
	s.mockChainRegistryClient = mock.NewMockChainRegistryClient(ctrl)
	s.mockContractRegistryClient = mock2.NewMockContractRegistryClient(ctrl)

	mockDB.EXPECT().TransactionRequest().Return(s.mockTxRequestDA).AnyTimes()

	s.validator = NewTransactionValidator(mockDB, s.mockChainRegistryClient, s.mockContractRegistryClient)
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateRequestHash() {
	txRequest := testutils.FakeSendTransactionRequest("chainName")

	s.T().Run("should validate tx successfully and return the request hash if data agent returns not found", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(nil, errors.NotFoundError("error"))
		expectedRequestHash := "d8055d499fcbb64b67a5eab35d5c8109"

		requestHash, err := s.validator.ValidateRequestHash(context.Background(), chainUUID, txRequest.Params, txRequest.IdempotencyKey)

		assert.Nil(t, err)
		assert.Equal(t, expectedRequestHash, requestHash)
	})

	s.T().Run("should return error if data agent fails", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(nil, errors.PostgresConnectionError("error"))

		requestHash, err := s.validator.ValidateRequestHash(context.Background(), chainUUID, txRequest.Params, txRequest.IdempotencyKey)

		assert.Empty(t, requestHash)
		assert.Equal(t, errors.PostgresConnectionError("error").ExtendComponent("transaction-validator"), err)
	})

	s.T().Run("should return AlreadyExistsError if data agent returns a request a different request hash", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(&models.TransactionRequest{
			IdempotencyKey: txRequest.IdempotencyKey,
			RequestHash:    "differentRequestHash",
		}, nil)

		requestHash, err := s.validator.ValidateRequestHash(context.Background(), chainUUID, txRequest.Params, txRequest.IdempotencyKey)

		assert.Empty(t, requestHash)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateChainExists() {
	s.T().Run("should validate chain successfully", func(t *testing.T) {
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).Return(&chainmodel.Chain{}, nil)
		err := s.validator.ValidateChainExists(context.Background(), chainUUID)
		assert.Nil(t, err)
	})

	s.T().Run("should fail with InvalidParameterError if ChainRegistryClient fails", func(t *testing.T) {
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).Return(nil, fmt.Errorf("error"))
		err := s.validator.ValidateChainExists(context.Background(), chainUUID)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateMethodSignature() {
	s.T().Run("should validate method signature successfully", func(t *testing.T) {
		txData, err := s.validator.ValidateMethodSignature("method(string,uint256)", []string{"val1", "15"})
		assert.Nil(t, err)
		assert.NotEmpty(t, txData)
	})

	s.T().Run("should fail with InvalidParameterError if ChainRegistryClient fails", func(t *testing.T) {
		_, err := s.validator.ValidateMethodSignature("method(string,uint256)", []string{"val1"})
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateContract() {
	ctx := context.Background()

	s.T().Run("should validate contract successfully", func(t *testing.T) {
		txRequest := testutils2.FakeTxRequestEntity()
		txRequest.Params.Args = []string{"300"}
		contract := testutils3.FakeContract()

		s.mockContractRegistryClient.EXPECT().GetContract(ctx, &contractregistry.GetContractRequest{
			ContractId: &abi2.ContractId{
				Name: txRequest.Params.ContractName,
				Tag:  txRequest.Params.ContractTag,
			},
		}).Return(&contractregistry.GetContractResponse{
			Contract: contract,
		}, nil)

		txData, err := s.validator.ValidateContract(ctx, txRequest.Params)
		assert.Nil(t, err)
		assert.NotEmpty(t, txData)
	})

	s.T().Run("should fail with InvalidParameterError if ContractRegistryClient fails", func(t *testing.T) {
		txRequest := testutils2.FakeTxRequestEntity()
		txRequest.Params.Args = []string{"300"}
		expectedErr := fmt.Errorf("error")

		s.mockContractRegistryClient.EXPECT().GetContract(ctx, gomock.Any()).Return(nil, expectedErr)

		txData, err := s.validator.ValidateContract(ctx, txRequest.Params)

		assert.Empty(t, txData)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with DataCorruptedError if bytecode decoding fails", func(t *testing.T) {
		txRequest := testutils2.FakeTxRequestEntity()
		txRequest.Params.Args = []string{"300"}
		contract := testutils3.FakeContract()
		contract.Bytecode = "Invalid bytecode"

		s.mockContractRegistryClient.EXPECT().GetContract(ctx, gomock.Any()).Return(&contractregistry.GetContractResponse{
			Contract: contract,
		}, nil)

		txData, err := s.validator.ValidateContract(ctx, txRequest.Params)

		assert.Empty(t, txData)
		assert.True(t, errors.IsDataCorruptedError(err))
	})

	s.T().Run("should fail with InvalidParameterError if invalid args", func(t *testing.T) {
		txRequest := testutils2.FakeTxRequestEntity()
		txRequest.Params.Args = []string{"InvalidArg"}
		contract := testutils3.FakeContract()

		s.mockContractRegistryClient.EXPECT().GetContract(ctx, gomock.Any()).Return(&contractregistry.GetContractResponse{
			Contract: contract,
		}, nil)

		txData, err := s.validator.ValidateContract(ctx, txRequest.Params)

		assert.Empty(t, txData)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
