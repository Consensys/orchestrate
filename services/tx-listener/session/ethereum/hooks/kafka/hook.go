package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	"github.com/gammazero/workerpool"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"

	encoding "github.com/ConsenSys/orchestrate/pkg/encoding/sarama"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/ethereum/abi"
	sdk "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	ierror "github.com/ConsenSys/orchestrate/pkg/types/error"
	types "github.com/ConsenSys/orchestrate/pkg/types/ethereum"
	"github.com/ConsenSys/orchestrate/pkg/types/tx"
	"github.com/ConsenSys/orchestrate/services/tx-listener/dynamic"
	"github.com/Shopify/sarama"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
)

const component = "tx-listener.session.ethereum.hook"

type Hook struct {
	conf     *Config
	ec       ethclient.ChainStateReader
	producer sarama.SyncProducer
	client   sdk.OrchestrateClient
	logger   *log.Logger
}

func NewHook(
	conf *Config,
	ec ethclient.ChainStateReader,
	producer sarama.SyncProducer,
	client sdk.OrchestrateClient,
) *Hook {
	return &Hook{
		conf:     conf,
		ec:       ec,
		producer: producer,
		client:   client,
		logger:   log.NewLogger().SetComponent(component),
	}
}

func (hk *Hook) AfterNewBlock(ctx context.Context, c *dynamic.Chain, block *ethtypes.Block, jobs []*entities.Job) error {
	blockLogCtx := log.WithFields(ctx, log.Field("chain", c.UUID), log.Field("block_number", block.Number().String()))
	logger := hk.logger.WithContext(blockLogCtx)

	var txResponses []*tx.TxResponse
	for _, job := range jobs {
		receiptLogCtx := log.WithFields(blockLogCtx, log.Field("receipt_tx_hash", job.Receipt.TxHash))
		// Register deployed contract
		err := hk.registerDeployedContract(receiptLogCtx, c, job.Receipt, block)
		if err != nil {
			hk.logger.WithContext(receiptLogCtx).WithError(err).Error("could not register deployed contract on registry")
		}

		txResponse := &tx.TxResponse{
			Id:            job.ScheduleUUID,
			JobUUID:       job.UUID,
			ContextLabels: job.Labels,
			Transaction: &types.Transaction{
				From:     job.Transaction.From,
				Nonce:    job.Transaction.Nonce,
				To:       job.Transaction.To,
				Value:    job.Transaction.Value,
				Gas:      job.Transaction.Gas,
				GasPrice: job.Transaction.GasPrice,
				Data:     job.Transaction.Data,
				Raw:      job.Transaction.Raw,
				TxHash:   job.Transaction.Hash,
			},
			Receipt: job.Receipt,
			Chain:   c.Name,
		}

		err = hk.decodeReceipt(receiptLogCtx, c, txResponse.Receipt)
		if err != nil {
			txResponse.Errors = []*ierror.Error{errors.FromError(err)}
		}

		txResponses = append(txResponses, txResponse)
	}

	// Update transactions to "MINED"
	// TODO: pass batch variable by environment variable
	wp := workerpool.New(20)
	for _, txResponse := range txResponses {
		if txResponse.GetJobUUID() == "" {
			continue
		}
		txResponse := txResponse
		wp.Submit(func() {
			_, err := hk.client.UpdateJob(
				ctx,
				txResponse.GetJobUUID(),
				&api.UpdateJobRequest{
					Status:  entities.StatusMined,
					Message: fmt.Sprintf("transaction mined in block %v", block.NumberU64()),
				},
			)
			if err != nil {
				logger.WithError(err).Warnf("failed to update status of %s to MINED", txResponse.Id)
			}
		})
	}
	wp.StopWait()

	// Prepare messages to be produced
	msgs, err := hk.prepareEnvelopeMsgs(txResponses, hk.conf.OutTopic, c.UUID)
	if err != nil {
		logger.WithError(err).Errorf("failed to prepare messages")
		return err
	}

	// Produce messages in Apache Kafka
	err = hk.produce(msgs)
	if err != nil {
		logger.WithError(err).Errorf("failed to produce message")
		return err
	}

	logger.Info("block processed")
	return nil
}

func (hk *Hook) decodeReceipt(ctx context.Context, c *dynamic.Chain, receipt *types.Receipt) error {
	hk.logger.WithContext(ctx).Debug("decoding receipt...")
	for _, l := range receipt.GetLogs() {
		if len(l.GetTopics()) == 0 {
			// This scenario is not supposed to happen
			return errors.InternalError("invalid receipt (no topics in log)")
		}

		logger := hk.logger.WithContext(ctx).WithField("sig_hash", utils.ShortString(l.Topics[0], 5)).
			WithField("address", l.GetAddress()).WithField("indexed", uint32(len(l.Topics)-1))

		logger.Debug("decoding receipt logs")
		eventResp, err := hk.client.GetContractEvents(
			ctx,
			l.GetAddress(),
			c.ChainID,
			&api.GetContractEventsRequest{
				SigHash:           l.Topics[0],
				IndexedInputCount: uint32(len(l.Topics) - 1),
			},
		)

		if err != nil {
			if errors.IsNotFoundError(err) {
				continue
			}

			logger.WithError(err).Error("failed to decode receipt logs")
			return err
		}

		if eventResp.Event == "" && len(eventResp.DefaultEvents) == 0 {
			logger.WithError(err).Warn("could not retrieve event ABI")
			continue
		}

		var mapping map[string]string
		event := &ethAbi.Event{}

		if eventResp.Event != "" {
			err = json.Unmarshal([]byte(eventResp.Event), event)
			if err != nil {
				logger.WithError(err).
					Warnf("could not unmarshal event ABI provided by the Contract Registry, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
				continue
			}
			mapping, err = abi.Decode(event, l)
		} else {
			for _, potentialEvent := range eventResp.DefaultEvents {
				// Try to unmarshal
				err = json.Unmarshal([]byte(potentialEvent), event)
				if err != nil {
					// If it fails to unmarshal, try the next potential event
					logger.WithError(err).Tracef("could not unmarshal potential event ABI, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
					continue
				}

				// Try to decode
				mapping, err = abi.Decode(event, l)
				if err == nil {
					// As the decoding is successful, stop looping
					break
				}
			}
		}

		if err != nil {
			// As all potentialEvents fail to unmarshal, go to the next log
			logger.WithError(err).Tracef("could not unmarshal potential event ABI, txHash: %s sigHash: %s, ", l.GetTxHash(), l.GetTopics()[0])
			continue
		}

		// Set decoded data on log
		l.DecodedData = mapping
		l.Event = GetAbi(event)

		logger.WithField("receipt_log", fmt.Sprintf("%v", mapping)).Debug("log decoded")
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

func (hk *Hook) registerDeployedContract(ctx context.Context, c *dynamic.Chain, receipt *types.Receipt, block *ethtypes.Block) error {
	if receipt.ContractAddress != "" && receipt.ContractAddress != "0x0000000000000000000000000000000000000000" {
		logger := hk.logger.WithContext(ctx).WithField("contract_address", receipt.ContractAddress)

		logger.Debug("register new deployed contract")
		code, err := hk.ec.CodeAt(ctx, c.URL, ethcommon.HexToAddress(receipt.ContractAddress), block.Number())
		if err != nil {
			return err
		}

		err = hk.client.SetContractAddressCodeHash(ctx, receipt.ContractAddress, c.ChainID,
			&api.SetContractCodeHashRequest{
				CodeHash: crypto.Keccak256Hash(code).String(),
			})
		if err != nil {
			logger.WithError(err).Error("failed to register contract")
			return err
		}
	}
	return nil
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
