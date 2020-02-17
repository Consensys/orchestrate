package sarama

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

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
	b := tx.NewEnvelope().
		SetID(uuid.NewV4().String()).
		SetChainName(r.ChainName).
		SetFrom(r.Creditor).
		SetValue(r.Amount).
		SetTo(r.Beneficiary).
		SetChainID(r.ChainID).
		SetChainUUID(r.ChainUUID)

	if authToken := authutils.AuthorizationFromContext(ctx); authToken != "" {
		_ = b.SetHeadersValue(multitenancy.AuthorizationMetadata, authToken)
	}
	if apiKey := authutils.APIKeyFromContext(ctx); apiKey != "" {
		_ = b.SetHeadersValue(authentication.APIKeyHeader, apiKey)
	}

	// Unmarshal envelope
	err := encoding.Marshal(b.TxEnvelopeAsRequest(), msg)
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
