package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mocks"
)

func TestGetCatalog_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepositoryDataAgent := mocks.NewMockRepositoryDataAgent(ctrl)
	usecase := NewGetCatalog(mockRepositoryDataAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		names := []string{"Contract0", "Contract1"}
		mockRepositoryDataAgent.EXPECT().FindAll(context.Background()).Return(names, nil)

		response, err := usecase.Execute(context.Background())

		assert.Equal(t, response, names)
		assert.Nil(t, err)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		mockRepositoryDataAgent.EXPECT().FindAll(context.Background()).Return(nil, dataAgentError)

		response, err := usecase.Execute(context.Background())

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(getCatalogComponent), err)
	})
}
