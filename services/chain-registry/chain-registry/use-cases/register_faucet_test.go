package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

func TestRegisterFaucet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	faucetAgent := mockstore.NewMockFaucetAgent(mockCtrl)

	registerFaucetUC := NewRegisterFaucet(faucetAgent)

	faucet := &models.Faucet{
		Name: "faucetName",
	}

	faucetAgent.EXPECT().RegisterFaucet(gomock.Any(), gomock.Eq(faucet)).Times(1)

	err := registerFaucetUC.Execute(context.Background(), faucet)
	assert.Nil(t, err)
}
