// +build unit

package contracts

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/services/api/store/mocks"
)

func TestGetCatalog_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryAgent := mocks.NewMockRepositoryAgent(ctrl)
	usecase := NewGetCatalogUseCase(repositoryAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		names := []string{"Contract0", "Contract1"}
		repositoryAgent.EXPECT().FindAll(gomock.Any()).Return(names, nil)

		response, err := usecase.Execute(context.Background())

		assert.Equal(t, response, names)
		assert.NoError(t, err)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		repositoryAgent.EXPECT().FindAll(gomock.Any()).Return(nil, dataAgentError)

		response, err := usecase.Execute(context.Background())

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(getCatalogComponent), err)
	})
}
