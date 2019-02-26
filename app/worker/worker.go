package worker

import (
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/app/worker/handlers"
	infList "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/infra"
)

// CreateWorker creates worker and attach it to application
func CreateWorker(infra *infra.Infra) *core.Worker {
	// Instantiate worker
	w := core.NewWorker(1)

	// Handler::loader
	w.Use(handCom.Loader(&infList.ReceiptUnmarshaller{}))

	// Handler::logger
	w.Use(handlers.Logger)

	// Handler::Producer
	w.Use(handCom.Producer(infra.Producer))

	return w
}
