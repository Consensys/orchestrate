package main

func main() {}

// import (
// 	// "os"
// 	// "os/signal"
// 	// "syscall"
// 	"context"
// 	"fmt"
// 	"math/big"
// 	"time"

// 	"github.com/Shopify/sarama"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/golang/protobuf/proto"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/handlers"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/infra"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf"
// 	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
// )

// var (
// 	chainURL      = "https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7"
// 	kafkaURL      = []string{"localhost:9092"}
// 	group         = "test-group"
// 	inTopic       = "test-in"
// 	outTopic      = "test-out"
// 	faucetAddress = common.HexToAddress("0x7E654d251Da770A068413677967F6d3Ea2FeA9E4")
// )

// // ERC20TransferRegistry holds ERC20 ABI
// var ERC20TransferRegistry = infra.NewDummyABIRegistry(
// 	[]byte(`{
// 		"constant": false,
// 		"inputs": [
// 			{
// 				"name": "_to",
// 				"type": "address"
// 			},
// 			{
// 				"name": "_value",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "transfer",
// 		"outputs": [
// 			{
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"payable": false,
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	}`),
// )

// func newEthClient(rawurl string) *infra.EthClient {
// 	ec, err := infra.Dial(rawurl)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Connected to Ethereum client")
// 	return ec
// }

// func newSaramaClient(kafkaURL []string) sarama.Client {
// 	config := sarama.NewConfig()
// 	config.Version = sarama.V1_0_0_0
// 	config.Consumer.Return.Errors = true
// 	config.Producer.Return.Errors = true
// 	config.Producer.Return.Successes = true

// 	// Create client
// 	client, err := sarama.NewClient(kafkaURL, config)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Sarama client ready")
// 	return client
// }

// func newSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
// 	// Create producer
// 	p, err := sarama.NewSyncProducerFromClient(client)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Producer ready")
// 	return p
// }

// // SaramaHandler is a sarama ConsumerGroupHandler
// type SaramaHandler struct {
// 	w *types.Worker
// }

// // NewSaramaHandler creates a new handler
// func NewSaramaHandler() *SaramaHandler {
// 	return &SaramaHandler{}
// }

// func logger(ctx *types.Context) {
// 	msg := ctx.Msg.(*sarama.ConsumerMessage)
// 	fmt.Printf("Entry %v\n", msg.Offset)

// 	ctx.Next()

// 	fmt.Printf(
// 		"Output %v:\nFrom: %v\nTo: %v\nNonce: %v\nValue: %v\nGasLimit: %v\nGasPrice: %v\nData: %q\nRaw: %v\nHash: %v\nErrors: %v\n\n",
// 		ctx.Msg.(*sarama.ConsumerMessage).Offset,
// 		ctx.T.Sender().Address.Hex(),
// 		ctx.T.Tx().To().Hex(),
// 		ctx.T.Tx().Nonce(),
// 		hexutil.EncodeBig(ctx.T.Tx().Value()),
// 		ctx.T.Tx().GasLimit(),
// 		hexutil.EncodeBig(ctx.T.Tx().GasPrice()),
// 		hexutil.Encode(ctx.T.Tx().Data()),
// 		hexutil.Encode(ctx.T.Tx().Raw()),
// 		ctx.T.Tx().Hash().Hex(),
// 		ctx.T.Errors,
// 	)
// }

// func ctxToProducerMessage(pb *tracepb.Trace) *sarama.ProducerMessage {
// 	b, _ := proto.Marshal(pb)
// 	msg := sarama.ProducerMessage{}
// 	msg.Value = sarama.ByteEncoder(b)
// 	msg.Topic = outTopic
// 	return &msg
// }

// func makeFaucetMessage(chainID *big.Int, a common.Address, value *big.Int) *sarama.ProducerMessage {
// 	// Create a trace for
// 	t := types.NewTrace()
// 	*t.Chain().ID = *chainID
// 	*t.Sender().Address = faucetAddress
// 	t.Tx().SetValue(value)
// 	t.Tx().SetTo(&a)

// 	pb := &tracepb.Trace{}
// 	protobuf.DumpTrace(t, pb)
// 	b, _ := proto.Marshal(pb)

// 	msg := sarama.ProducerMessage{}
// 	msg.Value = sarama.ByteEncoder(b)
// 	msg.Topic = inTopic // re-enter message in input queue

// 	return &msg
// }

// // Setup configure handler
// func (h *SaramaHandler) Setup(s sarama.ConsumerGroupSession) error {
// 	// Create worker
// 	h.w = types.NewWorker(50)

// 	// Fake logger
// 	h.w.Use(logger)

// 	// Sarama message loader
// 	h.w.Use(handlers.Loader(&infra.SaramaUnmarshaller{}))

// 	// Crafter
// 	crafter := infra.PayloadCrafter{}
// 	h.w.Use(handlers.Crafter(ERC20TransferRegistry, &crafter))

// 	// Gas Price
// 	h.w.Use(
// 		handlers.GasPricer(
// 			infra.NewSimpleGasManager(newEthClient(chainURL)),
// 		),
// 	)

// 	// Gas Limit
// 	h.w.Use(
// 		handlers.GasEstimator(
// 			infra.NewSimpleGasManager(newEthClient(chainURL)),
// 		),
// 	)

// 	// Faucet
// 	cfg := &infra.SimpleCreditControllerConfig{
// 		BalanceAt:    infra.NewEthBalanceAt(newEthClient(chainURL)),
// 		CreditAmount: big.NewInt(100000000000000000), // 0.1 ETH
// 		MaxBalance:   big.NewInt(200000000000000000), // 0.2 ETH
// 		CreditDelay:  time.Duration(60 * time.Second),
// 		BlackList:    map[string]struct{}{faucetAddress.Hex(): struct{}{}},
// 	}
// 	crediter := infra.NewSaramaCrediter(
// 		newSaramaSyncProducer(newSaramaClient(kafkaURL)),
// 		makeFaucetMessage,
// 	)
// 	controller := infra.NewSimpleCreditController(cfg, 50)
// 	h.w.Use(handlers.Faucet(crediter, controller))

// 	// Nonce
// 	h.w.Use(
// 		handlers.NonceHandler(
// 			infra.NewCacheNonceManager(
// 				infra.NewEthClientNonceCalibrate(newEthClient(chainURL)),
// 				40,
// 			),
// 		),
// 	)

// 	// Signer
// 	txSigner := infra.NewStaticSigner(
// 		[]string{
// 			"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E", // 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 (faucet account)
// 			"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A", // "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"
// 			"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC", // "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"
// 			"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E", // "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
// 			"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6", // "0xffbBa394DEf3Ff1df0941c6429887107f58d4e9b"
// 			"C4B172E72033581BC41C36FA0448FCF031E9A31C4A3E300E541802DFB7248307", // 0x664895b5fE3ddf049d2Fb508cfA03923859763C6
// 			"706CC0876DA4D52B6DCE6F5A0FF210AEFCD51DE9F9CFE7D1BF7B385C82A06B8C", // 0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb
// 			"1476C66DE79A57E8AB4CADCECCBE858C99E5EDF3BFFEA5404B15322B5421E18C", // 0x93f7274c9059e601be4512F656B57b830e019E41
// 			"A2426FE76ECA2AA7852B95A2CE9CC5CC2BC6C05BB98FDA267F2849A7130CF50D", // 0xbfc7137876d7Ac275019d70434B0f0779824a969
// 			"41B9C5E497CFE6A1C641EFCA314FF84D22036D1480AF5EC54558A5EDD2FEAC03", // 0xA8d8DB1d8919665a18212374d623fc7C0dFDa410
// 		},
// 	)
// 	h.w.Use(handlers.Signer(txSigner))

// 	// Sender
// 	h.w.Use(
// 		handlers.Sender(
// 			infra.NewSimpleSender(newEthClient(chainURL)),
// 		),
// 	)

// 	// Producer
// 	h.w.Use(
// 		handlers.Producer(
// 			infra.NewSaramaProducer(
// 				newSaramaSyncProducer(newSaramaClient(kafkaURL)),
// 				ctxToProducerMessage,
// 			),
// 		),
// 	)

// 	// Marker
// 	h.w.Use(handlers.Marker(infra.NewSimpleSaramaOffsetMarker(s)))

// 	return nil
// }

// // Cleanup cleans handler
// func (h *SaramaHandler) Cleanup(s sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // ConsumeClaim consume messages from queue
// func (h *SaramaHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
// 	in := make(chan interface{})
// 	go func() {
// 		// Pipe channels for interface compatibility
// 		for msg := range c.Messages() {
// 			in <- msg
// 		}
// 		close(in)
// 	}()
// 	h.w.Run(in)

// 	return nil
// }

// func main() {
// 	// Create client
// 	client := newSaramaClient(kafkaURL)

// 	// Create consumer
// 	g, err := sarama.NewConsumerGroupFromClient(group, client)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Consumer Group ready")
// 	defer func() { g.Close() }()

// 	// Track errors
// 	go func() {
// 		for err := range g.Errors() {
// 			fmt.Println("ERROR", err)
// 		}
// 	}()

// 	// // InitSignals redirect signals
// 	// signals = make(chan os.Signal 3)
// 	// signal.Notify(signals)

// 	// // ProcessSignal process signals
// 	// func (w *SaramaWorker) ProcessSignal(signal os.Signal) {
// 	// 	switch signal {
// 	// 	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
// 	// 		// Gracefully stop
// 	// 		fmt.Printf("Worker: gracefully stops...\n")
// 	// 		w.Stop()
// 	// 	default:
// 	// 		// Exit
// 	// 		fmt.Printf("Worker: unknown signal exits...\n")
// 	// 		w.Exit(1)
// 	// 	}
// 	// }

// 	//ctx, _ := context.WithCancel(context.Background())
// 	g.Consume(context.Background(), []string{inTopic}, NewSaramaHandler())
// }
