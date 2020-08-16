package sender

import (
	"context"
	"sync"

	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/envelope/storer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"
)

const component = "handler.sender"

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Tx Scheduler client
		txscheduler.Init()

		// Initialize Ethereum client
		ethclient.Init(ctx)

		ec := ethclient.GlobalClient()

		// Create Handler
		handler = engine.CombineHandlers(
			// Idempotency gate
			storer.TxAlreadySent(ec, txscheduler.GlobalClient()),
			// Sender
			Sender(ec, txscheduler.GlobalClient()),
		)

		log.Infof("sender: handler ready")
	})
}

// SetGlobalHandler sets global Sender Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Sender handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
