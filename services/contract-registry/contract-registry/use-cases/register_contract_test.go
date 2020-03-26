package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/testutils"
)

func TestRegisterContract_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContractDA := mocks.NewMockContractDataAgent(ctrl)

	usecase := NewRegisterContract(mockContractDA)

	t.Run("should execute use case successfully", func(t *testing.T) {
		contract := testutils.FakeContract()

		abiCompacted, _ := contract.GetABICompacted()
		mockContractDA.EXPECT().Insert(
			context.Background(),
			contract.GetName(),
			contract.GetTag(),
			abiCompacted,
			contract.GetBytecode(),
			contract.GetDeployedBytecode(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any()).Return(nil)

		err := usecase.Execute(context.Background(), contract)

		assert.Nil(t, err)
	})

	t.Run("should fail if it fails to extract artifacts", func(t *testing.T) {
		err := usecase.Execute(context.Background(), nil)
		assert.NotNil(t, err)
	})

	t.Run("should fail if it fails to name and tag", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.Id = nil
		err := usecase.Execute(context.Background(), contract)
		assert.NotNil(t, err)
	})

	t.Run("should fail with InvalidArg error if it fails to decode bytecode", func(t *testing.T) {
		contract := testutils.FakeContract()
		contract.DeployedBytecode = "hello!"
		err := usecase.Execute(context.Background(), contract)
		assert.True(t, errors.IsDataError(err))
	})

	t.Run("should fail if Insert fails", func(t *testing.T) {
		contract := testutils.FakeContract()
		abiCompacted, _ := contract.GetABICompacted()
		mockContractDA.EXPECT().Insert(
			context.Background(),
			contract.GetName(),
			contract.GetTag(),
			abiCompacted,
			contract.GetBytecode(),
			contract.GetDeployedBytecode(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any()).Return(fmt.Errorf("error"))

		err := usecase.Execute(context.Background(), contract)
		assert.NotNil(t, err)
	})
}
