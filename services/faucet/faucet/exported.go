package faucet

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet/sarama"
)

var (
	fct      Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if fct != nil {
			return
		}

		switch viper.GetString(typeViperKey) {
		case "sarama":
			// Initialize Sarama Faucet
			sarama.Init(ctx)

			// Set Faucet
			fct = sarama.GlobalFaucet()
		case "mock":
			// Initialize Mock Faucet
			mock.Init(ctx)

			// Set Faucet
			fct = mock.GlobalFaucet()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("faucet: unknown type")
		}
	})
}

// GlobalFaucet returns global Sarama Faucet
func GlobalFaucet() Faucet {
	return fct
}

// SetGlobalFaucet sets global Sarama Faucet
func SetGlobalFaucet(faucet Faucet) {
	fct = faucet
	log.Debug("faucet: set")
}
