package main

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	// Initialize client
	viper.Set("eth.client.url", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})
	ethclient.Init(context.Background())

	chain := big.NewInt(3)

	block, err := ethclient.GlobalClient().BlockByNumber(context.Background(), chain, nil)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "BlockByNumber"}).Fatal("Call failed")
	}

	log.WithFields(log.Fields{
		"chain.id":           chain.Text(16),
		"block.hash":         block.Hash().Hex(),
		"block.number":       block.Number().Text(10),
		"block.transactions": len(block.Transactions()),
	}).Infof("Latest block")

	blockHash := ethcommon.HexToHash("0x4d53ed90ecc4abeaca79840a1478ec011573a37347615b9a1bc69997806ce562")
	blockNumber := big.NewInt(5516994)

	header, err := ethclient.GlobalClient().HeaderByHash(context.Background(), chain, blockHash)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "HeaderByHash"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":        "HeaderByHash",
			"header.number": header.Number.Text(10),
		}).Info("Call succeeded")
	}

	header, err = ethclient.GlobalClient().HeaderByNumber(context.Background(), chain, blockNumber)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "HeaderByNumber"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":      "HeaderByNumber",
			"header.hash": header.Hash().Hex(),
		}).Info("Call succeeded")
	}

	block, err = ethclient.GlobalClient().BlockByNumber(context.Background(), chain, blockNumber)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "BlockByNumber"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":     "BlockByNumber",
			"block.hash": block.Hash().Hex(),
		}).Info("Call succeeded")
	}

	block, err = ethclient.GlobalClient().BlockByHash(context.Background(), chain, blockHash)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "BlockByHash"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":             "BlockByHash",
			"block.number":       block.Number().Text(10),
			"block.transactions": len(block.Transactions()),
		}).Info("Call succeeded")
	}

	txHash := ethcommon.HexToHash("0xdb695e527bb9c3e8ee2f607bf908dd98351e1bf4e1120c39df4ba435ca584aa5")
	tx, isPending, err := ethclient.GlobalClient().TransactionByHash(context.Background(), chain, txHash)
	from, _ := ethtypes.NewEIP155Signer(chain).Sender(tx)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "TransactionByHash"}).Errorf("Call failed")
	} else {

		log.WithFields(log.Fields{
			"method":              "TransactionByHash",
			"transaction.pending": isPending,
			"transaction.sender":  from.Hex(),
		}).Info("Call succeeded")
	}

	receipt, err := ethclient.GlobalClient().TransactionReceipt(context.Background(), chain, txHash)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "TransactionReceipt"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":         "TransactionReceipt",
			"receipt.status": receipt.Status,
			"receipt.logs":   len(receipt.Logs),
		}).Info("Call succeeded")
	}

	_, err = ethclient.GlobalClient().SyncProgress(context.Background(), chain)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "SyncProgress"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method": "SyncProgress",
		}).Info("Call succeeded")
	}

	balance, err := ethclient.GlobalClient().BalanceAt(context.Background(), chain, from, nil)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "BalanceAt"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":  "BalanceAt",
			"balance": balance.Text(10),
			"account": from.Hex(),
		}).Info("Call succeeded")
	}

	balance, err = ethclient.GlobalClient().PendingBalanceAt(context.Background(), chain, from)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "PendingBalanceAt"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":  "PendingBalanceAt",
			"balance": balance.Text(10),
			"account": from.Hex(),
		}).Info("Call succeeded")
	}

	nonce, err := ethclient.GlobalClient().NonceAt(context.Background(), chain, from, nil)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "NonceAt"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":  "NonceAt",
			"nonce":   nonce,
			"account": from.Hex(),
		}).Info("Call succeeded")
	}

	nonce, err = ethclient.GlobalClient().PendingNonceAt(context.Background(), chain, from)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "PendingNonceAt"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method":  "PendingNonceAt",
			"nonce":   nonce,
			"account": from.Hex(),
		}).Info("Call succeeded")
	}

	price, err := ethclient.GlobalClient().SuggestGasPrice(context.Background(), chain)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"method": "SuggestGasPrice"}).Errorf("Call failed")
	} else {
		log.WithFields(log.Fields{
			"method": "SuggestGasPrice",
			"chain":  chain.Text(10),
			"price":  price.Text(10),
		}).Info("Call succeeded")
	}
}
