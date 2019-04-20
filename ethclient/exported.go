package ethclient

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	mec = NewMultiEthClient()
}

var mec *MultiEthClient

// Init initialize Dials chains
//
// Ethereum clients URLs to Dial are read from viper configuration
// Cancelling input Context will stop multiclient
// If an error occurs during initialization, it will panic
func Init(ctx context.Context) {
	// Dial Ethereum client (URLs found in viper configuration)
	err := mec.MultiDial(ctx, viper.GetStringSlice(urlViperKey))
	if err != nil {
		log.WithError(err).Fatalf("ethereum: could not dial multi-client")
	}
	chainIDs := mec.Networks(ctx)
	log.Infof("ethereum: multi-client ready (connected to chains: %v)", chainIDs)

	// Wait for context to be done and then close
	go func() {
		<-ctx.Done()
		mec.Close()
	}()

}

// GlobalMultiClient returns global MultiEthClient
func GlobalMultiClient() *MultiEthClient {
	return mec
}
