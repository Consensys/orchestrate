package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

var (
	senders = []string{
		"0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
		// "0xf5956Eb46b377Ae41b41BDa94e6270208d8202bb",
		// "0x93f7274c9059e601be4512F656B57b830e019E41",
		// "0xbfc7137876d7Ac275019d70434B0f0779824a969",
		// "0xA8d8DB1d8919665a18212374d623fc7C0dFDa410",
	}
	// ERC1400Address of token contract to target
	ERC1400Address = "0x8f371DAA8A5325f53b754A7017Ac3803382bc847"
)

// NewMessage creates a New Sarama producer message
func NewMessage(i int) *sarama.ProducerMessage {
	var call *common.Call
	msgs := []string{
		// "constructor",
		"call",
	}

	switch msgs[i%len(msgs)] {
	case "constructor":
		call = &common.Call{
			Contract: &abi.Contract{Name: "ERC1400"},
			Method:   &abi.Method{Name: "constructor"},
			Args:     []string{"0xabcd", "0xabcd", "0x10", "[0xcd626bc764e1d553e0d75a42f5c4156b91a63f23,0xcd626bc764e1d553e0d75a42f5c4156b91a63f23]", "0xcd626bc764e1d553e0d75a42f5c4156b91a63f23", "0xabcd"},
		}

	case "call":
		call = &common.Call{
			Contract: &abi.Contract{Name: "ERC1400"},
			Method:   &abi.Method{Name: "setDocument"},
			Args:     []string{"0xabcd", "0xabcd", "0xabcd"},
		}
	}

	e := &envelope.Envelope{
		Chain:  &common.Chain{Id: "42"},
		Sender: &common.Account{Addr: senders[i%len(senders)]},
		Call:   call,
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				To: ERC1400Address,
			},
		},
	}

	msg := &sarama.ProducerMessage{
		Topic:     viper.GetString("kafka.topic.crafter"),
		Partition: -1,
	}

	_ = encoding.Marshal(e, msg)

	return msg
}

func main() {
	broker.InitSyncProducer(context.Background())

	rounds := 1
	for i := 0; i < rounds; i++ {
		partition, offset, err := broker.GlobalSyncProducer().SendMessage(NewMessage(i))
		if err != nil {
			log.WithError(err).Error("e2e: could not send message")
		}
		log.WithFields(log.Fields{
			"partition": partition,
			"offset":    offset,
		}).Info("e2e: message sent")
	}
}
