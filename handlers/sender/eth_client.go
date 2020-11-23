package sender

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
)

//go:generate mockgen -source=eth_client.go -destination=mocks/eth_client.go -package=mocks

type EthClient interface {
	ethclient.TransactionSender
	ethclient.QuorumTransactionSender
	ethclient.EEATransactionSender
}
