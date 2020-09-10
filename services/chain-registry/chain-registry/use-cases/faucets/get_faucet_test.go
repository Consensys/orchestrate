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

func TestGetFaucet_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetUC := NewGetFaucet(faucetAgent)
	faucetUUID := uuid.Must(uuid.NewV4()).String()

	expectedFaucet := &models.Faucet{
		UUID: faucetUUID,
		Name: "testFaucet",
	}
	faucetAgent.EXPECT().GetFaucet(gomock.Any(), gomock.Eq(faucetUUID), []string{}).Return(expectedFaucet, nil).Times(1)

	actualFaucet, err := getFaucetUC.Execute(context.Background(), faucetUUID, []string{})
	assert.NoError(t, err)
	assert.Equal(t, expectedFaucet, actualFaucet)
}

func TestGetFaucet_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetUC := NewGetFaucet(faucetAgent)
	faucetUUID := uuid.Must(uuid.NewV4()).String()
	tenantID := "tenantID_5"

	expectedFaucet := &models.Faucet{
		UUID:     faucetUUID,
		TenantID: tenantID,
		Name:     "testFaucet",
	}
	faucetAgent.EXPECT().GetFaucet(gomock.Any(), gomock.Eq(faucetUUID), []string{tenantID}).Return(expectedFaucet, nil).Times(1)

	actualFaucet, err := getFaucetUC.Execute(context.Background(), faucetUUID, []string{tenantID})
	assert.NoError(t, err)
	assert.Equal(t, expectedFaucet, actualFaucet)
}
