// +build unit

package validators

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"

	abi2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/abi"
	testutils3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client/mock"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/client/mock"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
	mock3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
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
	mockIdentityManagerClient  *mock3.MockIdentityManagerClient
}

func TestTransactionValidator(t *testing.T) {
	s := new(transactionsTestSuite)
	suite.Run(t, s)
}

func (s *transactionsTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockTxRequestDA = mocks.NewMockTransactionRequestAgent(ctrl)
	s.mockChainRegistryClient = mock.NewMockChainRegistryClient(ctrl)
	s.mockContractRegistryClient = mock2.NewMockContractRegistryClient(ctrl)
	s.mockIdentityManagerClient = mock3.NewMockIdentityManagerClient(ctrl)

	s.validator = NewTransactionValidator(s.mockChainRegistryClient, s.mockContractRegistryClient, s.mockIdentityManagerClient)
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateChainExists() {
	chainModel := &models.Chain{ChainID: "888"}

	s.T().Run("should validate chain successfully", func(t *testing.T) {
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).Return(chainModel, nil)
		chainID, err := s.validator.ValidateChainExists(context.Background(), chainUUID)

		assert.NoError(t, err)
		assert.Equal(t, chainID, chainModel.ChainID)
	})

	s.T().Run("should fail with InvalidParameterError if ChainRegistryClient returns NotFound", func(t *testing.T) {
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).
			Return(nil, errors.NotFoundError(("error")))
		_, err := s.validator.ValidateChainExists(context.Background(), chainUUID)

		assert.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
	
	s.T().Run("should fail with same error if ChainRegistryClient fails", func(t *testing.T) {
		expectedErr := errors.ServiceConnectionError(("error"))
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).
			Return(nil, expectedErr)
		_, err := s.validator.ValidateChainExists(context.Background(), chainUUID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.FromError(err).ExtendComponent(txValidatorComponent))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateMethodSignature() {
	s.T().Run("should validate method signature successfully", func(t *testing.T) {
		txData, err := s.validator.ValidateMethodSignature(
			"method(bool,bool,string,string,uint256,uint256)",
			testutils3.ParseIArray("true", "false", "val1", "false", "15", "0"))
		assert.NoError(t, err)
		assert.NotEmpty(t, txData)

		txData2, err := s.validator.ValidateMethodSignature(
			"method(bool,bool,string,string,uint256,uint256)",
			testutils3.ParseIArray(true, false, "val1", "false", 15, 0))
		assert.NoError(t, err)
		assert.NotEmpty(t, txData)

		assert.Equal(t, txData, txData2)
	})

	s.T().Run("should validate complex method signature successfully", func(t *testing.T) {
		txData, err := s.validator.ValidateMethodSignature(
			`constructor(
    string memory name,
    string memory symbol,
    uint256 granularity,
    uint256[] granularities,
    address[] memory controllers,
    address certificateSigner,
    bool certificateActivated,
    bool[] certificateActivates,
    bytes32[] memory defaultPartitions
  )`,
			testutils3.ParseIArray("WindToken",
				"WIN",
				"0x1",
				[]int{1, 2},
				[]string{
					"0xF112b55061C3a023E07f5347Bbb92F84F0FC529a",
				},
				"0xe31C41f0f70C5ff39f73B4B94bcCD767b3071630",
				false,
				[]bool{true, false},
				[]string{
					"0x7265736572766564000000000000000000000000000000000000000000000000",
					"0x6973737565640000000000000000000000000000000000000000000000000000",
					"0x6c6f636b65640000000000000000000000000000000000000000000000000000",
				}))
		assert.NoError(t, err)
		assert.NotEmpty(t, txData)
	})

	s.T().Run("should fail with InvalidParameterError if ChainRegistryClient fails", func(t *testing.T) {
		_, err := s.validator.ValidateMethodSignature("method(string,uint256)", testutils3.ParseIArray("val1"))
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateContract() {
	ctx := context.Background()

	s.T().Run("should validate contract successfully", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Params.Args = testutils3.ParseIArray("300")
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
		assert.NoError(t, err)
		assert.NotEmpty(t, txData)
	})

	s.T().Run("should fail with InvalidParameterError if ContractRegistryClient fails", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Params.Args = testutils3.ParseIArray("300")
		expectedErr := fmt.Errorf("error")

		s.mockContractRegistryClient.EXPECT().GetContract(ctx, gomock.Any()).Return(nil, expectedErr)

		txData, err := s.validator.ValidateContract(ctx, txRequest.Params)

		assert.Empty(t, txData)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with DataCorruptedError if bytecode decoding fails", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		txRequest.Params.Args = testutils3.ParseIArray("300")
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
		txRequest := testutils3.FakeTxRequest()
		txRequest.Params.Args = testutils3.ParseIArray("InvalidArg")
		contract := testutils3.FakeContract()

		s.mockContractRegistryClient.EXPECT().GetContract(ctx, gomock.Any()).Return(&contractregistry.GetContractResponse{
			Contract: contract,
		}, nil)

		txData, err := s.validator.ValidateContract(ctx, txRequest.Params)

		assert.Empty(t, txData)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateAccount() {
	ctx := context.Background()

	s.T().Run("should validate account successfully", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()

		s.mockIdentityManagerClient.EXPECT().GetAccount(ctx, txRequest.Params.From).
			Return(&identitymanager.AccountResponse{}, nil)

		err := s.validator.ValidateAccount(ctx, txRequest.Params.From)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with InvalidParameter error if IdentityManagerClient returns a NotFound error", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		expectedErr := errors.NotFoundError("not found")
		s.mockIdentityManagerClient.EXPECT().GetAccount(ctx, txRequest.Params.From).
			Return(&identitymanager.AccountResponse{}, expectedErr)

		err := s.validator.ValidateAccount(ctx, txRequest.Params.From)
		assert.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
	
	s.T().Run("should fail with same error if IdentityManagerClient returns an error", func(t *testing.T) {
		txRequest := testutils3.FakeTxRequest()
		expectedErr := errors.ConnectionError("not found")
		s.mockIdentityManagerClient.EXPECT().GetAccount(ctx, txRequest.Params.From).
			Return(&identitymanager.AccountResponse{}, expectedErr)

		err := s.validator.ValidateAccount(ctx, txRequest.Params.From)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.FromError(err).ExtendComponent(txValidatorComponent))
	})
}
