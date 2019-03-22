package infra

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
	commonpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
)

// SaramaCrediter allows to credit by sending messages to a Kafka topic
type SaramaCrediter struct {
	addresses map[string]common.Address

	p sarama.SyncProducer
	m *infSarama.Marshaller
}

func parseAddresses(addresses []string) (map[string]common.Address, error) {
	m := make(map[string]common.Address)
	for _, addr := range addresses {
		split := strings.Split(addr, ":")
		if len(split) != 2 {
			return nil, fmt.Errorf("Could not parse faucet address %q (expected format %q)", addr, "<chainID>:<address>")
		}

		if !common.IsHexAddress(split[1]) {
			return nil, fmt.Errorf("Invalid Ethereum address %q", split[1])
		}

		m[split[0]] = common.HexToAddress(split[1])
	}
	return m, nil
}

// NewSaramaCrediter creates a new SaramaCrediter
func NewSaramaCrediter(p sarama.SyncProducer) (*SaramaCrediter, error) {
	addresses, err := parseAddresses(viper.GetStringSlice("faucet.addresses"))
	if err != nil {
		return nil, err
	}

	return &SaramaCrediter{
		addresses: addresses,
		p:         p,
		m:         infSarama.NewMarshaller(),
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
	faucetTrace := &tracepb.Trace{}

	faucetTrace.Reset()
	fmt.Println(faucetTrace.Chain)
	faucetTrace.Chain = (&commonpb.Chain{}).SetID(r.ChainID)
	faucetTrace.Sender = &commonpb.Account{Addr: faucetAddress.Hex()}
	faucetTrace.Tx = &ethpb.Transaction{TxData: &ethpb.TxData{}}
	faucetTrace.Tx.TxData.SetValue(r.Value).SetTo(r.Address)

	// Create Producer message
	var msg sarama.ProducerMessage
	err := c.m.Marshal(faucetTrace, &msg)
	if err != nil {
		return sarama.ProducerMessage{}, err
	}
	msg.Topic = viper.GetString("faucet.topic")

	return msg, nil
}
