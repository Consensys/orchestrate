// +build unit

package contracts

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSearchContract_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	contract := testutils.FakeContract()
	address := testutils.FakeAddress()
	contractAgent := mocks.NewMockContractAgent(ctrl)
	usecase := NewSearchContractUseCase(contractAgent)

	t.Run("should execute use case by address successfully", func(t *testing.T) {
		contractAgent.EXPECT().
			FindOneByAddress(gomock.Any(), address.String()).
			Return(contract, nil)

		response, err := usecase.Execute(ctx, nil, address)

		assert.NoError(t, err)
		assert.Equal(t, contract.ABI, response.ABI)
	})

	t.Run("should execute use case by code_hash successfully", func(t *testing.T) {
		contractAgent.EXPECT().
			FindOneByCodeHash(gomock.Any(), contract.Bytecode.String()).
			Return(contract, nil)

		response, err := usecase.Execute(ctx, contract.Bytecode, nil)

		assert.NoError(t, err)
		assert.Equal(t, contract.ABI, response.ABI)
	})
}
