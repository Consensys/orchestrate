package steps

import "C"
import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/service/chanregistry"
)

type ScenarioContext struct {
	ScenarioID string

	// Topics -> chan *envelope.Envelope
	EnvelopesChan map[string]chan *envelope.Envelope

	// MetadataId -> *envelope.Envelope (keep envelopes to be processed and checked)
	Envelopes map[string]*envelope.Envelope

	// Used to keep alias of contract addresses for example
	Value map[string]interface{}

	Logger *log.Entry
}

// initScenarioContext initilize a scenario context - create a random scenario id - initilize a logger enrich with the scenario name - initilize envelope chan
func (sc *ScenarioContext) initScenarioContext(s interface{}) {

	var scenarioName string
	switch t := s.(type) {
	case *gherkin.Scenario:
		scenarioName = t.Name
	case *gherkin.ScenarioOutline:
		scenarioName = t.Name
	}

	sc.ScenarioID = uuid.NewV4().String()

	sc.Logger = log.StandardLogger().WithFields(log.Fields{
		"Sceneario": scenarioName,
	})

	topics := []string{
		viper.GetString("kafka.topic.crafter"),
		viper.GetString("kafka.topic.nonce"),
		viper.GetString("kafka.topic.signer"),
		viper.GetString("kafka.topic.sender"),
		viper.GetString("kafka.topic.decoded"),
	}
	if primary := viper.GetString("cucumber.chainid.primary"); primary != "" {
		topics = append(topics, fmt.Sprintf("%s-%s", viper.GetString("kafka.topic.decoder"), primary))
	}
	if secondary := viper.GetString("cucumber.chainid.secondary"); secondary != "" {
		topics = append(topics, fmt.Sprintf("%s-%s", viper.GetString("kafka.topic.decoder"), secondary))
	}

	sc.EnvelopesChan = make(map[string]chan *envelope.Envelope)
	sc.Envelopes = make(map[string]*envelope.Envelope)
	sc.Value = make(map[string]interface{})
	r := chanregistry.GlobalChanRegistry()
	for _, topic := range topics {
		sc.EnvelopesChan[topic] = r.NewEnvelopeChan(sc.ScenarioID, topic)
	}

}

func (sc *ScenarioContext) iHaveTheFollowingEnvelope(rawEnvelopes *gherkin.DataTable) error {

	head := rawEnvelopes.Rows[0].Cells

	for i := 1; i < len(rawEnvelopes.Rows); i++ {
		mapEnvelope := make(map[string]string)
		for j, cell := range head {
			// Replace "Aliases"
			switch {
			case cell.Value == "AliasTo":
				val := sc.Value[rawEnvelopes.Rows[i].Cells[j].Value]
				switch v := val.(type) {
				default:
					panic(fmt.Errorf("unexpected type %T", v))
				case string:
					mapEnvelope["to"] = val.(string)
				case *ethereum.Account:
					mapEnvelope["to"] = val.(*ethereum.Account).Hex()
				}
			case cell.Value == "AliasChainId":
				mapEnvelope["chainId"] = viper.GetString(fmt.Sprintf("cucumber.chainid.%s", rawEnvelopes.Rows[i].Cells[j].Value))
			default:
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
		sc.Logger.Infof("cucumber: sending envelope %+v", e)
		err := SendEnvelope(e)
		if err != nil {
			sc.Logger.Errorf("cucumber: sending envelopes to CoreStack failed with error %q", err)
			return err
		}
	}

	sc.Logger.Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) coreStackShouldReceiveEnvelopes() error {

	topic := viper.GetString("kafka.topic.crafter")
	e, err := ReadChanWithTimeout(sc.EnvelopesChan[topic], viper.GetInt64("cucumber.steps.timeout"), len(sc.Envelopes))
	if err != nil {
		sc.Logger.Errorf("cucumber: receiving envelopes from CoreStack failed with error %q", err)
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
	e, err := ReadChanWithTimeout(sc.EnvelopesChan[topic], viper.GetInt64("cucumber.steps.timeout"), len(sc.Envelopes))
	if err != nil {
		sc.Logger.Errorf("cucumber: reading from topic %s failed with error %q", topic, err)
		return err
	}

	for _, v := range e {
		if v.GetTx().GetTxData().GetData() == nil {
			err := fmt.Errorf("tx-crafter could not craft transaction")
			sc.Logger.Errorf("cucumber: step failed with error %q", err)
			return err
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
	e, err := ReadChanWithTimeout(sc.EnvelopesChan[topic], viper.GetInt64("cucumber.steps.timeout"), len(sc.Envelopes))
	if err != nil {
		sc.Logger.Errorf("cucumber: failed to read from topic %s with error %q", topic, err)
		return err
	}

	nonces := make(map[string]map[string]map[uint64]bool)
	for _, v := range e {
		chain := v.GetChain().ID().String()
		addr := v.GetFrom().Address().Hex()
		nonce := v.GetTx().GetTxData().GetNonce()

		if nonces[chain] == nil {
			nonces[chain] = make(map[string]map[uint64]bool)
		}
		if nonces[chain][addr] == nil {
			nonces[chain][addr] = make(map[uint64]bool)
		}
		if nonces[chain][addr][nonce] {
			err := fmt.Errorf("tx-nonce set 2 times the same nonce: %d", v.GetTx().GetTxData().GetNonce())
			sc.Logger.Errorf("cucumber: step failed with error %q", err)
			return err
		}
		nonces[chain][addr][nonce] = true
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) theTxsignerShouldSign() error {

	topic := viper.GetString("kafka.topic.sender")
	e, err := ReadChanWithTimeout(sc.EnvelopesChan[topic], viper.GetInt64("cucumber.steps.timeout"), len(sc.Envelopes))
	if err != nil {
		sc.Logger.Errorf("cucumber: failed to read from topic %s with error %q", topic, err)
		return err
	}

	for _, v := range e {
		if v.GetTx().GetRaw() == nil {
			err := fmt.Errorf("tx-signer could not sign")
			sc.Logger.Errorf("cucumber: step failed with error %q", err)
			return err
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

	return nil
}

func (sc *ScenarioContext) theTxlistenerShouldCatchTheTx() error {

	for chain, count := range GetChainCounts(sc.Envelopes) {

		topic := fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chain)
		e, err := ReadChanWithTimeout(sc.EnvelopesChan[topic], viper.GetInt64("cucumber.steps.miningtimeout"), int(count))
		if err != nil {
			sc.Logger.Errorf("cucumber: failed to read from topic %s failed with error %q", topic, err)
			return err
		}

		for _, v := range e {
			if v.GetReceipt().GetTxHash() == nil {
				err := fmt.Errorf("tx-listener could not catch the tx")
				sc.Logger.Errorf("cucumber: step failed with error %q", err)
				return err
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
	topic := viper.GetString("kafka.topic.decoded")
	channel := sc.EnvelopesChan[topic]
	e, err := ReadChanWithTimeout(channel, viper.GetInt64("cucumber.steps.timeout"), len(sc.Envelopes))
	if err != nil {
		sc.Logger.Errorf("cucumber: failed to read from topic %s with error %q", topic, err)
		return err
	}

	for _, v := range e {
		for _, log := range v.GetReceipt().GetLogs() {
			if len(log.GetTopics()) > 0 {
				if len(log.GetDecodedData()) == 0 {
					err := fmt.Errorf("tx-decoder could not decode the transaction")
					sc.Logger.Errorf("cucumber: step failed with error %q", err)
					return err
				}
			}
		}
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        topic,
	}).Info("cucumber: step check")

	return nil
}

func (sc *ScenarioContext) beforeStep(s *gherkin.Step) {

	sc.Logger = sc.Logger.WithFields(log.Fields{
		"Step": s.Text,
	})
}

func (sc *ScenarioContext) iShouldCatchTheirContractAddresses() error {

	topic := viper.GetString("kafka.topic.decoded")
	channel := sc.EnvelopesChan[topic]
	e, err := ReadChanWithTimeout(channel, viper.GetInt64("cucumber.steps.timeout"), len(sc.Envelopes))
	if err != nil {
		sc.Logger.Errorf("cucumber: failed to read from a topic %s with error %q", topic, err)
		return err
	}

	// Consume unused envelopes stuck in chan
	topics := []string{
		viper.GetString("kafka.topic.crafter"),
		viper.GetString("kafka.topic.nonce"),
		viper.GetString("kafka.topic.signer"),
		viper.GetString("kafka.topic.sender"),
	}
	for _, topic := range topics {
		_, _ = ReadChanWithTimeout(sc.EnvelopesChan[topic], viper.GetInt64("cucumber.steps.timeout"), len(e))
	}
	for chain, count := range GetChainCounts(sc.Envelopes) {
		_, _ = ReadChanWithTimeout(sc.EnvelopesChan[fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chain)], viper.GetInt64("cucumber.steps.timeout"), int(count))
	}

	for _, v := range e {
		if v.GetReceipt().GetContractAddress() == nil {
			return fmt.Errorf("could not deploy contract")
		}
		sc.Value[v.GetMetadata().GetExtra()["AliasContractInstance"]] = v.GetReceipt().GetContractAddress()

		// Envelope processed - remove from scenario context
		delete(sc.Envelopes, v.GetMetadata().GetId())
	}

	sc.Logger.WithFields(log.Fields{
		"EnvelopeReceived": len(e),
		"msg.Topic":        channel,
	}).Info("cucumber: step check")

	return nil
}

func FeatureContext(s *godog.Suite) {

	sc := &ScenarioContext{}

	s.BeforeScenario(sc.initScenarioContext)
	s.BeforeStep(sc.beforeStep)

	s.Step(`^I have the following envelope:$`, sc.iHaveTheFollowingEnvelope)
	s.Step(`^I send these envelopes to CoreStack$`, sc.iSendTheseEnvelopeToCoreStack)
	s.Step(`^CoreStack should receive envelopes$`, sc.coreStackShouldReceiveEnvelopes)
	s.Step(`^the tx-crafter should set the data$`, sc.theTxcrafterShouldSetTheData)
	s.Step(`^the tx-nonce should set the nonce$`, sc.theTxnonceShouldSetTheNonce)
	s.Step(`^the tx-signer should sign$`, sc.theTxsignerShouldSign)
	s.Step(`^the tx-sender should send the tx$`, sc.theTxsenderShouldSendTheTx)
	s.Step(`^the tx-listener should catch the tx$`, sc.theTxlistenerShouldCatchTheTx)
	s.Step(`^the tx-decoder should decode$`, sc.theTxdecoderShouldDecode)
	s.Step(`^I should catch their contract addresses$`, sc.iShouldCatchTheirContractAddresses)
}
