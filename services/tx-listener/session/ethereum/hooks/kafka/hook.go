package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/containous/traefik/v2/pkg/log"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	ethclientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/utils"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

type Hook struct {
	conf *Config

	registry svc.ContractRegistryClient
	ec       ethclient.ChainStateReader

	store evlpstore.EnvelopeStoreClient

	producer sarama.SyncProducer
}

func NewHook(conf *Config, registry svc.ContractRegistryClient, ec ethclient.ChainStateReader, store evlpstore.EnvelopeStoreClient, producer sarama.SyncProducer) *Hook {
	return &Hook{
		conf:     conf,
		registry: registry,
		ec:       ec,
		store:    store,
		producer: producer,
	}
}

func (hk *Hook) AfterNewBlock(ctx context.Context, node *dynamic.Node, block *ethtypes.Block, receipts []*ethtypes.Receipt) error {
	blockLogCtx := log.With(ctx, log.Str("block.number", block.Number().String()))
	var evlps []*envelope.Envelope
	for idx, r := range receipts {
		receiptLogCtx := log.With(blockLogCtx, log.Str("receipt.txhash", r.TxHash.Hex()))

		// Register deployed contract
		err := hk.registerDeployedContract(receiptLogCtx, node, r, block)
		if err != nil {
			log.FromContext(receiptLogCtx).WithError(err).Errorf("could not register deployed contract on registry")
			return err
		}

		// Cast receipt in internal format
		receipt := ethereum.FromGethReceipt(r)

		// Load envelope from envelope store
		evlp, err := hk.loadEnvelope(receiptLogCtx, node, receipt)
		switch {
		case err == nil:
		case !hk.conf.DisableExternalTx && errors.IsNotFoundError(err):
			// External transaction
			evlp = &envelope.Envelope{
				Metadata: &envelope.Metadata{Id: uuid.NewV4().String()},
			}
		case errors.IsNotFoundError(err):
			// We could not found the envelope and external transaction are diasbled so we skip
			log.FromContext(receiptLogCtx).WithError(err).Debugf("skipping external receipt")
			continue
		default:
			log.FromContext(receiptLogCtx).WithError(err).Errorf("failed to load envelope")
			return err
		}

		// Attach receipt to envelope
		evlp.Receipt = receipt.
			SetBlockHash(block.Hash()).
			SetBlockNumber(block.NumberU64()).
			SetTxIndex(uint64(idx))

		evlps = append(evlps, evlp)
	}

	// Prepare messages to be produced
	msgs, err := hk.prepareEnvelopeMsgs(evlps, hk.conf.TopicTxDecoder, node.ID)
	if err != nil {
		log.FromContext(blockLogCtx).WithError(err).Errorf("failed to prepare messages")
		return err
	}

	// Produce messages in Apache Kafka
	err = hk.produce(msgs)
	if err != nil {
		log.FromContext(blockLogCtx).WithError(err).Errorf("failed to produce message")
		return err
	}

	log.FromContext(blockLogCtx).Infof("block processed")

	return nil
}

func (hk *Hook) registerDeployedContract(ctx context.Context, node *dynamic.Node, receipt *ethtypes.Receipt, block *ethtypes.Block) error {
	if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
		log.FromContext(ctx).WithField("contract.address", receipt.ContractAddress.Hex()).Infof("new contract deployed")
		code, err := hk.ec.CodeAt(
			ethclientutils.RetryNotFoundError(ctx, true),
			node.URL,
			receipt.ContractAddress,
			block.Number(),
		)
		if err != nil {
			return err
		}

		_, err = hk.registry.SetAccountCodeHash(ctx,
			&svc.SetAccountCodeHashRequest{
				AccountInstance: &common.AccountInstance{
					Chain:   chain.FromBigInt(node.ChainID),
					Account: ethereum.HexToAccount(receipt.ContractAddress.Hex()),
				},
				CodeHash: crypto.Keccak256Hash(code).Bytes(),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hk *Hook) loadEnvelope(ctx context.Context, node *dynamic.Node, receipt *ethereum.Receipt) (*envelope.Envelope, error) {
	resp, err := hk.store.LoadByTxHash(
		ctx,
		&evlpstore.LoadByTxHashRequest{
			Chain:  chain.FromBigInt(node.ChainID),
			TxHash: receipt.GetTxHash(),
		})
	if err != nil {
		return nil, err
	}
	return resp.GetEnvelope(), nil
}

func (hk *Hook) prepareEnvelopeMsgs(evlps []*envelope.Envelope, topic, key string) ([]*sarama.ProducerMessage, error) {
	var msgs []*sarama.ProducerMessage
	for _, e := range evlps {
		msg, err := hk.prepareMsg(e, topic, key)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (hk *Hook) prepareMsg(pb proto.Message, topic, key string) (*sarama.ProducerMessage, error) {
	msg := &sarama.ProducerMessage{}

	err := encoding.Marshal(pb, msg)
	if err != nil {
		return nil, err
	}

	// Set topic to TxDecoder
	msg.Topic = topic

	// Set Message key to node ID
	msg.Key = sarama.StringEncoder(key)

	return msg, nil
}

func (hk *Hook) produce(msgs []*sarama.ProducerMessage) error {
	return hk.producer.SendMessages(msgs)
}
