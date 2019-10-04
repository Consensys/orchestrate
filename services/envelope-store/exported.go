package envelopestore

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/memory"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/pg"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
)

const (
	component   = "envelope-store"
	postgresOpt = "postgres"
	inMemoryOpt = "in-memory"
)

var (
	store    evlpstore.EnvelopeStoreServer
	initOnce = &sync.Once{}
)

func Init() {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		switch viper.GetString(typeViperKey) {
		case postgresOpt:
			// Initialize Sarama Faucet
			pg.Init()

			// Set Faucet
			store = pg.GlobalEnvelopeStore()
		case inMemoryOpt:
			// Initialize Mock Faucet
			memory.Init()

			// Set Faucet
			store = memory.GlobalEnvelopeStore()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("%s: unknown type", component)
		}
	})
}

func GlobalEnvelopeStoreServer() evlpstore.EnvelopeStoreServer {
	return store
}

// SetGlobalEnvelopeStoreServer sets EnvelopeStoreServer
func SetGlobalEnvelopeStoreServer(s evlpstore.EnvelopeStoreServer) {
	store = s
}
