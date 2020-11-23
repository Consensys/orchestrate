package chains

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

func TestGetChain_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	getChainUC := NewGetChain(chainAgent)
	chainUUID := uuid.Must(uuid.NewV4()).String()

	expectedChain := &models.Chain{
		UUID: chainUUID,
		Name: "testChain",
	}
	chainAgent.EXPECT().GetChain(gomock.Any(), gomock.Eq(chainUUID), []string{}).Return(expectedChain, nil).Times(1)

	actualChain, err := getChainUC.Execute(context.Background(), chainUUID, []string{})
	assert.NoError(t, err)
	assert.Equal(t, expectedChain, actualChain)
}

func TestGetChain_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	getChainUC := NewGetChain(chainAgent)
	chainUUID := uuid.Must(uuid.NewV4()).String()
	tenantID := "tenantID_3"

	expectedChain := &models.Chain{
		UUID:     chainUUID,
		TenantID: tenantID,
		Name:     "testChain",
	}
	chainAgent.EXPECT().GetChain(gomock.Any(), gomock.Eq(chainUUID), []string{tenantID}).Return(expectedChain, nil).Times(1)

	actualChain, err := getChainUC.Execute(context.Background(), chainUUID, []string{tenantID})
	assert.NoError(t, err)
	assert.Equal(t, expectedChain, actualChain)
}
