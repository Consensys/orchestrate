package usecases

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestUpdateFaucet_ByUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	updateChainUC := NewUpdateFaucet(faucetAgent)
	faucetUUID := uuid.Must(uuid.NewV4()).String()

	faucet := &models.Faucet{
		Name: "faucetName",
	}

	faucetAgent.EXPECT().UpdateFaucetByUUID(gomock.Any(), gomock.Eq(faucetUUID), gomock.Eq(faucet)).Times(1)

	err := updateChainUC.Execute(context.Background(), faucetUUID, faucet)
	assert.Nil(t, err)
}

func TestUpdateFaucet_FailUpdateUUID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	updateChainUC := NewUpdateFaucet(faucetAgent)
	faucetUUID := uuid.Must(uuid.NewV4()).String()

	faucet := &models.Faucet{
		UUID: uuid.Must(uuid.NewV4()).String(),
		Name: "faucetName",
	}

	err := updateChainUC.Execute(context.Background(), faucetUUID, faucet)
	assert.NotNil(t, err)
	assert.True(t, errors.IsConstraintViolatedError(err), "should be IsConstraintViolatedError")
}
