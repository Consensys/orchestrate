package infra

import (
	log "github.com/sirupsen/logrus"
	ethabi "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/abi"
)

func initCrafter(infra *Infra) {
	// Handler::Crafter
	contracts, err := ethabi.FromABIConfig()
	if err != nil {
		log.WithError(err).Fatalf("infra-crafter: could not initialize crafter")
	}

	// Attach crafter and ABI registry
	infra.Crafter = &ethabi.PayloadCrafter{}
	registry := ethabi.NewStaticRegistry()
	for _, contract := range contracts {
		registry.RegisterContract(contract)
	}
	infra.ABIRegistry = registry
}
