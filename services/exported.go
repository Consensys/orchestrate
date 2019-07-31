package services

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	mock "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/services/mock"
	pg "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/services/pg"
)

var (
	store    evlpstore.StoreServer
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		switch viper.GetString(typeViperKey) {
		case "pg":
			// Initialize Sarama Faucet
			pg.Init(ctx)

			// Set Faucet
			store = pg.GlobalEnvelopeStore()
		case "mock":
			// Initialize Mock Faucet
			mock.Init(ctx)

			// Set Faucet
			store = mock.GlobalEnvelopeStore()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("envelope-store: unknown type")
		}
	})
}

func GlobalEnvelopeStoreServer() evlpstore.StoreServer {
	return store
}

// SetGlobalEnvelopeStoreServer sets EnvelopeStoreServer
func SetGlobalEnvelopeStoreServer(s evlpstore.StoreServer) {
	store = s
}
