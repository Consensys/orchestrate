package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	genuuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
)

func TestDeleteChain_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	deleteChainUC := NewDeleteChain(chainAgent)
	uuid := genuuid.NewV4().String()

	chainAgent.EXPECT().DeleteChainByUUID(gomock.Any(), gomock.Eq(uuid)).Times(1)

	err := deleteChainUC.Execute(context.Background(), uuid, "")
	assert.Nil(t, err)
}

func TestDeleteChain_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	chainAgent := mockstore.NewMockChainAgent(mockCtrl)

	deleteChainUC := NewDeleteChain(chainAgent)
	uuid := genuuid.NewV4().String()
	tenantID := "tenantID_1"

	chainAgent.EXPECT().DeleteChainByUUIDAndTenant(gomock.Any(), gomock.Eq(uuid), gomock.Eq(tenantID)).Times(1)

	err := deleteChainUC.Execute(context.Background(), uuid, tenantID)
	assert.Nil(t, err)
}
