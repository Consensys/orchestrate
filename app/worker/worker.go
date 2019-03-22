package worker

import (
	commonHandlers "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/app/worker/handlers"
)

// CreateWorker creates worker and attach it to application
func CreateWorker(infra *infra.Infra, marker services.OffsetMarker) *worker.Worker {
	// Instantiate worker
	w := worker.NewWorker(worker.NewConfig())

	// Handler::loader
	w.Use(commonHandlers.Loader(infra.Unmarshaller))

	// Handler::logger
	w.Use(handlers.Logger)

	// Handler::marker
	w.Use(commonHandlers.Marker(marker))

	// Handler::nonce
	w.Use(
		handlers.NonceHandler(
			infra.NonceManager,
			infra.Mec.PendingNonceAt,
		),
	)

	// Handler::Producer
	w.Use(commonHandlers.Producer(infra.Producer))

	return w
}
