package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	genuuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestGetFaucet_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetUC := NewGetFaucet(faucetAgent)
	uuid := genuuid.NewV4().String()

	expectedFaucet := &models.Faucet{
		UUID: uuid,
		Name: "testFaucet",
	}
	faucetAgent.EXPECT().GetFaucetByUUID(gomock.Any(), gomock.Eq(uuid)).Return(expectedFaucet, nil).Times(1)

	actualFaucet, err := getFaucetUC.Execute(context.Background(), uuid, "")
	assert.Nil(t, err)
	assert.Equal(t, expectedFaucet, actualFaucet)
}

func TestGetFaucet_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetUC := NewGetFaucet(faucetAgent)
	uuid := genuuid.NewV4().String()
	tenantID := "tenantID_5"

	expectedFaucet := &models.Faucet{
		UUID:     uuid,
		TenantID: tenantID,
		Name:     "testFaucet",
	}
	faucetAgent.EXPECT().GetFaucetByUUIDAndTenant(gomock.Any(), gomock.Eq(uuid), tenantID).Return(expectedFaucet, nil).Times(1)

	actualFaucet, err := getFaucetUC.Execute(context.Background(), uuid, tenantID)
	assert.Nil(t, err)
	assert.Equal(t, expectedFaucet, actualFaucet)
}
