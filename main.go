package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	commonHandlers "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	types "gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers"
)

var opts Config

type handler struct {
	w *core.Worker
}

func newSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	// Create producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	log.Info("Producer ready")
	return p
}

func prepareMsg(t *types.Trace, msg *sarama.ProducerMessage) error {
	marshaller := infSarama.NewMarshaller()

	err := marshaller.Marshal(t, msg)
	if err != nil {
		return err
	}

	// Set topic
	msg.Topic = opts.App.OutTopic
	return nil
}

func newEthClient(rawurl string) *infEth.EthClient {
	ec, err := infEth.Dial(rawurl)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to Ethereum client")
	return ec
}

// Setup configure handler
func (h *handler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = core.NewWorker(50)

	// Worker::unmarchaller
	h.w.Use(commonHandlers.Loader(infSarama.NewUnmarshaller()))

	// Worker::logger
	h.w.Use(hand.LoggerHandler)

	// Worker::marker
	h.w.Use(commonHandlers.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Worker::signer
	txSigner := infEth.NewStaticSigner(
		[]string{
			"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E", // 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 (faucet account)
			"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A", // "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"
			"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC", // "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"
			"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E", // "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
			"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6", // "0xffbBa394DEf3Ff1df0941c6429887107f58d4e9b"
			"C4B172E72033581BC41C36FA0448FCF031E9A31C4A3E300E541802DFB7248307", // 0x664895b5fE3ddf049d2Fb508cfA03923859763C6
			"706CC0876DA4D52B6DCE6F5A0FF210AEFCD51DE9F9CFE7D1BF7B385C82A06B8C", // 0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb
			"1476C66DE79A57E8AB4CADCECCBE858C99E5EDF3BFFEA5404B15322B5421E18C", // 0x93f7274c9059e601be4512F656B57b830e019E41
			"A2426FE76ECA2AA7852B95A2CE9CC5CC2BC6C05BB98FDA267F2849A7130CF50D", // 0xbfc7137876d7Ac275019d70434B0f0779824a969
			"41B9C5E497CFE6A1C641EFCA314FF84D22036D1480AF5EC54558A5EDD2FEAC03", // 0xA8d8DB1d8919665a18212374d623fc7C0dFDa410
		},
	)
	h.w.Use(
		hand.Signer(txSigner),
	)

	// Worker::producer
	h.w.Use(
		commonHandlers.Producer(
			infSarama.NewProducer(
				newSaramaSyncProducer(newSaramaClient([]string{opts.Conn.Kafka.URL})),
				prepareMsg,
			),
		),
	)

	return nil
}

// ConsumeClaim consume messages from queue
func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	in := make(chan interface{})
	go func() {
		// Pipe channels for interface compatibility
		for msg := range c.Messages() {
			in <- msg
		}
		close(in)
	}()
	h.w.Run(in)

	return nil
}

// Cleanup cleans handler
func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func newSaramaClient(kafkaURL []string) sarama.Client {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create client
	client, err := sarama.NewClient(kafkaURL, config)
	if err != nil {
		panic(err)
	}
	log.Info("Sarama client ready")
	return client
}

func main() {
	LoadConfig(&opts)
	ConfigureLogger(opts.Log)
	log.Info("Start worker...")

	client := newSaramaClient([]string{opts.Conn.Kafka.URL})

	// Create consumer
	g, err := sarama.NewConsumerGroupFromClient(opts.App.ConsumerGroup, client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Consumer Group ready")
	defer func() { g.Close() }()

	// Track errors
	go func() {
		for err := range g.Errors() {
			log.Error("ERROR", err)
		}
	}()

	g.Consume(context.Background(), []string{opts.App.InTopic}, &handler{})
}
