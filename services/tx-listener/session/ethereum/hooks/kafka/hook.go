package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"github.com/Shopify/sarama"
	"github.com/containous/traefik/v2/pkg/log"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/common"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/error"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	svccontracts "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/dynamic"
)

type Hook struct {
	conf *Config

	registry          svccontracts.ContractRegistryClient
	ec                ethclient.ChainStateReader
	producer          sarama.SyncProducer
	txSchedulerClient client.TransactionSchedulerClient
}

func NewHook(
	conf *Config,
	registry svccontracts.ContractRegistryClient,
	ec ethclient.ChainStateReader,
	producer sarama.SyncProducer,
	txSchedulerClient client.TransactionSchedulerClient,
) *Hook {
	return &Hook{
		conf:              conf,
		registry:          registry,
		ec:                ec,
		producer:          producer,
		txSchedulerClient: txSchedulerClient,
	}
}

func (hk *Hook) AfterNewBlock(ctx context.Context, c *dynamic.Chain, block *ethtypes.Block, jobs []*entities.Job) error {
	blockLogCtx := log.With(ctx, log.Str("block.number", block.Number().String()))
	var txResponses []*tx.TxResponse

	for _, job := range jobs {
		receiptLogCtx := log.With(blockLogCtx, log.Str("receipt.txhash", job.Receipt.TxHash))

		// Register deployed contract
		err := hk.registerDeployedContract(receiptLogCtx, c, job.Receipt, block)
		if err != nil {
			log.FromContext(receiptLogCtx).WithError(err).Errorf("could not register deployed contract on registry")
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
	for _, txResponse := range txResponses {
		if txResponse.GetJobUUID() == "" {
			continue
		}

		_, err := hk.txSchedulerClient.UpdateJob(
			ctx,
			txResponse.GetJobUUID(),
			&txschedulertypes.UpdateJobRequest{
				Status:  utils.StatusMined,
				Message: fmt.Sprintf("Transaction mined in block %v", block.NumberU64()),
			})

		if err != nil {
			log.FromContext(blockLogCtx).WithError(err).Warnf("failed to update status of %s to MINED", txResponse.Id)
		}
	}

	// Prepare messages to be produced
	msgs, err := hk.prepareEnvelopeMsgs(txResponses, hk.conf.OutTopic, c.UUID)
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

	log.FromContext(blockLogCtx).Infof("block %v processed", block.NumberU64())

	return nil
}

func (hk *Hook) decodeReceipt(ctx context.Context, c *dynamic.Chain, receipt *types.Receipt) error {
	for _, l := range receipt.GetLogs() {
		if len(l.GetTopics()) == 0 {
			// This scenario is not supposed to happen
			return errors.InternalError("invalid receipt (no topics in log)")
		}

		// Retrieve event ABI from contract-registry
		eventResp, err := hk.registry.GetEventsBySigHash(
			ctx,
			&svccontracts.GetEventsBySigHashRequest{
				SigHash: l.Topics[0],
				AccountInstance: &common.AccountInstance{
					ChainId: c.ChainID,
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
			mapping, err = abi.Decode(event, l)
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
				mapping, err = abi.Decode(event, l)
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

func (hk *Hook) registerDeployedContract(ctx context.Context, c *dynamic.Chain, receipt *types.Receipt, block *ethtypes.Block) error {
	if receipt.ContractAddress != "" && receipt.ContractAddress != "0x0000000000000000000000000000000000000000" {
		log.FromContext(ctx).WithField("contract.address", receipt.ContractAddress).Infof("new contract deployed")
		code, err := hk.ec.CodeAt(ctx, c.URL, ethcommon.HexToAddress(receipt.ContractAddress), block.Number())
		if err != nil {
			return err
		}

		_, err = hk.registry.SetAccountCodeHash(ctx,
			&svccontracts.SetAccountCodeHashRequest{
				AccountInstance: &common.AccountInstance{
					ChainId: c.ChainID,
					Account: receipt.ContractAddress,
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
