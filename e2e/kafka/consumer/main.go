package main

import (
	// "os"
	// "os/signal"
	// "syscall"
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
)

func erc20Getter() handlers.ABIGetter {
	var method abi.Method
	json.Unmarshal(
		[]byte(`{
			"constant": false,
			"inputs": [
				{
					"name": "_to",
					"type": "address"
				},
				{
					"name": "_value",
					"type": "uint256"
				}
			],
			"name": "transfer",
			"outputs": [
				{
					"name": "",
					"type": "bool"
				}
			],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		}`),
		&method,
	)
	return handlers.NewDummyABIGetter(&method)
}

func newNonceEthClient(ec *ethclient.Client) handlers.NewNonceFunc {
	return func(chainID *big.Int, a common.Address) (uint64, error) {
		fmt.Printf("Getting nonce for %v", a.Hex())
		return ec.PendingNonceAt(context.Background(), a)
	}
}

// SaramaHandler is a sarama ConsumerGroupHandler
type SaramaHandler struct {
	w *infra.Worker
}

// NewSaramaHandler creates a new handler
func NewSaramaHandler() *SaramaHandler {
	return &SaramaHandler{}
}

func logger(ctx *infra.Context) {
	msg := ctx.Msg.(*sarama.ConsumerMessage)
	fmt.Printf("Entry %v\n%v\n", msg.Offset, string(msg.Value))

	ctx.Next()

	fmt.Printf(
		"Output:\nNonce: %v\nData: %q\nErrors: %v\n",
		ctx.T.Tx().Nonce(),
		hexutil.Encode(ctx.T.Tx().Data()),
		ctx.T.Errors,
	)
}

// Setup configure handler
func (h *SaramaHandler) Setup(s sarama.ConsumerGroupSession) error {
	// Create worker
	h.w = infra.NewWorker(50)

	// Fake logger
	h.w.Use(logger)

	// Sarama message loader
	h.w.Use(handlers.Loader(&handlers.SaramaUnmarshaller{}))

	// Crafer
	h.w.Use(handlers.Crafter(erc20Getter()))

	// Nonce
	chainURL := "https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"
	ec, err := ethclient.Dial(chainURL)
	if err != nil {
		fmt.Printf("Could not connect to eth client: %v\n", err)
		return err
	}
	fmt.Printf("Coonected to Ethereum client")
	h.w.Use(handlers.NonceHandler(handlers.NewCacheNonce(newNonceEthClient(ec), 20)))

	// Marker
	h.w.Use(handlers.Marker(handlers.NewSimpleSaramaOffsetMarker(s)))

	return nil
}

// Cleanup cleans handler
func (h *SaramaHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim consume messages from queue
func (h *SaramaHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	in := make(chan interface{})
	go func() {
		// Pipe for type compatibility
		for msg := range c.Messages() {
			in <- msg
		}
		close(in)
	}()
	h.w.Run(in)

	return nil
}

func main() {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true

	// Create client
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { client.Close() }()
	fmt.Println("Client ready")

	// Create consumer
	g, err := sarama.NewConsumerGroupFromClient("test-group", client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Consumer Group ready")
	defer func() { g.Close() }()

	// Track errors
	go func() {
		for err := range g.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// // InitSignals redirect signals
	// signals = make(chan os.Signal 3)
	// signal.Notify(signals)

	// // ProcessSignal process signals
	// func (w *SaramaWorker) ProcessSignal(signal os.Signal) {
	// 	switch signal {
	// 	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
	// 		// Gracefully stop
	// 		fmt.Printf("Worker: gracefully stops...\n")
	// 		w.Stop()
	// 	default:
	// 		// Exit
	// 		fmt.Printf("Worker: unknown signal exits...\n")
	// 		w.Exit(1)
	// 	}
	// }

	//ctx, _ := context.WithCancel(context.Background())
	g.Consume(context.Background(), []string{"test"}, NewSaramaHandler())
}
