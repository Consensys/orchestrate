// +build unit

package usecases

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/store/mock"
)

func TestGetTags_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	contractName := "myContract"
	mockTagDataAgent := mock.NewMockTagDataAgent(ctrl)
	usecase := NewGetTags(mockTagDataAgent)

	t.Run("should execute use case successfully", func(t *testing.T) {
		tags := []string{"latest", "v1.0.0"}
		mockTagDataAgent.EXPECT().FindAllByName(context.Background(), contractName).Return(tags, nil)

		response, err := usecase.Execute(context.Background(), contractName)

		assert.Equal(t, response, tags)
		assert.NoError(t, err)
	})

	t.Run("should fail if data agent fails", func(t *testing.T) {
		dataAgentError := fmt.Errorf("error")
		mockTagDataAgent.EXPECT().FindAllByName(context.Background(), contractName).Return(nil, dataAgentError)

		response, err := usecase.Execute(context.Background(), contractName)

		assert.Nil(t, response)
		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(getTagsComponent), err)
	})
}
