package validators

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
)

type transactionsTestSuite struct {
	suite.Suite
	validator       TransactionValidator
	mockTxRequestDA *mocks.MockTransactionRequestAgent
}

func TestTransactionValidator(t *testing.T) {
	s := new(transactionsTestSuite)
	suite.Run(t, s)
}

func (s *transactionsTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockTxRequestDA = mocks.NewMockTransactionRequestAgent(ctrl)
	s.validator = NewTransactionValidator(s.mockTxRequestDA)
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateRequestHash() {
	txRequest := testutils.FakeTransactionRequest()

	s.T().Run("should validate tx successfully and return the request hash if data agent returns not found", func(t *testing.T) {
		s.mockTxRequestDA.EXPECT().FindOneByIdempotencyKey(gomock.Any(), txRequest.IdempotencyKey).Return(nil, errors.NotFoundError("error"))
		expectedRequestHash := "ec5269a2be64b487bc1741dead29bc40"

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
