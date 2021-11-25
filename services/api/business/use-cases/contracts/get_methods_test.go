// +build unit

package contracts

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models"
)

func TestGetMethods_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	selector := []byte{58, 56}
	method := &models.MethodModel{
		ABI: "eventABI",
	}

	methodAgent := mocks.NewMockMethodAgent(ctrl)
	usecase := NewGetMethodsUseCase(methodAgent)

	t.Run("should execute use case successfully if method is found", func(t *testing.T) {
		methodAgent.EXPECT().FindOneByAccountAndSelector(gomock.Any(), chainID, contractAddress.Hex(), selector).Return(method, nil)

		responseABI, eventsABI, err := usecase.Execute(ctx, chainID, contractAddress, selector)

		assert.Equal(t, responseABI, method.ABI)
		assert.Nil(t, eventsABI)
		assert.NoError(t, err)
	})

	t.Run("should fail if data agent returns connection error", func(t *testing.T) {
		pgError := errors.PostgresConnectionError("error")
		methodAgent.EXPECT().FindOneByAccountAndSelector(gomock.Any(), chainID, contractAddress.Hex(), selector).Return(nil, pgError)

		responseABI, eventsABI, err := usecase.Execute(ctx, chainID, contractAddress, selector)

		assert.Equal(t, errors.FromError(pgError).ExtendComponent(getMethodsComponent), err)
		assert.Empty(t, responseABI)
		assert.Nil(t, eventsABI)
	})

	t.Run("should execute use case successfully if method is not found", func(t *testing.T) {
		methodAgent.EXPECT().FindOneByAccountAndSelector(gomock.Any(), chainID, contractAddress.Hex(), selector).Return(nil, nil)
		methodAgent.EXPECT().FindDefaultBySelector(gomock.Any(), selector).Return([]*models.MethodModel{method, method}, nil)

		responseABI, eventsABI, err := usecase.Execute(ctx, chainID, contractAddress, selector)

		assert.Equal(t, eventsABI, []string{method.ABI, method.ABI})
		assert.Empty(t, responseABI)
		assert.NoError(t, err)
	})

	t.Run("should fail if data agent returns error on find default", func(t *testing.T) {
		pgError := errors.PostgresConnectionError("error")
		methodAgent.EXPECT().FindOneByAccountAndSelector(gomock.Any(), chainID, contractAddress.Hex(), selector).Return(nil, nil)
		methodAgent.EXPECT().FindDefaultBySelector(gomock.Any(), selector).Return(nil, pgError)

		responseABI, eventsABI, err := usecase.Execute(ctx, chainID, contractAddress, selector)

		assert.Equal(t, errors.FromError(pgError).ExtendComponent(getMethodsComponent), err)
		assert.Empty(t, responseABI)
		assert.Nil(t, eventsABI)
	})
}
