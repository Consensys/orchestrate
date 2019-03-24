package worker

import (
	"context"
	"math/big"

	"github.com/spf13/viper"
	commonHandlers "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/worker"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/app/worker/handlers"
)

// CreateWorker creates worker and attach it to application
func CreateWorker(infra *infra.Infra, marker services.OffsetMarker) *worker.Worker {
	// Instantiate worker
	w := worker.NewWorker(context.Background(), worker.NewConfig())

	// Handler::loader
	w.Use(commonHandlers.Loader(infra.Unmarshaller))

	// Handler::logger
	w.Use(handlers.Logger)

	// Handler::marker
	w.Use(commonHandlers.Marker(marker))

	// Handler::Faucet
	creditAmount := big.NewInt(0)
	creditAmount.SetString(viper.GetString("faucet.amount"), 10)
	w.Use(handlers.Faucet(infra.Faucet, creditAmount))

	// Handler::Crafter
	w.Use(handlers.Crafter(infra.ABIRegistry, infra.Crafter))

	// Handler:Gas
	w.Use(handlers.GasPricer(infra.GasManager))    // Gas Price
	w.Use(handlers.GasEstimator(infra.GasManager)) // Gas Limit

	// Handler::Producer
	w.Use(commonHandlers.Producer(infra.Producer))

	return w
}
