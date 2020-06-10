package transactionscheduler

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client/mock"

	"github.com/Shopify/sarama/mocks"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	mockclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
)

func TestApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrl)
	keyChecker := mockauth.NewMockChecker(ctrl)
	mockSyncProducer := mocks.NewSyncProducer(t, nil)

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "postgres"

	_, err := New(
		cfg,
		postgres.GetManager(),
		jwtChecker, keyChecker,
		mockclient.NewMockChainRegistryClient(ctrl),
		mock.NewMockContractRegistryClient(ctrl),
		mockSyncProducer,
		"tx-crafter-topic",
	)
	assert.NoError(t, err, "Creating App should not error")
}
