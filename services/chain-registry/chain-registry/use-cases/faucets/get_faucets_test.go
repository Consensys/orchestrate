package faucets

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestGetFaucets_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetsUC := NewGetFaucets(faucetAgent)

	filters := make(map[string]string)
	expectedFaucet := []*models.Faucet{
		{
			UUID: uuid.Must(uuid.NewV4()).String(),
			Name: "testFaucet",
		},
	}
	faucetAgent.EXPECT().GetFaucets(gomock.Any(), []string{}, gomock.Eq(filters)).Return(expectedFaucet, nil).Times(1)

	actualFaucets, err := getFaucetsUC.Execute(context.Background(), []string{}, filters)
	assert.NoError(t, err)
	assert.Equal(t, expectedFaucet, actualFaucets)
}

func TestGetFaucets_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetsUC := NewGetFaucets(faucetAgent)
	tenantID := "tenantID_6"

	filters := make(map[string]string)
	expectedFaucet := []*models.Faucet{
		{
			UUID:     uuid.Must(uuid.NewV4()).String(),
			TenantID: tenantID,
			Name:     "testFaucet",
		},
	}
	faucetAgent.EXPECT().GetFaucets(gomock.Any(), []string{tenantID}, gomock.Eq(filters)).Return(expectedFaucet, nil).Times(1)

	actualFaucets, err := getFaucetsUC.Execute(context.Background(), []string{tenantID}, filters)
	assert.NoError(t, err)
	assert.Equal(t, expectedFaucet, actualFaucets)
}
