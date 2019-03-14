package worker

import (
	commonHandlers "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/handlers"
	coreWorker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	coreServices "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/app/worker/handlers"
)

// CreateWorker creates worker and attach it to application
func CreateWorker(infra *infra.Infra, marker coreServices.OffsetMarker) *coreWorker.Worker {
	// Instantiate worker
	w := coreWorker.NewWorker(coreWorker.NewConfig())

	// Handler::loader
	w.Use(commonHandlers.Loader(infra.Unmarshaller))

	// Handler::logger
	w.Use(handlers.Logger)

	// Handler::marker
	w.Use(commonHandlers.Marker(marker))

	// Handler::decoder
	w.Use(handlers.Decoder(infra.ABIRegistry))

	// Handler::Producer
	w.Use(commonHandlers.Producer(infra.Producer))

	return w
}
