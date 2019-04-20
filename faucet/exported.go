package faucet

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet/sarama"
)

var (
	fct      Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		switch viper.GetString("faucet") {
		case "sarama":
			// Initialize Sarama Faucet
			sarama.Init(ctx)

			// Set Faucet
			fct = sarama.GlobalFaucet()
		default:
			// Initialize Mock Faucet
			mock.Init(ctx)

			// Set Faucet
			fct = mock.GlobalFaucet()
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
}
