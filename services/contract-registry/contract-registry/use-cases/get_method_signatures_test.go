// +build unit

package usecases

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/contract-registry/use-cases/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetMethodSignatures_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetContractUC := mocks.NewMockGetContractUseCase(ctrl)
	ctx := context.Background()

	usecase := NewGetMethodSignatures(mockGetContractUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		contract := testutils.FakeContract()

		mockGetContractUC.EXPECT().Execute(ctx, contract.Id).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Id, "transfer")

		assert.NoError(t, err)
		assert.Equal(t, signatures[0], "transfer(address,uint256)")
	})

	t.Run("should execute use case successfully if method name is constructor", func(t *testing.T) {
		contract := testutils.FakeContract()

		mockGetContractUC.EXPECT().Execute(ctx, contract.Id).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Id, constructorMethodName)

		assert.NoError(t, err)
		assert.Equal(t, signatures[0], "constructor(uint256)")
	})

	t.Run("should execute use case successfully and return an empty array if nothing is found", func(t *testing.T) {
		contract := testutils.FakeContract()

		mockGetContractUC.EXPECT().Execute(ctx, contract.Id).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Id, "inexistentMethod")

		assert.NoError(t, err)
		assert.Empty(t, signatures)
	})

	t.Run("should fail with same error if get contract fails", func(t *testing.T) {
		contract := testutils.FakeContract()
		expectedErr := fmt.Errorf("error")

		mockGetContractUC.EXPECT().Execute(ctx, contract.Id).Return(nil, expectedErr)

		signatures, err := usecase.Execute(ctx, contract.Id, constructorMethodName)

		assert.Nil(t, signatures)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getMethodSignaturesComponent), err)
	})

	t.Run("should fail with DataCorruptedError if fails to get the ABI", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.Abi = "wrongABI"

		mockGetContractUC.EXPECT().Execute(ctx, contract.Id).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Id, constructorMethodName)

		assert.Nil(t, signatures)
		assert.True(t, errors.IsDataCorruptedError(err))
	})
}
