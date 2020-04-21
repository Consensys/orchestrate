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

func TestGetChain_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	getChainUC := NewGetChain(chainAgent)
	uuid := genuuid.NewV4().String()

	expectedChain := &models.Chain{
		UUID: uuid,
		Name: "testChain",
	}
	chainAgent.EXPECT().GetChainByUUID(gomock.Any(), gomock.Eq(uuid)).Return(expectedChain, nil).Times(1)

	actualChain, err := getChainUC.Execute(context.Background(), uuid, "")
	assert.Nil(t, err)
	assert.Equal(t, expectedChain, actualChain)
}

func TestGetChain_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	getChainUC := NewGetChain(chainAgent)
	uuid := genuuid.NewV4().String()
	tenantID := "tenantID_3"

	expectedChain := &models.Chain{
		UUID:     uuid,
		TenantID: tenantID,
		Name:     "testChain",
	}
	chainAgent.EXPECT().GetChainByUUIDAndTenant(gomock.Any(), gomock.Eq(uuid), tenantID).Return(expectedChain, nil).Times(1)

	actualChain, err := getChainUC.Execute(context.Background(), uuid, tenantID)
	assert.Nil(t, err)
	assert.Equal(t, expectedChain, actualChain)
}
