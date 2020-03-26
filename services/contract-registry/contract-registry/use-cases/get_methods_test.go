package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/testutils"
)

func TestGetMethods_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	account := testutils.FakeAccount()
	selector := []byte{58, 56}
	method := &models.MethodModel{
		ABI: "eventABI",
	}

	mockMethodDataAgent := mocks.NewMockMethodDataAgent(ctrl)
	usecase := NewGetMethods(mockMethodDataAgent)

	t.Run("should execute use case successfully if method is found", func(t *testing.T) {
		mockMethodDataAgent.EXPECT().FindOneByAccountAndSelector(context.Background(), account, selector).Return(method, nil)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, selector)

		assert.Equal(t, responseABI, method.ABI)
		assert.Nil(t, eventsABI)
		assert.Nil(t, err)
	})

	t.Run("should fail if data agent returns connection error", func(t *testing.T) {
		pgError := errors.PostgresConnectionError("error")
		mockMethodDataAgent.EXPECT().FindOneByAccountAndSelector(context.Background(), account, selector).Return(nil, pgError)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, selector)

		assert.Equal(t, errors.FromError(pgError).ExtendComponent(getMethodsComponent), err)
		assert.Empty(t, responseABI)
		assert.Nil(t, eventsABI)
	})

	t.Run("should execute use case successfully if method is not found", func(t *testing.T) {
		mockMethodDataAgent.EXPECT().FindOneByAccountAndSelector(context.Background(), account, selector).Return(nil, nil)
		mockMethodDataAgent.EXPECT().FindDefaultBySelector(context.Background(), selector).Return([]*models.MethodModel{method, method}, nil)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, selector)

		assert.Equal(t, eventsABI, []string{method.ABI, method.ABI})
		assert.Empty(t, responseABI)
		assert.Nil(t, err)
	})

	t.Run("should fail if data agent returns error on find default", func(t *testing.T) {
		pgError := errors.PostgresConnectionError("error")
		mockMethodDataAgent.EXPECT().FindOneByAccountAndSelector(context.Background(), account, selector).Return(nil, nil)
		mockMethodDataAgent.EXPECT().FindDefaultBySelector(context.Background(), selector).Return(nil, pgError)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, selector)

		assert.Equal(t, errors.FromError(pgError).ExtendComponent(getMethodsComponent), err)
		assert.Empty(t, responseABI)
		assert.Nil(t, eventsABI)
	})
}
