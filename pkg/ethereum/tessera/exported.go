package tessera

import (
	"context"
	"sync"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	client   Client
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		client = NewEnclaveClient()
		tesseraEndpoints := viper.GetStringMapString(URLsViperKey)
		log.Infof("connecting to %d Tessera URLs", len(tesseraEndpoints))

		for chainID, endpoint := range tesseraEndpoints {
			log.Infof("adding Tessera client for endpoint '%s' for chain id %s", endpoint, chainID)

			enclaveHTTPClient := CreateEnclaveHTTPEndpoint(endpoint)
			client.AddClient(chainID, enclaveHTTPClient)

			checkIfEndpointAccessible(chainID, endpoint)
		}
	})
}

func checkIfEndpointAccessible(chainID, endpoint string) {
	_, err := client.GetStatus(chainID)
	if err != nil {
		log.Errorf("status check failed for Tessera endpoint '%s' with error %s", endpoint, err)
	}
}

// GlobalClient returns global Client
func GlobalClient() Client {
	return client
}

// SetGlobalClient sets global Client
func SetGlobalMultiClient(ec Client) {
	client = ec
}
