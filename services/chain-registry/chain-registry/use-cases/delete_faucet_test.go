package usecases

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
)

func TestDeleteFaucet_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	deleteAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	registerFaucetUC := NewDeleteFaucet(deleteAgent)
	faucetUUID := uuid.Must(uuid.NewV4()).String()

	deleteAgent.EXPECT().DeleteFaucet(gomock.Any(), gomock.Eq(faucetUUID), []string{}).Times(1)

	err := registerFaucetUC.Execute(context.Background(), faucetUUID, []string{})
	assert.NoError(t, err)
}

func TestDeleteFaucet_ByUUIDAndTenantID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	deleteFaucetUC := NewDeleteFaucet(faucetAgent)
	faucetUUID := uuid.Must(uuid.NewV4()).String()
	tenantID := "tenantID_2"

	faucetAgent.EXPECT().DeleteFaucet(gomock.Any(), gomock.Eq(faucetUUID), []string{tenantID}).Times(1)

	err := deleteFaucetUC.Execute(context.Background(), faucetUUID, []string{tenantID})
	assert.NoError(t, err)
}
