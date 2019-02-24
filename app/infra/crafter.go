package infra

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	infraCraft "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/infra"
)

func parseAbis(abis []string) (map[string]string, error) {
	reg := `(?P<contract_name>[a-zA-Z0-9]+):(?P<abi>\[.+\])`
	pattern := regexp.MustCompile(reg)
	m := make(map[string]string)
	for _, abi := range abis {
		match := pattern.FindStringSubmatch(abi)
		if len(match) != 3 {
			return nil, fmt.Errorf("Could not parse abi (expected format %q): %v ", abi, reg)
		}
		m[match[1]] = match[2]
	}
	return m, nil
}

func initCrafter(infra *Infra) {
	// Handler::Crafter
	abis, err := parseAbis(viper.GetStringSlice("abis"))
	if err != nil {
		log.WithError(err).Fatalf("infra-crafter: could not initialize crafter")
	}

	// Attach crafter and ABI registry
	infra.Crafter = &ethereum.PayloadCrafter{}
	infra.ABIRegistry = infraCraft.LoadABIRegistry(abis)
}
