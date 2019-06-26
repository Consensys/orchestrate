package steps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// GetChainCounts returns a mapping that counts the number of tx per chain
func GetChainCounts(envelopes map[string]*envelope.Envelope) map[string]uint {
	chains := make(map[string]uint)
	for _, v := range envelopes {
		chains[v.GetChain().ID().String()]++
	}
	return chains
}

// ChanTimeout constrains channel to receive message or to timeout
func ChanTimeout(c chan *envelope.Envelope, seconds int64, expectedItems int) ([]*envelope.Envelope, error) {
	envelopesChan := make([]*envelope.Envelope, expectedItems)
	for i := 0; i < expectedItems; i++ {
		select {
		case msg := <-c:
			envelopesChan[i] = msg
		case <-time.After(time.Duration(seconds) * time.Second):
			return nil, fmt.Errorf("timeout: not receiving msg after %d seconds", seconds)
		}
	}
	return envelopesChan, nil
}

// SendEnvelope sends an envelope to kafka
func SendEnvelope(e *envelope.Envelope) error {

	p := broker.GlobalSyncProducer()

	msg := &sarama.ProducerMessage{
		Topic: viper.GetString("kafka.topic.crafter"),
	}

	_ = encoding.Marshal(e, msg)

	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"kafka.out.partition": partition,
		"kafka.out.offset":    offset,
	}).Info("e2e: message sent")

	return nil
}

// EnvelopeCrafter crafts a envelope from a string mapping
func EnvelopeCrafter(m map[string]string) *envelope.Envelope {

	e := &envelope.Envelope{}

	for k, v := range m {
		switch k {
		case "chainId":
			chainID, err := strconv.Atoi(v)
			if err != nil {
				panic("Failed to parse chain id")
			}
			e.Chain = chain.CreateChainInt(int64(chainID))
		case "from":
			e.From = ethereum.HexToAccount(v)
		case "contractName":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetCall() == nil {
				e.Args.Call = &args.Call{}
			}
			if e.GetArgs().GetCall().GetContract() == nil {
				e.Args.Call.Contract = &abi.Contract{}
			}
			if e.GetArgs().GetCall().GetContract().GetId() == nil {
				e.Args.Call.Contract.Id = &abi.ContractId{}
			}

			e.GetArgs().GetCall().GetContract().GetId().Name = v
		case "contractTag":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetCall() == nil {
				e.Args.Call = &args.Call{}
			}
			if e.GetArgs().GetCall().GetContract() == nil {
				e.Args.Call.Contract = &abi.Contract{}
			}
			if e.GetArgs().GetCall().GetContract().GetId() == nil {
				e.Args.Call.Contract.Id = &abi.ContractId{}
			}

			e.GetArgs().GetCall().GetContract().GetId().Tag = v
		case "methodSignature":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetCall() == nil {
				e.Args.Call = &args.Call{}
			}
			if e.GetArgs().GetCall().GetMethod() == nil {
				e.Args.Call.Method = &abi.Method{}
			}

			e.GetArgs().GetCall().GetMethod().Signature = v
		case "args":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetCall() == nil {
				e.Args.Call = &args.Call{}
			}

			e.GetArgs().GetCall().Args = strings.Split(v, ",")
		case "to":
			if e.GetTx() == nil {
				e.Tx = &ethereum.Transaction{}
			}
			if e.GetTx().GetTxData() == nil {
				e.Tx.TxData = &ethereum.TxData{}
			}

			e.GetTx().GetTxData().To = ethereum.HexToAccount(v)
		case "gas":
			gas, _ := strconv.ParseUint(v, 10, 32)
			if e.GetTx() != nil {
				if e.GetTx().GetTxData() != nil {
					e.GetTx().GetTxData().SetGas(gas)
				} else {
					e.GetTx().TxData = &ethereum.TxData{
						Gas: gas,
					}
				}
			} else {
				e.Tx = &ethereum.Transaction{
					TxData: &ethereum.TxData{
						Gas: gas,
					},
				}
			}
		case "gasPrice":
			gasPrice, err := strconv.Atoi(v)
			if err != nil {
				panic("Failed to parse gas price")
			}
			if e.GetTx() == nil {
				e.Tx = &ethereum.Transaction{}
			}
			if e.GetTx().GetTxData() == nil {
				e.Tx.TxData = &ethereum.TxData{}
			}

			e.GetTx().GetTxData().GasPrice = ethereum.IntToQuantity(int64(gasPrice))
		case "value":
			value, err := strconv.Atoi(v)
			if err != nil {
				panic("Failed to parse value")
			}
			if e.GetTx() == nil {
				e.Tx = &ethereum.Transaction{}
			}
			if e.GetTx().GetTxData() == nil {
				e.Tx.TxData = &ethereum.TxData{}
			}

			e.GetTx().GetTxData().Value = ethereum.IntToQuantity(int64(value))
		case "metadataID":
			if e.GetMetadata() != nil {
				e.GetMetadata().Id = v
			} else {
				e.Metadata = &envelope.Metadata{Id: v}
			}
		default:
			if e.GetMetadata() != nil {
				if e.GetMetadata().GetExtra() != nil {
					e.GetMetadata().Extra[k] = v
				} else {
					extra := make(map[string]string)
					extra[k] = v
					e.GetMetadata().Extra = extra
				}
			} else {
				extra := make(map[string]string)
				extra[k] = v
				e.Metadata = &envelope.Metadata{Extra: extra}
			}

			e.Metadata.Extra[k] = v
		}
	}
	return e
}
