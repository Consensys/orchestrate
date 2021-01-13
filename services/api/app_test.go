// +build unit

package api

import (
	ethclientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/mock"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"

	"github.com/Shopify/sarama/mocks"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	keymanagerclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client/mock"
)

func TestApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "postgres"

	kCfg := sarama.NewKafkaTopicConfig(viper.New())
	_, err := NewAPI(
		cfg,
		postgres.GetManager(),
		mockauth.NewMockChecker(ctrl), mockauth.NewMockChecker(ctrl),
		keymanagerclient.NewMockKeyManagerClient(ctrl),
		ethclientmock.NewMockClient(ctrl),
		mocks.NewSyncProducer(t, nil),
		kCfg,
	)
	assert.NoError(t, err, "Creating App should not error")
}
