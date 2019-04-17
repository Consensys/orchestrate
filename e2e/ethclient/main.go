package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})

	// Initialize Multi-Client
	ethclient.Init(context.Background())

	// Retrieve last block on every client
	mec := ethclient.MultiClient()
	for _, chain := range mec.Networks(context.Background()) {
		header, err := mec.HeaderByNumber(context.Background(), chain, nil)
		if err != nil {
			log.WithError(err).Errorf("Error retrieving block header")
			continue
		}

		log.WithFields(log.Fields{
			"chain.id":      chain.Text(16),
			"header.number": header.Number.Text(10),
		}).Infof("Latest block")
	}
}
