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
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/services/chanregistry"
)

type ScenarioContext struct {
	Name        string
	ScenarioID  string

	// Topics -> chan *envelope.Envelope
	EnvelopesChan map[string]chan *envelope.Envelope

	// MetadataId -> *envelope.Envelope
	Envelopes map[string]*envelope.Envelope

	ChainIDs []*big.Int

	Value map[string]interface{}

	Logger *log.Entry
}

func (sc *ScenarioContext) initScenarioContext(s interface{}) {

	sc.Logger = log.NewEntry(log.StandardLogger())

	var scenarioName string
	switch t := s.(type) {
	case *gherkin.Scenario:
		scenarioName = t.Name
	case *gherkin.ScenarioOutline:
		scenarioName = t.Name
	}

	sc.ScenarioID = uuid.NewV4().String()

	sc.Logger = sc.Logger.WithFields(log.Fields{
		"Sceneario": scenarioName,
	})

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
	sc.Value = make(map[string]interface{})

	r := chanregistry.GlobalChanRegistry()
	for _, topic := range topics {
		sc.EnvelopesChan[topic] = r.NewEnvelopeChan(sc.ScenarioID, topic)
	}

	sc.ChainIDs = chainIDs
}

func (sc *ScenarioContext) iHaveTheFollowingEnvelope(rawEnvelopes *gherkin.DataTable) error {

	head := rawEnvelopes.Rows[0].Cells

	for i := 1; i < len(rawEnvelopes.Rows); i++ {
		mapEnvelope := make(map[string]string)
		for j, cell := range head {
			// Replace "to" contract alias to address
			if cell.Value == "to" {
				mapEnvelope[cell.Value] = sc.Value[rawEnvelopes.Rows[i].Cells[j].Value].(string)
			} else {
				mapEnvelope[cell.Value] = rawEnvelopes.Rows[i].Cells[j].Value
			}
		}

		if mapEnvelope["metadataID"] == "" {
			mapEnvelope["metadataID"] = uuid.NewV4().String()
		}

		mapEnvelope["ScenarioID"] = sc.ScenarioID
		e := EnvelopeCrafter(mapEnvelope)
		sc.Envelopes[mapEnvelope["metadataID"]] = e
	}

	sc.Logger.Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) iSendTheseEnvelopeToCoreStack() error {

	for _, e := range sc.Envelopes {
		err := SendEnvelope(e)
		if err != nil {
			return err
		}
	}

	sc.Logger.Info("cucumber: step check")

	return nil
}


func (sc *ScenarioContext) coreStackShouldReceiveThem() error {

	topic := viper.GetString("kafka.topic.crafter")
	e, err := ChanTimeout(sc.EnvelopesChan[topic], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxcrafterShouldSetTheData() error {

	topic := viper.GetString("kafka.topic.nonce")
	e, err := ChanTimeout(sc.EnvelopesChan[topic], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		if v.GetTx().GetTxData().GetData() == "" {
			return fmt.Errorf("tx-crafter could not craft transaction")
		}
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxnonceShouldSetTheNonce() error {

	topic := viper.GetString("kafka.topic.signer")
	e, err := ChanTimeout(sc.EnvelopesChan[topic], 10, len(sc.Envelopes))
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

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxsignerShouldSign() error {

	topic := viper.GetString("kafka.topic.sender")
	e, err := ChanTimeout(sc.EnvelopesChan[topic], 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		if v.GetTx().GetRaw() == "" {
			return fmt.Errorf("tx-signer could not sign")
		}
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxsenderShouldSendTheTx() error {
	// TODO call API envelope store

	// for _, v := range sc.Envelopes {
	// 	status, _, err := grpcStore.GlobalEnvelopeStore().GetStatus(context.Background(), v.GetMetadata().GetId())
	// 	log.Infof("Status: %s", status)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (sc *ScenarioContext) theTxlistenerShouldCatchTheTx() error {

	for _, chain := range sc.ChainIDs {
		topic := fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chain.String())
		e, err := ChanTimeout(sc.EnvelopesChan[topic], 10, len(sc.Envelopes))
		if err != nil {
			return err
		}
	
		for _, v := range e {
			if v.GetReceipt().GetContractAddress() == "" {
				return fmt.Errorf("tx-listener could not catch the tx")
			}
		}

		sc.Logger.WithFields(log.Fields{
			"EnvelopeReceived": len(e),
			"msg.Topic":        topic,
		}).Info("cucumber: step check")
	}




	return nil
}

func (sc *ScenarioContext) theTxdecoderShouldDecode() error {

	topic := sc.EnvelopesChan[viper.GetString("kafka.topic.decoded")]
	e, err := ChanTimeout(topic, 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	for _, v := range e {
		for _, log := range v.GetReceipt().GetLogs() {
			if len(log.GetTopics()) > 0 {
				if len(log.GetDecodedData()) == 0 {
					return fmt.Errorf("tx-decoder could not decode the transaction")
				}
			}
		}
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic": topic,
	}).Info("cucumber: step check")

	return nil
}


func (sc *ScenarioContext) beforeStep(s *gherkin.Step) {

	sc.Logger = sc.Logger.WithFields(log.Fields{
		"Step": s.Text,
	})
}

func (sc *ScenarioContext) iShouldCatchTheirContractAddresses() error {

	topic := sc.EnvelopesChan[viper.GetString("kafka.topic.decoded")]
	e, err := ChanTimeout(topic, 10, len(sc.Envelopes))
	if err != nil {
		return err
	}

	if sc.Value == nil {
		sc.Value = make(map[string]interface{})
	}
	
	for _, v := range e {
		if v.GetReceipt().GetContractAddress() == "" {
			return fmt.Errorf("Could not deploy contract")
		}
		sc.Value[v.GetMetadata().GetExtra()["Alias"]] = v.GetReceipt().GetContractAddress()
	}
	
	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil

}

func FeatureContext(s *godog.Suite) {

	sc := &ScenarioContext{}

	s.BeforeScenario(sc.initScenarioContext)
	s.BeforeStep(sc.beforeStep)

	s.Step(`^I have the following envelope:$`, sc.iHaveTheFollowingEnvelope)
	s.Step(`^I send these envelopes to CoreStack$`, sc.iSendTheseEnvelopeToCoreStack)
	s.Step(`^CoreStack should receive them$`, sc.coreStackShouldReceiveThem)
	s.Step(`^the tx-crafter should set the data$`, sc.theTxcrafterShouldSetTheData)
	s.Step(`^the tx-nonce should set the nonce$`, sc.theTxnonceShouldSetTheNonce)
	s.Step(`^the tx-signer should sign$`, sc.theTxsignerShouldSign)
	s.Step(`^the tx-sender should send the tx$`, sc.theTxsenderShouldSendTheTx)
	s.Step(`^the tx-listener should catch the tx$`, sc.theTxlistenerShouldCatchTheTx)
	s.Step(`^the tx-decoder should decode$`, sc.theTxdecoderShouldDecode)
	s.Step(`^I should catch their contract addresses$`, sc.iShouldCatchTheirContractAddresses)
}