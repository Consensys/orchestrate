package steps

import (
	"context"
	"fmt"
	"math/big"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"

	// grpcStore "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/cucumber/chanregistry"
)

type ScenarioContext struct {
	ScenarioID string

	// Topics -> chan *envelope.Envelope
	EnvelopesChan map[string]chan *envelope.Envelope

	// MetadataId -> *envelope.Envelope
	Envelopes map[string]*envelope.Envelope

	ChainIDs []*big.Int

	Value map[string]interface{}
}

func (sc *ScenarioContext) initScenarioContext(interface{}) {

	sc.ScenarioID = uuid.NewV4().String()

	chainIDs := rpc.GlobalClient().Networks(context.Background())
	topics := []string{
		viper.GetString("kafka.topic.crafter"),
		viper.GetString("kafka.topic.nonce"),
		viper.GetString("kafka.topic.signer"),
		viper.GetString("kafka.topic.sender"),
		viper.GetString("kafka.topic.decoded"),
	}
	for _, chainID := range chainIDs {
		topics = append(topics, fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chainID.String()))
	}

	sc.EnvelopesChan = make(map[string]chan *envelope.Envelope)
	sc.Envelopes = make(map[string]*envelope.Envelope)

	r := chanregistry.GlobalChanRegistry()
	for _, v := range topics {
		sc.EnvelopesChan[v] = r.NewEnvelopeChan(sc.ScenarioID, v)
	}

	sc.ChainIDs = chainIDs
}

func (sc *ScenarioContext) iHaveTheFollowingEnvelope(rawEnvelopes *gherkin.DataTable) error {

	head := rawEnvelopes.Rows[0].Cells

	for i := 1; i < len(rawEnvelopes.Rows); i++ {
		mapEnvelope := make(map[string]string)
		for j, cell := range head {
			mapEnvelope[cell.Value] = rawEnvelopes.Rows[i].Cells[j].Value
		}

		if mapEnvelope["metadataID"] == "" {
			mapEnvelope["metadataID"] = uuid.NewV4().String()
		}

		mapEnvelope["ScenarioID"] = sc.ScenarioID
		e := EnvelopeCrafter(mapEnvelope)
		sc.Envelopes[mapEnvelope["metadataID"]] = e
	}

	log.WithFields(log.Fields{
		"ScenarioID": sc.ScenarioID,
		"Step":       "iHaveTheFollowingEnvelope",
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) iSendTheseEnvelopeToCoreStack() error {

	for _, e := range sc.Envelopes {
		err := SendEnvelope(e)
		if err != nil {
			return err
		}
	}

	e, err := ChanTimeout(sc.EnvelopesChan[viper.GetString("kafka.topic.crafter")], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"ScenarioID":       sc.ScenarioID,
		"EnvelopeReceived": len(e),
		"msg.Topic":        viper.GetString("kafka.topic.crafter"),
		"Step":             "iSendTheseEnvelopeToCoreStack",
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxcrafterShouldSetTheData() error {
	e, err := ChanTimeout(sc.EnvelopesChan[viper.GetString("kafka.topic.nonce")], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		if v.GetTx().GetTxData().GetData() == "" {
			return fmt.Errorf("tx-crafter could not craft transaction")
		}
	}

	log.WithFields(log.Fields{
		"ScenarioID":       sc.ScenarioID,
		"EnvelopeReceived": len(e),
		"msg.Topic":        viper.GetString("kafka.topic.nonce"),
		"Step":             "theTxcrafterShouldSetTheData",
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxnonceShouldSetTheNonce() error {
	e, err := ChanTimeout(sc.EnvelopesChan[viper.GetString("kafka.topic.signer")], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}
	nonces := make(map[uint64]bool, len(sc.Envelopes))
	for _, v := range e {
		if nonces[v.GetTx().GetTxData().GetNonce()] {
			return fmt.Errorf("tx-nonce set 2 times the same nonce")
		}
		nonces[v.GetTx().GetTxData().GetNonce()] = true
	}

	log.WithFields(log.Fields{
		"ScenarioID":       sc.ScenarioID,
		"EnvelopeReceived": len(e),
		"msg.Topic":        viper.GetString("kafka.topic.signer"),
		"Step":             "theTxnonceShouldSetTheNonce",
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxsignerShouldSign() error {
	e, err := ChanTimeout(sc.EnvelopesChan[viper.GetString("kafka.topic.sender")], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		if v.GetTx().GetRaw() == "" {
			return fmt.Errorf("tx-signer could not sign")
		}
	}

	log.WithFields(log.Fields{
		"ScenarioID":       sc.ScenarioID,
		"EnvelopeReceived": len(e),
		"msg.Topic":        viper.GetString("kafka.topic.sender"),
		"Step":             "theTxsignerShouldSign",
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxsenderShouldSendTheTx() error {
	// TODO call API envelope store

	// for _, v := range sc.Envelopes {
	// 	status, _, err := grpcStore.GlobalEnvelopeStore().GetStatus(context.Background(), v.GetMetadata().GetId())
	// 	log.Infof("Status: %s", status)
	// 	if err != nil {
	// 		return fmt.Errorf("transaction not stored")
	// 	}
	// }

	return nil
}

func (sc *ScenarioContext) theTxlistenerShouldCatchTheTx() error {
	e, err := ChanTimeout(sc.EnvelopesChan[fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), sc.ChainIDs[0].String())], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		if v.GetReceipt().GetContractAddress() == "" {
			return fmt.Errorf("tx-listener could not catch the tx")
		}
	}

	log.WithFields(log.Fields{
		"ScenarioID":       sc.ScenarioID,
		"EnvelopeReceived": len(e),
		"msg.Topic":        fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), sc.ChainIDs[0].String()),
		"Step":             "theTxlistenerShouldCatchTheTx",
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxdecoderShouldDecode() error {
	e, err := ChanTimeout(sc.EnvelopesChan[viper.GetString("kafka.topic.decoded")], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		if v.GetReceipt().GetContractAddress() == "" {
			return fmt.Errorf("tx-listener could not catch the tx")
		}
	}

	log.WithFields(log.Fields{
		"ScenarioID":       sc.ScenarioID,
		"EnvelopeReceived": len(e),
		// "Metadata":   e.GetMetadata().GetId(),
		"msg.Topic": viper.GetString("kafka.topic.decoded"),
		"Step":      "theTxdecoderShouldDecode",
	}).Info("cucumber: step check")

	return nil
}

func FeatureContext(s *godog.Suite) {

	sc := &ScenarioContext{}

	s.BeforeScenario(sc.initScenarioContext)

	s.Step(`^I have the following envelope:$`, sc.iHaveTheFollowingEnvelope)
	s.Step(`^I send these envelope to CoreStack$`, sc.iSendTheseEnvelopeToCoreStack)
	s.Step(`^the tx-crafter should set the data$`, sc.theTxcrafterShouldSetTheData)
	s.Step(`^the tx-nonce should set the nonce$`, sc.theTxnonceShouldSetTheNonce)
	s.Step(`^the tx-signer should sign$`, sc.theTxsignerShouldSign)
	s.Step(`^the tx-sender should send the tx$`, sc.theTxsenderShouldSendTheTx)
	s.Step(`^the tx-listener should catch the tx$`, sc.theTxlistenerShouldCatchTheTx)
	s.Step(`^the tx-decoder should decode$`, sc.theTxdecoderShouldDecode)
}
