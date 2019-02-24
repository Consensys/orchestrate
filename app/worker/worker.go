package worker

import (
	"github.com/spf13/viper"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/app/worker/handlers"
)

// CreateWorker creates worker and attach it to application
func CreateWorker(infra *infra.Infra, marker services.OffsetMarker) *core.Worker {
	// Instantiate worker
	w := core.NewWorker(uint(viper.GetInt("worker.slots")))

	// Handler::loader
	w.Use(handCom.Loader(infra.Unmarshaller))

	// Handler::logger
	w.Use(handlers.Logger)

	// Handler::marker
	w.Use(handCom.Marker(marker))

	// Handler::Sender
	w.Use(handlers.Sender(infra.TxSender))

	return w
}
