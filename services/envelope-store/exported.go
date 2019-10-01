package envelopestore

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/envelope-store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/envelope-store/pg"
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
		case "pg":
			// Initialize Sarama Faucet
			pg.Init()

			// Set Faucet
			store = pg.GlobalEnvelopeStore()
		case "mock":
			// Initialize Mock Faucet
			mock.Init()

			// Set Faucet
			store = mock.GlobalEnvelopeStore()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("envelope-store: unknown type")
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
