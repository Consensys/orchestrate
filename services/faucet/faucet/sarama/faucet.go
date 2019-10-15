package sarama

import (
	"context"
	"math/big"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/ethereum"
)

// Faucet allows to credit by sending messages to a Kafka topic
type Faucet struct {
	p sarama.SyncProducer
}

// NewFaucet creates a New Faucet that can send message to a Kafka Topic
func NewFaucet(p sarama.SyncProducer) *Faucet {
	return &Faucet{
		p: p,
	}
}

func (f *Faucet) prepareMsg(r *types.Request, msg *sarama.ProducerMessage) error {
	// Create Trace for Crediting message
	e := &envelope.Envelope{
		Chain: (&chain.Chain{}).SetID(r.ChainID),
		From:  ethereum.HexToAccount(r.Creditor.Hex()),
		Tx: &ethereum.Transaction{
			TxData: (&ethereum.TxData{}).SetValue(r.Amount).SetTo(r.Beneficiary),
		},
	}

	// Unmarshal envelope
	err := encoding.Marshal(e, msg)
	if err != nil {
		return err
	}

	// Message should be sent to crafter topic
	msg.Topic = viper.GetString("kafka.topic.crafter")
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(r.ChainID, r.Creditor))

	return nil
}

// Credit process a Faucet credit request
func (f *Faucet) Credit(ctx context.Context, r *types.Request) (*big.Int, bool, error) {
	// Prepare Message
	msg := &sarama.ProducerMessage{}
	err := f.prepareMsg(r, msg)
	if err != nil {
		return big.NewInt(0), false, errors.FromError(err).ExtendComponent(component)
	}

	// Send message
	partition, offset, err := f.p.SendMessage(msg)
	if err != nil {
		return big.NewInt(0), false, errors.FromError(err).ExtendComponent(component)
	}

	log.WithFields(log.Fields{
		"kafka.out.partition": partition,
		"kafka.out.offset":    offset,
		"kafka.out.topic":     msg.Topic,
	}).Tracef("faucet: message produced")

	return r.Amount, true, nil
}