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
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models"
)

func TestGetContract_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	contract := testutils.FakeContract()
	artifactAgent := mocks.NewMockArtifactAgent(ctrl)
	usecase := NewGetContractUseCase(artifactAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		artifactAgent.EXPECT().
			FindOneByNameAndTag(gomock.Any(), contract.Name, contract.Tag).
			Return(&models.ArtifactModel{
				ID:               1,
				ABI:              contract.ABI,
				Bytecode:         contract.Bytecode.String(),
				DeployedBytecode: contract.DeployedBytecode.String(),
				Codehash:         "",
			}, nil)

		response, err := usecase.Execute(ctx, contract.Name, contract.Tag)

		assert.NoError(t, err)
		assert.Equal(t, contract.Bytecode, response.Bytecode)
		assert.Equal(t, contract.DeployedBytecode, response.DeployedBytecode)
		assert.Equal(t, contract.ABI, response.ABI)
		assert.Equal(t, contract.Constructor, response.Constructor)
		assert.Len(t, response.Methods, 11)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		artifactAgent.EXPECT().FindOneByNameAndTag(gomock.Any(), contract.Name, contract.Tag).Return(nil, dataAgentError)

		response, err := usecase.Execute(ctx, contract.Name, contract.Tag)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(getContractComponent), err)
	})
}
