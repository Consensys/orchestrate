package infra

import (
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
)

func initGasManager(infra *Infra) {
	infra.GasManager = ethereum.NewGasManager(infra.Mec)
}
