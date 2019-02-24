package worker

import (
	"math/big"

	"github.com/spf13/viper"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/app/worker/handlers"
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
	w.Use(handCom.Producer(infra.Producer))

	return w
}
