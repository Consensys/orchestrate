package rpc

import (
	"context"
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/rpc/geth"
)

var (
	client   *Client
	config   *geth.Config
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		if config == nil {
			config = geth.NewConfig()
		}

		client = NewClient(config)
		rpcUrls := viper.GetStringSlice(urlViperKey)
		log.Infof("Connecting to %d RPC URLs", len(rpcUrls))

		chains := []string{}
		for _, url := range viper.GetStringSlice(urlViperKey) {
			chainID, err := client.Dial(ctx, url)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"eth-client": url,
					"error":      err,
				}).Fatalf("ethereum: could not dial client")
			}
			log.Infof("Chain id for RPC URL %s is  %s", url, chainID.String())
			chains = append(chains, chainID.String())
		}

		log.WithFields(log.Fields{
			"chains": chains,
		}).Infof("ethereum: multi-client ready")
	})
}

// Dial
func Dial(ctx context.Context, rawurl string) (*big.Int, error) {
	return client.Dial(ctx, rawurl)
}

// GlobalClient returns global Client
func GlobalClient() *Client {
	return client
}

// SetGlobalClient sets global Client
func SetGlobalMultiClient(ec *Client) {
	client = ec
}
