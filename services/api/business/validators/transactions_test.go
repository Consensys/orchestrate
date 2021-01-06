// +build unit

package validators

import (
	"context"
	"testing"

	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client/mock"

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
	validator               TransactionValidator
	mockTxRequestDA         *mocks.MockTransactionRequestAgent
	mockChainRegistryClient *mock.MockChainRegistryClient
	mockGetContractUC       *mocks2.MockGetContractUseCase
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
	s.mockGetContractUC = mocks2.NewMockGetContractUseCase(ctrl)

	s.validator = NewTransactionValidator(s.mockChainRegistryClient)
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
			Return(nil, errors.NotFoundError("error"))
		_, err := s.validator.ValidateChainExists(context.Background(), chainUUID)

		assert.Error(t, err)
		assert.True(t, errors.IsInvalidParameterError(err))
	})

	s.T().Run("should fail with same error if ChainRegistryClient fails", func(t *testing.T) {
		expectedErr := errors.ServiceConnectionError("error")
		s.mockChainRegistryClient.EXPECT().GetChainByUUID(gomock.Any(), chainUUID).
			Return(nil, expectedErr)
		_, err := s.validator.ValidateChainExists(context.Background(), chainUUID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.FromError(err).ExtendComponent(validatorComponent))
	})
}
