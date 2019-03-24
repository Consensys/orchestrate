package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	ethURLs := []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	}
	log.Infof("Connecting to EthClients: %v", ethURLs)
	mec, err := ethclient.MultiDial(ethURLs)
	if err != nil {
		log.WithError(err).Fatalf("infra-ethereum: could not dial multi-client")
	}

	chains := []string{}
	for _, chain := range mec.Networks(context.Background()) {
		chains = append(chains, chain.Text(16))
	}

	log.Infof("infra-ethereum: multi-client ready (connected to chains: %v)", chains)
}
