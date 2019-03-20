package worker

import (
	handcommon "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/app/worker/handlers"
	inflistener "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/infra"
)

// CreateWorker creates worker and attach it to application
func CreateWorker(infra *infra.Infra) *worker.Worker {
	// Instantiate worker
	w := worker.NewWorker(worker.NewConfig())

	// Handler::loader
	w.Use(handcommon.Loader(&inflistener.ReceiptUnmarshaller{}))

	// Handler::logger
	w.Use(handlers.Logger)

	// Handler::Producer
	w.Use(handcommon.Producer(infra.Producer))

	return w
}
