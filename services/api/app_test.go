// +build unit

package api

import (
	ethclientmock "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/mock"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/broker/sarama"

	"github.com/Shopify/sarama/mocks"

	mockauth "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/mock"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	keymanagerclient "github.com/ConsenSys/orchestrate/services/key-manager/client/mock"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
