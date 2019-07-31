package mock

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "envelope-store.mock"

var (
	store    *EnvelopeStore
	initOnce = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Initialize Grpc store
		store = NewEnvelopeStore()

		log.Infof("envelope-store.mock: store ready")
	})
}

func GlobalEnvelopeStore() *EnvelopeStore {
	return store
}

// SetGlobalEnvelopeStore set global mock store
func SetGlobalEnvelopeStore(s *EnvelopeStore) {
	store = s
}
