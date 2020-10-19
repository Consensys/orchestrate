// +build unit

package identitymanager

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	mock2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client/mock"
	mock3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client/mock"
)

func TestApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrl)
	keyChecker := mockauth.NewMockChecker(ctrl)
	keyManagerClient := mock.NewMockKeyManagerClient(ctrl)
	chainRegistryClient := mock2.NewMockChainRegistryClient(ctrl)
	txSchedulerClient := mock3.NewMockTransactionSchedulerClient(ctrl)
	

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "postgres"

	keyManagerClient.EXPECT().Checker().Return(func() error {return nil})
	chainRegistryClient.EXPECT().Checker().Return(func() error {return nil})
	txSchedulerClient.EXPECT().Checker().Return(func() error {return nil})
	_, err := NewIdentityManager(cfg, postgres.GetManager(), jwtChecker, keyChecker, keyManagerClient, chainRegistryClient, txSchedulerClient)
	assert.NoError(t, err, "Creating App should not error")
}
