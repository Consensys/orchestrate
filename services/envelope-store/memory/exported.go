package memory

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

const component = "envelope-store.in-memory"

var (
	store    *EnvelopeStore
	initOnce = &sync.Once{}
)

// Init initialize mock envelope store
func Init() {
	initOnce.Do(func() {
		if store != nil {
			return
		}

		// Initialize gRPC store
		store = NewEnvelopeStore()

		log.Infof("%s: store ready", component)
	})
}

func GlobalEnvelopeStore() *EnvelopeStore {
	return store
}

// SetGlobalEnvelopeStore set global mock store
func SetGlobalEnvelopeStore(s *EnvelopeStore) {
	store = s
}
