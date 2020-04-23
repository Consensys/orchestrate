package validators

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types/testutils"
)

type transactionsTestSuite struct {
	suite.Suite
	validator               TransactionValidator
	mockChainRegistryClient *mock.MockChainRegistryClient
}

func TestTransactionValidator(t *testing.T) {
	s := new(transactionsTestSuite)
	suite.Run(t, s)
}

func (s *transactionsTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockChainRegistryClient = mock.NewMockChainRegistryClient(ctrl)
	s.validator = NewTransaction(s.mockChainRegistryClient)
}

func (s *transactionsTestSuite) TestTransactionValidator_ValidateTx() {
	s.T().Run("should validate tx successfully", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), txRequest.ChainID).Return(nil, nil)

		err := s.validator.ValidateTx(context.Background(), txRequest)

		assert.Nil(t, err)
	})

	s.T().Run("should fail with client error if client fails to fetch chain", func(t *testing.T) {
		txRequest := testutils.FakeTransactionRequest()
		expectedError := fmt.Errorf("error")
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), txRequest.ChainID).Return(nil, expectedError)

		err := s.validator.ValidateTx(context.Background(), txRequest)

		assert.Equal(t, "transaction-validator", errors.FromError(err).Component)
		assert.Equal(t, expectedError.Error(), errors.FromError(err).Message)
	})
}
