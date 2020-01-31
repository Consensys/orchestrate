package sarama

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"

	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"

	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
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

func (f *Faucet) prepareMsg(ctx context.Context, r *types.Request, msg *sarama.ProducerMessage) error {
	// Create Trace for Crediting message
	e := &envelope.Envelope{
		Chain: (&chain.Chain{}).SetChainID(r.ChainID).SetUUID(r.ChainUUID).SetName(r.ChainName),
		From:  ethereum.HexToAccount(r.Creditor.Hex()),
		Tx: &ethereum.Transaction{
			TxData: (&ethereum.TxData{}).SetValue(r.Amount).SetTo(r.Beneficiary),
		},
		Metadata: &envelope.Metadata{
			Id: uuid.NewV4().String(),
		},
	}

	if authToken := authutils.AuthorizationFromContext(ctx); authToken != "" {
		e.SetMetadataValue(multitenancy.AuthorizationMetadata, authToken)
	}
	if apiKey := authutils.APIKeyFromContext(ctx); apiKey != "" {
		e.SetMetadataValue(authentication.APIKeyHeader, apiKey)
	}

	// Unmarshal envelope
	err := encoding.Marshal(e, msg)
	if err != nil {
		return err
	}

	// Message should be sent to crafter topic
	msg.Topic = viper.GetString(broker.TxCrafterViperKey)
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(r.ChainID, r.Creditor))

	return nil
}

// Credit process a Faucet credit request
func (f *Faucet) Credit(ctx context.Context, r *types.Request) (*big.Int, error) {
	// Prepare Message
	msg := &sarama.ProducerMessage{}
	err := f.prepareMsg(ctx, r, msg)
	if err != nil {
		return big.NewInt(0), errors.FromError(err).ExtendComponent(component)
	}

	// Send message
	partition, offset, err := f.p.SendMessage(msg)
	if err != nil {
		return big.NewInt(0), errors.FromError(err).ExtendComponent(component)
	}

	log.WithFields(log.Fields{
		"kafka.out.partition": partition,
		"kafka.out.offset":    offset,
		"kafka.out.topic":     msg.Topic,
	}).Tracef("faucet: message produced")

	return r.Amount, nil
}
