package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	grpcclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/client"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

const component = "contract-registry.client"

var (
	client   svc.ContractRegistryClient
	conn     *grpc.ClientConn
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		var err error
		conn, err = grpcclient.DialContextWithDefaultOptions(
			ctx,
			viper.GetString(ContractRegistryURLViperKey),
		)
		if err != nil {
			log.WithError(errors.FromError(err).ExtendComponent(component)).Fatalf("%s: failed to dial grpc server", component)
		}

		client = svc.NewContractRegistryClient(conn)

		log.WithFields(log.Fields{
			"url": viper.GetString(ContractRegistryURLViperKey),
		}).Infof("%s: client ready", component)
	})
}

func Close() {
	_ = conn.Close()
}

func GlobalClient() svc.ContractRegistryClient {
	return client
}

// SetGlobalClient sets ContractRegistry global configuration
func SetGlobalClient(c svc.ContractRegistryClient) {
	client = c
}
