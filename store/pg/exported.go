package pg

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/database/postgres"
)

var (
	component     = "envelope-store.pg"
	envelopeStore *EnvelopeStore
	initOnce      = &sync.Once{}
)

// InitStore initilialize envelope store
func initStore() {
	opts := postgres.NewOptions()
	envelopeStore = NewEnvelopeStoreFromPGOptions(opts)
	log.WithFields(log.Fields{
		"db.address":  opts.Addr,
		"db.database": opts.Database,
		"db.user":     opts.User,
	}).Infof("envelope-store: postgres store ready")
}

// Init initialize Sender Handler
func Init() {
	initOnce.Do(func() {
		if envelopeStore != nil {
			return
		}

		// Initialize Grpc store
		initStore()

		log.Infof("grpc: store ready")
	})
}

func GlobalEnvelopeStore() *EnvelopeStore {
	return envelopeStore
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStore(s *EnvelopeStore) {
	envelopeStore = s
}
