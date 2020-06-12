package sarama

import (
	"context"
	"math/big"
	"reflect"

	"github.com/Shopify/sarama"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
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

func (f *Faucet) prepareMsg(ctx context.Context, r *types.Request, elected string, msg *sarama.ProducerMessage) error {
	// Create Trace for Crediting message
	b := tx.NewEnvelope().
		SetID(uuid.Must(uuid.NewV4()).String()).
		SetChainName(r.ChainName).
		SetFrom(r.FaucetsCandidates[elected].Creditor).
		SetValue(r.FaucetsCandidates[elected].Amount).
		SetTo(r.Beneficiary).
		SetChainID(r.ChainID).
		SetChainUUID(r.ChainUUID).
		SetContextLabelsValue("faucet.parentTxID", r.ParentTxID)

	if authToken := authutils.AuthorizationFromContext(ctx); authToken != "" {
		_ = b.SetHeadersValue(multitenancy.AuthorizationMetadata, authToken)
	}
	if apiKey := authutils.APIKeyFromContext(ctx); apiKey != "" {
		_ = b.SetHeadersValue(authutils.APIKeyHeader, apiKey)
	}

	// Unmarshal envelope
	err := encoding.Marshal(b.TxEnvelopeAsRequest(), msg)
	if err != nil {
		return err
	}

	// Message should be sent to crafter topic
	msg.Topic = viper.GetString(broker.TxCrafterViperKey)
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(r.ChainID, r.FaucetsCandidates[elected].Creditor))

	return nil
}

// Credit process a Faucet credit request
func (f *Faucet) Credit(ctx context.Context, r *types.Request) (*big.Int, error) {
	// Elect final faucet
	if len(r.FaucetsCandidates) == 0 {
		return nil, errors.FaucetWarning("no faucet request").ExtendComponent(component)
	}

	// Select a first faucet candidate for comparison
	r.ElectedFaucet = ElectFaucet(r.FaucetsCandidates)

	// Prepare Message
	msg := &sarama.ProducerMessage{}
	err := f.prepareMsg(ctx, r, r.ElectedFaucet, msg)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Send message
	partition, offset, err := f.p.SendMessage(msg)
	if err != nil {
		return nil, errors.KafkaConnectionError("could not send faucet transaction - got %v", err).ExtendComponent(component)
	}

	log.WithFields(log.Fields{
		"kafka.out.partition": partition,
		"kafka.out.offset":    offset,
		"kafka.out.topic":     msg.Topic,
	}).Tracef("faucet: message produced")

	return r.FaucetsCandidates[r.ElectedFaucet].Amount, nil
}

// ElectFaucet is currently selecting the remaining faucet candidates with the highest amount
func ElectFaucet(faucetsCandidates map[string]types.Faucet) string {
	// Select a first faucet candidate for comparison
	electedFaucet := reflect.ValueOf(faucetsCandidates).MapKeys()[0].String()
	for key, candidate := range faucetsCandidates {
		if candidate.Amount.Cmp(faucetsCandidates[electedFaucet].Amount) > 0 {
			electedFaucet = key
		}
	}
	return electedFaucet
}
