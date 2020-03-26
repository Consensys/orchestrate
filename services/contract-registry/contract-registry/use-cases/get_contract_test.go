package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/testutils"
)

func TestGetContract_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	contract := testutils.FakeContract()
	mockArtifactDataAgent := mocks.NewMockArtifactDataAgent(ctrl)
	usecase := NewGetContract(mockArtifactDataAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		mockArtifactDataAgent.EXPECT().FindOneByNameAndTag(context.Background(), contract.GetName(), contract.GetTag()).Return(&models.ArtifactModel{
			ID:               1,
			Abi:              contract.GetAbi(),
			Bytecode:         contract.GetBytecode(),
			DeployedBytecode: contract.GetDeployedBytecode(),
			Codehash:         "",
		}, nil)

		response, err := usecase.Execute(context.Background(), contract.GetId())

		assert.Equal(t, response, contract)
		assert.Nil(t, err)
	})

	t.Run("should fail if fails to extract contract name and tag", func(t *testing.T) {
		response, err := usecase.Execute(context.Background(), nil)

		assert.Nil(t, response)
		assert.NotNil(t, err)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		mockArtifactDataAgent.EXPECT().FindOneByNameAndTag(context.Background(), contract.GetName(), contract.GetTag()).Return(nil, dataAgentError)

		response, err := usecase.Execute(context.Background(), contract.GetId())

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(getContractComponent), err)
	})
}
