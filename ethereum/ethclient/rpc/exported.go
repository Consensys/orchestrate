package rpc

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	httpclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client"
)

const component = "ethclient.rpc"

var (
	client   *Client
	clientV2 *ClientV2
	config   *Config
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		if config == nil {
			config = NewConfig()
		}

		// Set Client
		httpclient.Init(ctx)
		clientV2 = NewClientV2(config, httpclient.GlobalClient())
		client = NewClient(clientV2)

		rpcUrls := viper.GetStringSlice(URLViperKey)
		log.Infof("Connecting to %d RPC URLs", len(rpcUrls))

		var chains []string
		for _, url := range viper.GetStringSlice(URLViperKey) {
			chainID, err := client.Dial(ctx, url)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"eth-client-url": url,
					"error":          err,
				}).Fatalf("ethereum: could not dial client")
			}
			log.Infof("Created client for chain id %s for RPC URL %s", chainID.String(), url)
			chains = append(chains, chainID.String())
		}

		log.WithFields(log.Fields{
			"chains": chains,
		}).Infof("ethereum: multi-client ready")

	})
}

// GlobalClient returns global Client
func GlobalClient() *Client {
	return client
}

// GlobalClient returns global Client
func GlobalClientV2() *ClientV2 {
	return clientV2
}

// SetGlobalClient sets global Client
func SetGlobalMultiClient(ec *Client) {
	client = ec
}
