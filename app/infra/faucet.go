package infra

import (
	log "github.com/sirupsen/logrus"
	infraFaucet "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/infra"
)

func initFaucet(infra *Infra) {
	crediter, err := infraFaucet.NewSaramaCrediter(infra.SaramaProducer)
	if err != nil {
		log.WithError(err).Fatalf("infra-faucet: could not create faucet")
	}

	faucet, err := infraFaucet.CreateFaucet(infra.Mec.PendingBalanceAt, crediter.Credit)
	if err != nil {
		log.WithError(err).Fatalf("infra-faucet: could not create faucet")
	}

	infra.Faucet = faucet
}
