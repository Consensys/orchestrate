package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/abi/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/common"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
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

func (hk *Hook) AfterNewBlock(ctx context.Context, c *dynamic.Chain, block *ethtypes.Block, receipts []*ethtypes.Receipt) error {
	blockLogCtx := log.With(ctx, log.Str("block.number", block.Number().String()))
	var evlps []*tx.TxResponse

	for idx, r := range receipts {
		b := tx.NewBuilder()

		receiptLogCtx := log.With(blockLogCtx, log.Str("receipt.txhash", r.TxHash.Hex()))

		// Register deployed contract
		err := hk.registerDeployedContract(receiptLogCtx, c, r, block)
		if err != nil {
			log.FromContext(receiptLogCtx).WithError(err).Errorf("could not register deployed contract on registry")
			return err
		}

		// Cast receipt in internal format
		receipt := ethereum.FromGethReceipt(r)

		// Load envelope from envelope store
		req, err := hk.loadEnvelope(receiptLogCtx, c, receipt)
		isExternalTx := errors.IsNotFoundError(err)
		if req != nil {
			b, err = req.Builder()
			if err != nil {
				log.FromContext(receiptLogCtx).WithError(err).Errorf("loaded invalid envelope - id: %s", req.GetID())
				return err
			}
		}
		switch {
		case err == nil:
		case c.Listener.ExternalTxEnabled && isExternalTx:
			// External transaction that we listen to, we create an envelope for it
			envelopeUUID := uuid.NewV4().String()
			log.FromContext(receiptLogCtx).WithField("uuid", envelopeUUID).Debugf("External transaction received")
			_ = b.SetID(envelopeUUID)
		case !c.Listener.ExternalTxEnabled && isExternalTx:
			// External transaction that we skip
			log.FromContext(receiptLogCtx).WithError(err).Debugf("Skipping external transaction")
			continue
		default:
			log.FromContext(receiptLogCtx).WithError(err).Errorf("Failed to load envelope")
			return err
		}

		// Attach receipt to envelope
		_ = b.SetReceipt(receipt.
			SetBlockHash(block.Hash()).
			SetBlockNumber(block.NumberU64()).
			SetTxIndex(uint64(idx)))

		err = hk.decodeReceipt(receiptLogCtx, c, b.Receipt)
		if err != nil {
			b.Errors = append(b.Errors, errors.FromError(err))
		}

		evlps = append(evlps, b.TxResponse())
	}

	// Prepare messages to be produced
	msgs, err := hk.prepareEnvelopeMsgs(evlps, hk.conf.OutTopic, c.UUID)
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

func (hk *Hook) decodeReceipt(ctx context.Context, c *dynamic.Chain, receipt *ethereum.Receipt) error {
	for _, l := range receipt.GetLogs() {
		if len(l.GetTopics()) == 0 {
			// This scenario is not supposed to happen
			return errors.InternalError("invalid receipt (no topics in log)")
		}

		// Retrieve event ABI from contract-registry
		eventResp, err := hk.registry.GetEventsBySigHash(
			ctx,
			&svc.GetEventsBySigHashRequest{
				SigHash: l.Topics[0],
				AccountInstance: &common.AccountInstance{
					ChainId: c.ChainID.String(),
					Account: l.GetAddress(),
				},
				IndexedInputCount: uint32(len(l.Topics) - 1),
			},
		)
		if err != nil || (eventResp.GetEvent() == "" && len(eventResp.GetDefaultEvents()) == 0) {
			log.FromContext(ctx).WithError(err).Tracef("could not retrieve event ABI, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
			continue
		}

		// Decode log
		var mapping map[string]string
		event := &ethAbi.Event{}

		if eventResp.GetEvent() != "" {
			err = json.Unmarshal([]byte(eventResp.GetEvent()), event)
			if err != nil {
				log.FromContext(ctx).WithError(err).Warnf("could not unmarshal event ABI provided by the Contract Registry, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
				continue
			}
			mapping, err = decoder.Decode(event, l)
		} else {
			for _, potentialEvent := range eventResp.GetDefaultEvents() {
				// Try to unmarshal
				err = json.Unmarshal([]byte(potentialEvent), event)
				if err != nil {
					// If it fails to unmarshal, try the next potential event
					log.FromContext(ctx).WithError(err).Tracef("could not unmarshal potential event ABI, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
					continue
				}

				// Try to decode
				mapping, err = decoder.Decode(event, l)
				if err == nil {
					// As the decoding is successful, stop looping
					break
				}
			}
		}

		if err != nil {
			// As all potentialEvents fail to unmarshal, go to the next log
			log.FromContext(ctx).WithError(err).Tracef("could not unmarshal potential event ABI, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
			continue
		}

		// Set decoded data on log
		l.DecodedData = mapping
		l.Event = GetAbi(event)

		receiptLogCtx := log.With(ctx, log.Str("receipt.log", fmt.Sprintf("%v", mapping)))
		log.FromContext(receiptLogCtx).Debug("decoder: log decoded")
	}
	return nil
}

// GetAbi creates a string ABI (format EventName(argType1, argType2)) from an event
func GetAbi(e *ethAbi.Event) string {
	inputs := make([]string, len(e.Inputs))
	for i := range e.Inputs {
		inputs[i] = fmt.Sprintf("%v", e.Inputs[i].Type)
	}
	return fmt.Sprintf("%v(%v)", e.Name, strings.Join(inputs, ","))
}

func (hk *Hook) registerDeployedContract(ctx context.Context, c *dynamic.Chain, receipt *ethtypes.Receipt, block *ethtypes.Block) error {
	if receipt.ContractAddress.Hex() != "0x0000000000000000000000000000000000000000" {
		log.FromContext(ctx).WithField("contract.address", receipt.ContractAddress.Hex()).Infof("new contract deployed")
		code, err := hk.ec.CodeAt(
			ethclientutils.RetryNotFoundError(ctx, true),
			c.URL,
			receipt.ContractAddress,
			block.Number(),
		)
		if err != nil {
			return err
		}

		_, err = hk.registry.SetAccountCodeHash(ctx,
			&svc.SetAccountCodeHashRequest{
				AccountInstance: &common.AccountInstance{
					ChainId: c.ChainID.String(),
					Account: receipt.ContractAddress.Hex(),
				},
				CodeHash: crypto.Keccak256Hash(code).String(),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hk *Hook) loadEnvelope(ctx context.Context, c *dynamic.Chain, receipt *ethereum.Receipt) (*tx.TxEnvelope, error) {
	ctx = multitenancy.WithTenantID(ctx, c.TenantID)
	resp, err := hk.store.LoadByTxHash(
		ctx,
		&evlpstore.LoadByTxHashRequest{
			ChainId: c.ChainID.String(),
			TxHash:  receipt.GetTxHash(),
		})
	if err != nil {
		return nil, err
	}
	return resp.GetEnvelope(), nil
}

func (hk *Hook) prepareEnvelopeMsgs(evlps []*tx.TxResponse, topic, key string) ([]*sarama.ProducerMessage, error) {
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

	// Set Message key to chain UUID
	msg.Key = sarama.StringEncoder(key)

	return msg, nil
}

func (hk *Hook) produce(msgs []*sarama.ProducerMessage) error {
	return hk.producer.SendMessages(msgs)
}
