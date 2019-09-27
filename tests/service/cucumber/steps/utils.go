package steps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/chain"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/faucet/types/ethereum"
)

// GetChainCounts returns a mapping that counts the number of tx per chain
func GetChainCounts(envelopes map[string]*envelope.Envelope) map[string]uint {
	chains := make(map[string]uint)
	for _, v := range envelopes {
		chains[v.GetChain().ID().String()]++
	}
	return chains
}

// ReadChanWithTimeout constrains channel to receive message or to timeout
func ReadChanWithTimeout(c chan *envelope.Envelope, seconds int64, expectedItems int) ([]*envelope.Envelope, error) {
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

// SendEnvelope sends the inputs Envelope to the provided kafka topic
func SendEnvelope(e *envelope.Envelope, topic string) error {

	p := broker.GlobalSyncProducer()

	msg := &sarama.ProducerMessage{
		Topic: topic,
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

// ParseEnvelope parses a string mapping into an envelope
func ParseEnvelope(m map[string]string) *envelope.Envelope {
	e := &envelope.Envelope{}
	for k, v := range m {
		switch k {
		case "chainId":
			chainID, err := strconv.Atoi(v)
			if err != nil {
				panic("Failed to parse chain id")
			}
			e.Chain = chain.FromInt(int64(chainID))
		case "from":
			e.From = ethereum.HexToAccount(v)
		case contractName:
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
		case contractTag:
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
		case "privateFrom":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetPrivate() == nil {
				e.GetArgs().Private = &args.Private{}
			}
			e.GetArgs().GetPrivate().PrivateFrom = v
		case "privateFor":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetPrivate() == nil {
				e.GetArgs().Private = &args.Private{}
			}
			e.GetArgs().GetPrivate().PrivateFor = strings.Split(v, ",")
		case "privateTxType":
			if e.GetArgs() == nil {
				e.Args = &envelope.Args{}
			}
			if e.GetArgs().GetPrivate() == nil {
				e.GetArgs().Private = &args.Private{}
			}
			e.GetArgs().GetPrivate().PrivateTxType = v
		case "protocol":
			protocolValue, err := strconv.Atoi(v)
			if err != nil {
				panic("Failed to parse value")
			}
			e.Protocol = &chain.Protocol{
				Type: chain.ProtocolType(
					int64(protocolValue),
				),
			}
		case "nonce":
			nonce, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				panic("Failed to parse nonce")
			}
			if e.GetTx() == nil {
				e.Tx = &ethereum.Transaction{}
			}
			if e.GetTx().GetTxData() == nil {
				e.Tx.TxData = &ethereum.TxData{}
			}

			e.GetTx().GetTxData().SetNonce(nonce)
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
