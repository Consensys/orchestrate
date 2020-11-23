package txupdater

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Tx Updater Handler
func Init() {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize client
		txscheduler.Init()

		// Create Handler
		handler = TransactionUpdater(txscheduler.GlobalClient())

		log.Infof("logger: transaction updater handler ready")
	})
}

// SetGlobalHandler sets global OpenTracing Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global OpenTracing handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
