package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client/dialer"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	gogrpc "google.golang.org/grpc"
)

const component = "contract-registry.client"

var (
	client   svc.ContractRegistryClient
	conn     *gogrpc.ClientConn
	initOnce = &sync.Once{}
)

func Init(ctx context.Context, contractRegistryURL string) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		var err error
		client, err = dialer.DialContextWithDefaultOptions(ctx, contractRegistryURL)
		if err != nil {
			log.WithError(err).Fatalf("could not dial contract-registry server")
		}

		log.WithFields(log.Fields{
			"url.full": contractRegistryURL,
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
