package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/workerpool"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/pkg/utils/envelope"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Shopify/sarama"
	encoding "github.com/consensys/orchestrate/pkg/encoding/sarama"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/ethereum/abi"
	sdk "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	ierror "github.com/consensys/orchestrate/pkg/types/error"
	types "github.com/consensys/orchestrate/pkg/types/ethereum"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"google.golang.org/protobuf/proto"
)

const component = "tx-listener.session.ethereum.hook"

type Hook struct {
	conf     *Config
	ec       ethclient.MultiClient
	producer sarama.SyncProducer
	client   sdk.OrchestrateClient
	logger   *log.Logger
}

func NewHook(
	conf *Config,
	ec ethclient.MultiClient,
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
				From:       utils.StringerToString(job.Transaction.From),
				Nonce:      utils.ValueToString(job.Transaction.Nonce),
				To:         utils.StringerToString(job.Transaction.To),
				Value:      utils.StringerToString(job.Transaction.Value),
				Gas:        utils.ValueToString(job.Transaction.Gas),
				GasPrice:   utils.StringerToString(job.Transaction.GasPrice),
				GasFeeCap:  utils.StringerToString(job.Transaction.GasFeeCap),
				GasTipCap:  utils.StringerToString(job.Transaction.GasTipCap),
				Data:       utils.StringerToString(job.Transaction.Data),
				Raw:        utils.StringerToString(job.Transaction.Raw),
				TxHash:     utils.StringerToString(job.Transaction.Hash),
				AccessList: envelope.ConvertFromAccessList(job.Transaction.AccessList),
				TxType:     string(job.Transaction.TransactionType),
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

		updateReq := &api.UpdateJobRequest{
			Status:  entities.StatusMined,
			Message: fmt.Sprintf("transaction mined in block %v", block.NumberU64()),
		}

		if txResponse.Receipt.EffectiveGasPrice != "" {
			effectiveGas, _ := hexutil.DecodeBig(txResponse.Receipt.EffectiveGasPrice)
			updateReq.Transaction = &entities.ETHTransaction{
				GasPrice: (*hexutil.Big)(effectiveGas),
			}
		}

		wp.Submit(func() {
			_, err := hk.client.UpdateJob(
				ctx,
				txResponse.GetJobUUID(),
				updateReq,
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
			return errors.DependencyFailureError("invalid receipt (no topics in log)")
		}

		logger := hk.logger.WithContext(ctx).WithField("sig_hash", utils.ShortString(l.Topics[0], 5)).
			WithField("address", l.GetAddress()).WithField("indexed", uint32(len(l.Topics)-1))

		logger.Debug("decoding receipt logs")
		eventResp, err := hk.client.GetContractEvents(
			ctx,
			l.GetAddress(),
			c.ChainID,
			&api.GetContractEventsRequest{
				SigHash:           hexutil.MustDecode(l.Topics[0]),
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
			logger.WithError(err).Trace("could not retrieve event ABI")
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
	if receipt.ContractAddress == "" || receipt.ContractAddress == "0x0000000000000000000000000000000000000000" {
		return nil
	}

	logger := hk.logger.WithContext(ctx).WithField("contract_address", receipt.ContractAddress)
	logger.Debug("register new deployed contract")
	var code []byte
	var err error
	if receipt.PrivacyGroupId != "" {
		// Fetch EEA deployed contract code
		code, err = hk.ec.PrivCodeAt(ctx, c.URL, ethcommon.HexToAddress(receipt.ContractAddress), receipt.PrivacyGroupId, block.Number())
	} else {
		code, err = hk.ec.CodeAt(ctx, c.URL, ethcommon.HexToAddress(receipt.ContractAddress), block.Number())
	}

	if err != nil {
		return err
	}

	err = hk.client.SetContractAddressCodeHash(ctx, receipt.ContractAddress, c.ChainID,
		&api.SetContractCodeHashRequest{
			CodeHash: crypto.Keccak256Hash(code).Bytes(),
		})
	if err != nil {
		logger.WithError(err).Error("failed to register contract")
		return err
	}

	logger.Info("contract has been registered successfully")
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
