package infra

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
)

func initEthereum(infra *Infra, wait *sync.WaitGroup) {
	// Create Ethereum multi-client
	mec, err := ethclient.MultiDial(viper.GetStringSlice("eth.clients"))

	if err != nil {
		log.WithError(err).Fatalf("infra-ethereum: could not dial multi-client")
	}

	chainIDs := mec.Networks(context.Background())
	log.Infof("infra-ethereum: multi-client ready (connected to chains: %v)", chainIDs)

	// Attach Ethereum client and sender
	infra.Mec = mec
	infra.TxSender = mec

	// Wait for app to be done and then close
	go func() {
		<-infra.ctx.Done()
		// TODO: impelent mec.Close() and then close
	}()

	wait.Done()
}
