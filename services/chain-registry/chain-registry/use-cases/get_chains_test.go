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

func TestGetChains_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	getChainsUC := NewGetChains(chainAgent)
	uuid := genuuid.NewV4().String()

	filters := make(map[string]string)
	expectedChain := []*models.Chain{
		&models.Chain{
			UUID: uuid,
			Name: "testChain",
		},
	}
	chainAgent.EXPECT().GetChains(gomock.Any(), gomock.Eq(filters)).Return(expectedChain, nil).Times(1)

	actualChains, err := getChainsUC.Execute(context.Background(), "", filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedChain, actualChains)
}

func TestGetChains_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	getChainsUC := NewGetChains(chainAgent)
	uuid := genuuid.NewV4().String()
	tenantID := "tenantID_4"

	filters := make(map[string]string)
	expectedChain := []*models.Chain{
		&models.Chain{
			UUID:     uuid,
			TenantID: tenantID,
			Name:     "testChain",
		},
	}
	chainAgent.EXPECT().GetChainsByTenant(gomock.Any(), gomock.Eq(filters), tenantID).Return(expectedChain, nil).Times(1)

	actualChains, err := getChainsUC.Execute(context.Background(), tenantID, filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedChain, actualChains)
}
