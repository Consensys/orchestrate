package pg

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

const component = "envelope-store.pg"

var (
	store    *EnvelopeStore
	initOnce = &sync.Once{}
)

// Init initialize Postgres Envelope Store
func Init() {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Initialize gRPC store
		opts := postgres.NewOptions()
		store = NewEnvelopeStoreFromPGOptions(opts)
		log.WithFields(log.Fields{
			"db.address":  opts.Addr,
			"db.database": opts.Database,
			"db.user":     opts.User,
		}).Infof("envelope-store.pg: store ready")
	})
}

func GlobalEnvelopeStore() *EnvelopeStore {
	return store
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStore(s *EnvelopeStore) {
	store = s
}
