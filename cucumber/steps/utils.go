package steps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	ethCommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

func ChanTimeout(c chan *envelope.Envelope, seconds int64, expectedItems int) ([]*envelope.Envelope, error) {
	ch := make([]*envelope.Envelope, expectedItems)
	for i := 0; i < expectedItems; i++ {
		select {
		case msg := <-c:
			ch[i] = msg
		case <-time.After(time.Duration(seconds) * time.Second):
			return nil, fmt.Errorf("error not receiving msg")
		}
	}
	return ch, nil
}

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

func EnvelopeCrafter(m map[string]string) *envelope.Envelope {

	e := &envelope.Envelope{}

	for k, v := range m {
		switch k {
		case "chainId":
			if e.GetChain() != nil {
				e.GetChain().Id = v
			} else {
				e.Chain = &common.Chain{Id: v}
			}
		case "from":
			if e.GetSender() != nil {
				e.GetSender().Addr = v
			} else {
				e.Sender = &common.Account{Addr: v}
			}
		case "contractName":
			if e.GetCall() != nil {
				if e.GetCall().GetContract() != nil {
					e.GetCall().GetContract().Name = v
				} else {
					e.GetCall().Contract = &abi.Contract{Name: v}
				}
			} else {
				e.Call = &common.Call{
					Contract: &abi.Contract{
						Name: v,
					},
				}
			}
		case "contractTag":
			if e.GetCall() != nil {
				if e.GetCall().GetContract() != nil {
					e.GetCall().GetContract().Tag = v
				} else {
					e.GetCall().Contract = &abi.Contract{Tag: v}
				}
			} else {
				e.Call = &common.Call{
					Contract: &abi.Contract{
						Tag: v,
					},
				}
			}
		case "method":
			if e.GetCall() != nil {
				if e.GetCall().GetMethod() != nil {
					e.GetCall().GetMethod().Signature = v
				} else {
					e.GetCall().Method = &abi.Method{Signature: v}
				}
			} else {
				e.Call = &common.Call{
					Method: &abi.Method{
						Signature: v,
					},
				}
			}
		case "args":
			if e.GetCall() != nil {
				e.GetCall().Args = strings.Split(v, ":")
			} else {
				e.Call = &common.Call{
					Args: strings.Split(v, ":"),
				}
			}
		case "to":
			if e.GetTx() != nil {
				if e.GetTx().GetTxData() != nil {
					// Todo replace interface instead of Address
					e.GetTx().GetTxData().SetTo(ethCommon.HexToAddress(v))
				} else {
					e.GetTx().TxData = &ethereum.TxData{
						To: v,
					}
				}
			} else {
				e.Tx = &ethereum.Transaction{
					TxData: &ethereum.TxData{
						To: v,
					},
				}
			}
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
			if e.GetTx() != nil {
				if e.GetTx().GetTxData() != nil {
					e.GetTx().GetTxData().GasPrice = v
				} else {
					e.GetTx().TxData = &ethereum.TxData{
						GasPrice: v,
					}
				}
			} else {
				e.Tx = &ethereum.Transaction{
					TxData: &ethereum.TxData{
						GasPrice: v,
					},
				}
			}
		case "value":
			if e.GetTx() != nil {
				if e.GetTx().GetTxData() != nil {
					e.GetTx().GetTxData().Value = v
				} else {
					e.GetTx().TxData = &ethereum.TxData{
						Value: v,
					}
				}
			} else {
				e.Tx = &ethereum.Transaction{
					TxData: &ethereum.TxData{
						Value: v,
					},
				}
			}
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
