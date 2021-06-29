package api

import (
	"context"

	ethclient "github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient/rpc"

	qkm "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager"

	"github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	authjwt "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/jwt"
	authkey "github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/spf13/viper"
)

// New Utility function used to initialize a new service
func New(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	authjwt.Init(ctx)
	authkey.Init(ctx)
	sarama.InitSyncProducer(ctx)
	ethclient.Init(ctx)
	qkm.Init()

	config := NewConfig(viper.GetViper())
	pgmngr := postgres.GetManager()

	return NewAPI(
		config,
		pgmngr,
		authjwt.GlobalChecker(),
		authkey.GlobalChecker(),
		qkm.GlobalClient(),
		ethclient.GlobalClient(),
		sarama.GlobalSyncProducer(),
		sarama.NewKafkaTopicConfig(viper.GetViper()),
	)
}

func Run(ctx context.Context) error {
	appli, err := New(ctx)
	if err != nil {
		return err
	}
	return appli.Run(ctx)
}
