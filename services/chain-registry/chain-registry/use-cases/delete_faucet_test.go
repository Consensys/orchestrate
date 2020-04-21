package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	genuuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
)

func TestDeleteFaucet_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	deleteAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	registerFaucetUC := NewDeleteFaucet(deleteAgent)
	uuid := genuuid.NewV4().String()

	deleteAgent.EXPECT().DeleteFaucetByUUID(gomock.Any(), gomock.Eq(uuid)).Times(1)

	err := registerFaucetUC.Execute(context.Background(), uuid, "")
	assert.Nil(t, err)
}

func TestDeleteFaucet_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	deleteFaucetUC := NewDeleteFaucet(faucetAgent)
	uuid := genuuid.NewV4().String()
	tenantID := "tenantID_2"

	faucetAgent.EXPECT().DeleteFaucetByUUIDAndTenant(gomock.Any(), gomock.Eq(uuid), gomock.Eq(tenantID)).Times(1)

	err := deleteFaucetUC.Execute(context.Background(), uuid, tenantID)
	assert.Nil(t, err)
}
