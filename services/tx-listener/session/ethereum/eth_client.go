package ethereum

import (
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
)

//go:generate mockgen -source=eth_client.go -destination=mocks/eth_client.go -package=mocks

type EthClient interface {
	ethclient.ChainStateReader
	ethclient.ChainLedgerReader
	ethclient.EEAChainStateReader
	ethclient.EEAChainLedgerReader
}
