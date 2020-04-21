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

func TestGetFaucets_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetsUC := NewGetFaucets(faucetAgent)
	uuid := genuuid.NewV4().String()

	filters := make(map[string]string)
	expectedFaucet := []*models.Faucet{
		&models.Faucet{
			UUID: uuid,
			Name: "testFaucet",
		},
	}
	faucetAgent.EXPECT().GetFaucets(gomock.Any(), gomock.Eq(filters)).Return(expectedFaucet, nil).Times(1)

	actualFaucets, err := getFaucetsUC.Execute(context.Background(), "", filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedFaucet, actualFaucets)
}

func TestGetFaucets_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	getFaucetsUC := NewGetFaucets(faucetAgent)
	uuid := genuuid.NewV4().String()
	tenantID := "tenantID_6"

	filters := make(map[string]string)
	expectedFaucet := []*models.Faucet{
		&models.Faucet{
			UUID:     uuid,
			TenantID: tenantID,
			Name:     "testFaucet",
		},
	}
	faucetAgent.EXPECT().GetFaucetsByTenant(gomock.Any(), gomock.Eq(filters), tenantID).Return(expectedFaucet, nil).Times(1)

	actualFaucets, err := getFaucetsUC.Execute(context.Background(), tenantID, filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedFaucet, actualFaucets)
}
