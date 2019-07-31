package pg

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/database/postgres"
)

var (
	component = "envelope-store.pg"
	store     *EnvelopeStore
	initOnce  = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Initialize Grpc store
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
