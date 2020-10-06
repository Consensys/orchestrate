package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	healthz "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client/dialer"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
)

const component = "contract-registry.client"

var (
	client   svc.ContractRegistryClient
	checker  healthz.Check
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		contractRegistryURL := viper.GetString(GRPCURLViperKey)
		var err error
		client, err = dialer.DialContextWithDefaultOptions(ctx, contractRegistryURL)
		if err != nil {
			log.WithError(err).Fatalf("could not dial contract-registry server")
		}

		checker = healthz.HTTPGetCheck(fmt.Sprintf("%s/live", viper.GetString(MetricsURLViperKey)), time.Second)
		log.Infof("%s: client ready - url: %s", component, contractRegistryURL)
	})
}

func GlobalClient() svc.ContractRegistryClient {
	return client
}

func GlobalChecker() healthz.Check {
	return checker
}

func SetGlobalChecker(c healthz.Check) {
	checker = c
}

// SetGlobalClient sets ContractRegistry global configuration
func SetGlobalClient(c svc.ContractRegistryClient) {
	client = c
}
