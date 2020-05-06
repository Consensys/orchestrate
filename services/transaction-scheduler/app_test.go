package transactionscheduler

import (
	"testing"

	"github.com/Shopify/sarama/mocks"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	mockclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client/mock"
)

func TestApp(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrlr)
	keyChecker := mockauth.NewMockChecker(ctrlr)
	mockSyncProducer := mocks.NewSyncProducer(t, nil)

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "postgres"

	_, err := New(
		cfg,
		postgres.GetManager(),
		jwtChecker, keyChecker,
		mockclient.NewMockChainRegistryClient(ctrlr),
		mockSyncProducer,
		"tx-crafter-topic",
	)
	assert.NoError(t, err, "Creating App should not error")
}
