package mock

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	envelopeStore *EnvelopeStore
	initOnce      = &sync.Once{}
)

// InitStore initilialize envelope store
func initStore() {
	envelopeStore = NewEnvelopeStore()
}

// Init initialize Sender Handler
func Init() {
	initOnce.Do(func() {
		if envelopeStore != nil {
			return
		}

		// Initialize Grpc store
		initStore()

		log.Infof("mock: store ready")
	})
}

func GlobalEnvelopeStore() *EnvelopeStore {
	return envelopeStore
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStore(s *EnvelopeStore) {
	envelopeStore = s
}
