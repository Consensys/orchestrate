package validators

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	models2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
)

type transactionsTestSuite struct {
	suite.Suite
	validator               TransactionValidator
	mockTxRequestDA         *mocks2.MockTransactionRequestAgent
	mockChainRegistryClient *mock.MockChainRegistryClient
}

func TestTransactionValidator(t *testing.T) {
	s := new(transactionsTestSuite)
	suite.Run(t, s)
}

func (s *transactionsTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	mockDB := mocks2.NewMockDB(ctrl)
	s.mockTxRequestDA = mocks2.NewMockTransactionRequestAgent(ctrl)
	mockDB.EXPECT().TransactionRequest().Return(s.mockTxRequestDA).AnyTimes()

	s.mockChainRegistryClient = mock.NewMockChainRegistryClient(ctrl)
	s.validator = NewTransactionValidator(mockDB, s.mockChainRegistryClient)
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateRequestHash() {
	txRequest := testutils.FakeTransactionRequest()

	s.T().Run("should validate tx successfully and return the request hash if data agent returns not found", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(nil, errors.NotFoundError("error"))
		expectedRequestHash := "2c9e0e4c6834e516fb99b0bdb00d4973"

		requestHash, err := s.validator.ValidateRequestHash(context.Background(), txRequest.Params, txRequest.IdempotencyKey)

		assert.Nil(t, err)
		assert.Equal(t, expectedRequestHash, requestHash)
	})

	s.T().Run("should return error if data agent fails", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(nil, errors.PostgresConnectionError("error"))

		requestHash, err := s.validator.ValidateRequestHash(context.Background(), txRequest.Params, txRequest.IdempotencyKey)

		assert.Empty(t, requestHash)
		assert.Equal(t, errors.PostgresConnectionError("error").ExtendComponent("transaction-validator"), err)
	})

	s.T().Run("should return AlreadyExistsError if data agent returns a request a different request hash", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(&models.TransactionRequest{
			IdempotencyKey: txRequest.IdempotencyKey,
			RequestHash:    "differentRequestHash",
		}, nil)

		requestHash, err := s.validator.ValidateRequestHash(context.Background(), txRequest.Params, txRequest.IdempotencyKey)

		assert.Empty(t, requestHash)
		assert.True(t, errors.IsAlreadyExistsError(err))
	})
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateChainExists() {
	chainUUID := "chainUUID"

	s.T().Run("should validate chain successfully", func(t *testing.T) {
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).Return(&models2.Chain{}, nil)
		err := s.validator.ValidateChainExists(context.Background(), chainUUID)
		assert.Nil(t, err)
	})

	s.T().Run("should fail with InvalidParameterError if ChainRegistryClient fails", func(t *testing.T) {
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).Return(nil, fmt.Errorf("error"))
		err := s.validator.ValidateChainExists(context.Background(), chainUUID)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
