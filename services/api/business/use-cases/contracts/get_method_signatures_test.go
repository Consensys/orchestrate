// +build unit

package contracts

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	mocks2 "github.com/consensys/orchestrate/services/api/business/use-cases/mocks"
)

func TestGetMethodSignatures_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockGetContractUC := mocks2.NewMockGetContractUseCase(ctrl)

	usecase := NewGetMethodSignaturesUseCase(mockGetContractUC)

	t.Run("should execute use case successfully", func(t *testing.T) {
		contract := testutils.FakeContract()

		mockGetContractUC.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Name, contract.Tag, "transfer")

		assert.NoError(t, err)
		assert.Equal(t, signatures[0], "transfer(address,uint256)")
	})

	t.Run("should execute use case successfully if method name is constructor", func(t *testing.T) {
		contract := testutils.FakeContract()

		mockGetContractUC.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Name, contract.Tag, constructorMethodName)

		assert.NoError(t, err)
		assert.Equal(t, signatures[0], "constructor")
	})

	t.Run("should execute use case successfully and return an empty array if nothing is found", func(t *testing.T) {
		contract := testutils.FakeContract()

		mockGetContractUC.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Name, contract.Tag, "inexistentMethod")

		assert.NoError(t, err)
		assert.Empty(t, signatures)
	})

	t.Run("should fail with same error if get contract fails", func(t *testing.T) {
		contract := testutils.FakeContract()
		expectedErr := fmt.Errorf("error")

		mockGetContractUC.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(nil, expectedErr)

		signatures, err := usecase.Execute(ctx, contract.Name, contract.Tag, constructorMethodName)

		assert.Nil(t, signatures)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getMethodSignaturesComponent), err)
	})

	t.Run("should fail with DataCorruptedError if fails to get the ABI", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.ABI = "wrongABI"

		mockGetContractUC.EXPECT().Execute(gomock.Any(), contract.Name, contract.Tag).Return(contract, nil)

		signatures, err := usecase.Execute(ctx, contract.Name, contract.Tag, constructorMethodName)

		assert.Nil(t, signatures)
		assert.True(t, errors.IsDataCorruptedError(err))
	})
}
