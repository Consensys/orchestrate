// +build unit

package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

func TestGetEvents_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	account := testutils.FakeAccount()
	sigHash := "sigHash"
	indexedInputCount := uint32(1)
	eventModel := &models.EventModel{
		ABI: "eventABI",
	}

	mockEventDataAgent := mock.NewMockEventDataAgent(ctrl)
	usecase := NewGetEvents(mockEventDataAgent)

	t.Run("should execute use case successfully if event is found", func(t *testing.T) {
		mockEventDataAgent.EXPECT().FindOneByAccountAndSigHash(context.Background(), account, sigHash, indexedInputCount).Return(eventModel, nil)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, sigHash, indexedInputCount)

		assert.Equal(t, responseABI, eventModel.ABI)
		assert.Nil(t, eventsABI)
		assert.Nil(t, err)
	})

	t.Run("should fail if data agent returns connection error", func(t *testing.T) {
		pgError := errors.PostgresConnectionError("error")
		mockEventDataAgent.EXPECT().FindOneByAccountAndSigHash(context.Background(), account, sigHash, indexedInputCount).Return(nil, pgError)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, sigHash, indexedInputCount)

		assert.Equal(t, errors.FromError(pgError).ExtendComponent(getEventsComponent), err)
		assert.Empty(t, responseABI)
		assert.Nil(t, eventsABI)
	})

	t.Run("should execute use case successfully if event is not found", func(t *testing.T) {
		mockEventDataAgent.EXPECT().FindOneByAccountAndSigHash(context.Background(), account, sigHash, indexedInputCount).Return(nil, nil)
		mockEventDataAgent.EXPECT().FindDefaultBySigHash(context.Background(), sigHash, indexedInputCount).Return([]*models.EventModel{eventModel, eventModel}, nil)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, sigHash, indexedInputCount)

		assert.Equal(t, eventsABI, []string{eventModel.ABI, eventModel.ABI})
		assert.Empty(t, responseABI)
		assert.Nil(t, err)
	})

	t.Run("should fail if data agent returns error on find default", func(t *testing.T) {
		pgError := errors.PostgresConnectionError("error")
		mockEventDataAgent.EXPECT().FindOneByAccountAndSigHash(context.Background(), account, sigHash, indexedInputCount).Return(nil, nil)
		mockEventDataAgent.EXPECT().FindDefaultBySigHash(context.Background(), sigHash, indexedInputCount).Return(nil, pgError)

		responseABI, eventsABI, err := usecase.Execute(context.Background(), account, sigHash, indexedInputCount)

		assert.Equal(t, errors.FromError(pgError).ExtendComponent(getEventsComponent), err)
		assert.Empty(t, responseABI)
		assert.Nil(t, eventsABI)
	})
}
