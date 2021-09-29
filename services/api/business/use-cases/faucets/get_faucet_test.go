package faucets

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/multitenancy"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	"github.com/consensys/orchestrate/services/api/store/mocks"
	"github.com/consensys/orchestrate/services/api/store/models/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetFaucet_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDB(ctrl)
	faucetAgent := mocks.NewMockFaucetAgent(ctrl)
	mockDB.EXPECT().Faucet().Return(faucetAgent).AnyTimes()

	usecase := NewGetFaucetUseCase(mockDB)

	tenantID := multitenancy.DefaultTenant
	tenants := []string{tenantID}

	t.Run("should execute use case successfully", func(t *testing.T) {
		faucet := testutils.FakeFaucetModel()
		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), faucet.UUID, tenants).Return(faucet, nil)

		resp, err := usecase.Execute(ctx, faucet.UUID, tenants)

		assert.NoError(t, err)
		assert.Equal(t, parsers.NewFaucetFromModel(faucet), resp)
	})

	t.Run("should fail with same error if get faucet fails", func(t *testing.T) {
		expectedErr := errors.NotFoundError("error")

		faucetAgent.EXPECT().FindOneByUUID(gomock.Any(), "uuid", tenants).Return(nil, expectedErr)

		resp, err := usecase.Execute(ctx, "uuid", tenants)

		assert.Nil(t, resp)
		assert.Error(t, err)
		assert.Equal(t, errors.FromError(expectedErr).ExtendComponent(getFaucetCandidateComponent), err)
	})
}
