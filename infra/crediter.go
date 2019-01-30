package infra

import (
	"context"
	"math/big"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
)

// SaramaCrediter allows to credit by sending messages to a Kafka topic
type SaramaCrediter struct {
	conf      FaucetConfig
	addresses map[string]common.Address

	p sarama.SyncProducer
	m *infSarama.Marshaller
}

// NewSaramaCrediter creates a new SaramaCrediter
func NewSaramaCrediter(conf FaucetConfig, p sarama.SyncProducer) (*SaramaCrediter, error) {
	addresses := map[string]common.Address{}
	for k, v := range conf.Addresses {
		addresses[k] = common.HexToAddress(v)
	}

	return &SaramaCrediter{
		conf:      conf,
		addresses: addresses,
		p:         p,
		m: infSarama.NewMarshaller(),
	}, nil
}

// Credit credit a given request by sending a message to a Kafka topic
func (c *SaramaCrediter) Credit(ctx context.Context, r *services.FaucetRequest) (*big.Int, bool, error) {
	// Prepare Faucet message
	msg, err := c.PrepareFaucetMsg(r)
	if err != nil {
		return big.NewInt(0), false, err
	}

	// Send message
	_, _, err = c.p.SendMessage(&msg)
	if err != nil {
		return big.NewInt(0), false, err
	}

	return r.Value, true, nil
}

// PrepareFaucetMsg creates a credit message to send to a specific topic
func (c *SaramaCrediter) PrepareFaucetMsg(r *services.FaucetRequest) (sarama.ProducerMessage, error) {
	// Determine Address of the faucet for requested chain
	faucetAddress := c.addresses[r.ChainID.Text(10)]

	// Create Trace for Crediting message
	faucetTrace := types.NewTrace()
	faucetTrace.Chain().ID.Set(r.ChainID)
	faucetTrace.Sender().Address = &faucetAddress
	faucetTrace.Tx().SetValue(r.Value)
	faucetTrace.Tx().SetTo(&r.Address)

	// Create Producer message
	var msg sarama.ProducerMessage
	err := c.m.Marshal(faucetTrace, &msg)
	if err != nil {
		return sarama.ProducerMessage{}, err
	}
	msg.Topic = c.conf.Topic

	return msg, nil
}
