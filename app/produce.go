package app

import (
	"context"
	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/rpc"
)

// NewMessage creates a New Sarama producer message
func NewMessage(i int) *sarama.ProducerMessage {

	call := &common.Call{
		Contract: &abi.Contract{Name: "ERC20"},
		Method: &abi.Method{Name: "constructor"},
		Args:   []string{},
	}

	chainIDs := rpc.GlobalClient().Networks(context.Background())

	e := &envelope.Envelope{
		Chain:  &common.Chain{Id: chainIDs[0].String()},
		Sender: &common.Account{Addr: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"},
		Call:   call,
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Gas: uint64(2000000),
			},
		},
		Metadata: &envelope.Metadata{
			Id: uuid.NewV4().String(),
		},
	}

	msg := &sarama.ProducerMessage{
		Topic: viper.GetString("kafka.topic.crafter"),
	}

	_ = encoding.Marshal(e, msg)

	log.WithFields(log.Fields{
		"msg": msg,
	}).Info("e2e")

	return msg
}

func SendTx() {

	p := broker.GlobalSyncProducer()

	rounds := 50
	for i := 0; i < rounds; i++ {
		partition, offset, err := p.SendMessage(NewMessage(i))
		if err != nil {
			log.WithError(err).Errorf("producer: could not produce message")
		}
		log.WithFields(log.Fields{
			"kafka.out.partition": partition,
			"kafka.out.offset":    offset,
		}).Info("e2e: message sent")
	}

}
