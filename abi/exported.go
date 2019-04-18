package abi

import (
	log "github.com/sirupsen/logrus"
)

var (
	crafter  Crafter
	registry Registry
)

func init() {
	crafter = &PayloadCrafter{}
	registry = NewStaticRegistry()
}

// InitRegistry initialize ABI Registry
func InitRegistry() {
	// Read ABIs from ABI viper configuration
	contracts, err := FromABIConfig()
	if err != nil {
		log.WithError(err).Fatalf("abi: could not initialize ABI registry")
	}

	// Register contracts
	for _, contract := range contracts {
		registry.RegisterContract(contract)
	}
}

// SetGlobalRegistry sets global ABI registry
func SetGlobalRegistry(r Registry) {
	registry = r
}

// GlobalRegistry returns global ABI registry
func GlobalRegistry() Registry {
	return registry
}

// SetGlobalCrafter sets global crafter
func SetGlobalCrafter(c Crafter) {
	crafter = c
}

// GlobalCrafter returns global ABI registry
func GlobalCrafter() Crafter {
	return crafter
}
