package base

import (
	"context"
	"math/big"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/handler/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
)

// Handler implements TxListenerHandler interface
//
// It uses a pkg Engine to listen to chains messages
type Handler struct {
	base.Handler

	client   sarama.Client
	producer sarama.SyncProducer
}

// NewHandler creates a new EngineConsumerGroupHandler
func NewHandler(e *engine.Engine, client sarama.Client, producer sarama.SyncProducer, conf *base.Config) *Handler {
	return &Handler{
		Handler:  *(base.NewHandler(e, conf)),
		client:   client,
		producer: producer,
	}
}

func (h *Handler) GetInitialPosition(chain *big.Int) (blockNumber, txIndex int64, err error) {
	position, ok := h.Conf.Start.Positions[chain.Text(10)]
	if !ok {
		blockNumber = h.Conf.Start.Default.BlockNumber
		txIndex = h.Conf.Start.Default.TxIndex
	} else {
		blockNumber = position.BlockNumber
		txIndex = position.TxIndex
	}

	if blockNumber >= -1 {
		return blockNumber, txIndex, nil
	}

	// BlockNumber == -2 means we should start listening from position of the last produced message
	decoderTopic := utils.KafkaChainTopic(viper.GetString("kafka.topic.decoder"), chain)

	// Retrieve last record
	lastRecord, err := h.getLastRecord(decoderTopic, 0)
	if err != nil || lastRecord == nil {
		// If we have never produced then we start from latest
		return -1, 0, nil
	}

	// Parse last record into envelope
	e := &envelope.Envelope{}
	err = proto.Unmarshal(lastRecord.Value, e)
	if err != nil {
		return -1, 0, err
	}
	return int64(e.GetReceipt().GetBlockNumber()), int64(e.GetReceipt().GetTxIndex() + 1), nil
}

func (h *Handler) getLastRecord(topic string, partition int32) (*sarama.Record, error) {
	// Retrieve last offset that has been produced for topic-partition
	lastOffset, err := h.client.GetOffset(topic, partition, -1)
	if err != nil {
		return nil, err
	}

	// Get broker Leader fo topic-partition
	broker, err := h.client.Leader(topic, partition)
	if err != nil {
		return nil, err
	}

	// Fetch block containing last produced record on topic partition
	req := &sarama.FetchRequest{
		MinBytes:    h.client.Config().Consumer.Fetch.Min,
		MaxWaitTime: int32(h.client.Config().Consumer.MaxWaitTime / time.Millisecond),
	}
	req.AddBlock(topic, 0, lastOffset-1, h.client.Config().Consumer.Fetch.Max)
	req.Version = 4
	req.Isolation = sarama.ReadUncommitted
	response, err := broker.Fetch(req)

	if err != nil {
		return nil, err
	}

	// Parse block to retrieve record of interest
	block := response.GetBlock(topic, partition)
	if len(block.RecordsSet) == 0 {
		return nil, nil
	}
	records := block.RecordsSet[0]
	record := records.RecordBatch.Records[0]

	return record, nil
}

// Pipe take a channel of types.TxListenerReceipt and pipes it into a channel of interface{}
//
// Pipe will stop forwarding messages either
// - receipt channel is closed
// - ctx has been canceled
func Pipe(ctx context.Context, receiptChan <-chan *types.TxListenerReceipt) <-chan interface{} {
	interfaceChan := make(chan interface{})

	// Start a goroutine that pipe messages
	go func() {
	pipeLoop:
		for {
			select {
			case msg, ok := <-receiptChan:
				if !ok {
					// Sarama channel has been closed so we exit loop
					break pipeLoop
				}
				interfaceChan <- msg
			case <-ctx.Done():
				// Context has been cancel so we exit loop
				break pipeLoop
			}
		}
		close(interfaceChan)
	}()

	return interfaceChan
}
