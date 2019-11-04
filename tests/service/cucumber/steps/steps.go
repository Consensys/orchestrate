package steps

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/services/contract-registry"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/services/contract-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/parser"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
)

var TOPICS = [...]string{
	"tx.crafter",
	"tx.nonce",
	"tx.signer",
	"tx.sender",
	"tx.decoded",
	"tx.recover",
	"wallet.generator",
	"wallet.generated",
}

// NewID creates a 8 character long random id
func NewID() string {
	u := uuid.NewV4()
	buf := make([]byte, 8)
	hex.Encode(buf, u[0:4])
	return string(buf)
}

// ScenarioID generates a random scenario ID
func ScenarioID(def *gherkin.ScenarioDefinition) string {
	return fmt.Sprintf("|%v|-%v", fmt.Sprintf("%-20v", def.Name)[:20], NewID())
}

// ScenarioContext is container for scenario context data
type ScenarioContext struct {
	ID         string
	Definition *gherkin.ScenarioDefinition

	// Parser to parse cucumber/gherkin entries
	parser *parser.Parser

	// trackers track envelopes that are generated within the test session
	// as they are processed in the system
	trackers []*tracker

	// defaultTracker allows to capture envelopes that are generated
	// within the system (to be captured those envelopes should have scenario.id set)
	defaultTracker *tracker

	// chanReg to register envelopes channels on trackers
	chanReg *chanregistry.ChanRegistry

	// RegistryClient
	Registry registry.RegistryClient

	// Producer to producer envelopes in topics
	producer sarama.SyncProducer

	logger *log.Entry
}

func NewScenarioContext(
	chanReg *chanregistry.ChanRegistry,
	reg registry.RegistryClient,
	producer sarama.SyncProducer,
	p *parser.Parser,
) *ScenarioContext {
	return &ScenarioContext{
		chanReg:  chanReg,
		Registry: reg,
		producer: producer,
		parser:   p,
		logger:   log.NewEntry(log.StandardLogger()),
	}
}

// initScenarioContext initialize a scenario context - create a random scenario id - initialize a logger enrich with the scenario name - initialize envelope chan
func (sc *ScenarioContext) init(s interface{}) {
	// Extract scenario definition
	switch t := s.(type) {
	case *gherkin.Scenario:
		sc.Definition = &t.ScenarioDefinition
	case *gherkin.ScenarioOutline:
		sc.Definition = &t.ScenarioDefinition
	}

	// Compute scenario ID
	sc.ID = ScenarioID(sc.Definition)

	// Prepare default tracker
	sc.defaultTracker = sc.newTracker(nil)

	// Enrich logger
	sc.logger = sc.logger.WithFields(log.Fields{
		"scenario.name": sc.Definition.Name,
		"scenario.id":   sc.ID,
	})
}

func (sc *ScenarioContext) newTracker(e *envelope.Envelope) *tracker {
	if e != nil {
		sc.setMetadata(e)
	}
	// Set envelope metadata so it can be tracked

	// Create tracker and attach envelope
	t := newTracker()
	t.current = e

	// Initialize output channels on tracker and register channels on channel registry
	for _, topic := range TOPICS {
		// Create channel
		// TODO: make chan size configurable
		ch := make(chan *envelope.Envelope, 30)

		// Add channel as a tracker output
		t.addOutput(topic, ch)

		// Register channel on channel registry
		if e != nil {
			sc.chanReg.Register(LongKeyOf(topic, sc.ID, e.GetMetadata().GetId()), ch)
		} else {
			sc.chanReg.Register(ShortKeyOf(topic, sc.ID), ch)
		}
	}

	return t
}

func (sc *ScenarioContext) setMetadata(e *envelope.Envelope) {
	// Prepare envelope metadata
	e.SetMetadataValue("debug", "true")
	e.SetMetadataValue("scenario.id", sc.ID)
	e.SetMetadataValue("scenario.name", sc.Definition.Name)
	e.GetMetadata().Id = uuid.NewV4().String()
}

func (sc *ScenarioContext) newTrackers(envelopes []*envelope.Envelope) []*tracker {
	// Create a tracker for every envelope
	var trackers []*tracker
	for _, e := range envelopes {
		// Create a tracker
		trackers = append(trackers, sc.newTracker(e))
	}

	return trackers
}

func (sc *ScenarioContext) setTrackers(trackers []*tracker) {
	sc.trackers = trackers
}

func (sc *ScenarioContext) sendEnvelope(topic string, e *envelope.Envelope) error {
	// Prepare message to be sent
	msg := &sarama.ProducerMessage{Topic: viper.GetString(fmt.Sprintf("topic.%v", topic))}
	err := encoding.Marshal(e, msg)
	if err != nil {
		return err
	}

	// Send message
	_, _, err = sc.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"metadata.id":   e.GetMetadata().GetId(),
		"scenario.id":   sc.ID,
		"scenario.name": sc.Definition.Name,
	}).Debugf("scenario: envelope sent")

	return nil
}

func (sc *ScenarioContext) iRegisterTheFollowingContract(table *gherkin.DataTable) error {
	// Parse table
	contracts, err := sc.parser.ParseContracts(sc.ID, table)
	if err != nil {
		return err
	}

	// Register contracts on the registry
	for _, contract := range contracts {
		_, err := sc.Registry.RegisterContract(
			context.Background(),
			&registry.RegisterContractRequest{
				Contract: contract,
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (sc *ScenarioContext) iSendEnvelopesToTopic(topic string, table *gherkin.DataTable) error {
	// Parse table
	envelopes, err := sc.parser.ParseEnvelopes(sc.ID, table)
	if err != nil {
		return err
	}

	// Set trackers for each envelope
	sc.setTrackers(sc.newTrackers(envelopes))

	// Send envelopes
	for _, t := range sc.trackers {
		err := sc.sendEnvelope(topic, t.current)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sc *ScenarioContext) iHaveDeployedContract(alias string, table *gherkin.DataTable) error {
	// Parse table
	envelopes, err := sc.parser.ParseEnvelopes(sc.ID, table)
	if err != nil {
		return err
	}

	// Set trackers
	trackers := sc.newTrackers(envelopes)

	if len(trackers) != 1 {
		return fmt.Errorf("%v: should deploy exactly 1 contract", sc.ID)
	}

	// Send envelope
	err = sc.sendEnvelope("tx.crafter", trackers[0].current)
	if err != nil {
		return err
	}

	// Catch envelope after it has been decoded
	err = trackers[0].load("tx.decoded", 30*time.Second)
	if err != nil {
		return fmt.Errorf("%v: no receipt for contract %q deployment", sc.ID, alias)
	}

	// Alias contract address
	if trackers[0].current.GetReceipt().GetContractAddress() == nil {
		return fmt.Errorf("%v: contract %q could not be deployed", sc.ID, alias)
	}
	sc.parser.Aliases.Set(sc.ID, alias, trackers[0].current.GetReceipt().GetContractAddress().Hex())

	return nil
}

func (sc *ScenarioContext) envelopeShouldBeInTopic(topic string) error {
	for i, t := range sc.trackers {
		err := t.load(topic, 5*time.Second)
		if err != nil {
			return fmt.Errorf("%v: envelope n°%v not in topic %q", sc.ID, i, topic)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHavePayloadSet() error {
	for i, t := range sc.trackers {
		if t.current.GetTx().GetTxData().GetData() == nil {
			return fmt.Errorf("%v: payload not set envelope n°%v", sc.ID, i)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveNonceSet() error {
	nonces := make(map[string]map[string]map[uint64]bool)
	for _, t := range sc.trackers {
		chain := t.current.GetChain().ID().String()
		addr := t.current.GetFrom().Address().Hex()
		nonce := t.current.GetTx().GetTxData().GetNonce()
		if _, ok := nonces[chain]; !ok {
			nonces[chain] = make(map[string]map[uint64]bool)
		}
		if _, ok := nonces[chain][addr]; !ok {
			nonces[chain][addr] = make(map[uint64]bool)
		}
		if _, ok := nonces[chain][addr][nonce]; ok {
			return fmt.Errorf("%v: nonce %d attributed more than once", sc.ID, t.current.GetTx().GetTxData().GetNonce())
		}
		nonces[chain][addr][nonce] = true
	}

	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveRawAndHashSet() error {
	for i, t := range sc.trackers {
		if t.current.GetTx().GetRaw() == nil {
			return fmt.Errorf("%v: raw not set on envelope n°%v", sc.ID, i)
		}

		if t.current.GetTx().GetHash() == nil {
			return fmt.Errorf("%v: hash not set on envelope n°%v", sc.ID, i)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveFromSet() error {
	for i, t := range sc.trackers {
		if t.current.GetFrom() == nil {
			return fmt.Errorf("%v: from not set on envelope n°%v", sc.ID, i)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveLogDecoded() error {
	for i, t := range sc.trackers {
		for _, l := range t.current.GetReceipt().GetLogs() {
			if len(l.GetTopics()) > 0 && len(l.GetDecodedData()) == 0 {
				return fmt.Errorf("%v: log have not been decoded on envelope n°%v", sc.ID, i)
			}
		}
	}

	return nil
}

// FeatureContext is a initializer for cucumber scenario methods
func FeatureContext(s *godog.Suite) {
	sc := NewScenarioContext(
		chanregistry.GlobalChanRegistry(),
		registryclient.GlobalContractRegistryClient(),
		broker.GlobalSyncProducer(),
		parser.GlobalParser(),
	)

	s.BeforeScenario(sc.init)

	s.BeforeStep(func(s *gherkin.Step) {
		log.WithFields(log.Fields{
			"step.text":     s.Text,
			"scenario.name": sc.Definition.Name,
			"scenario.id":   sc.ID,
		}).Debugf("scenario: step starts")
	})
	s.AfterStep(func(s *gherkin.Step, err error) {
		log.WithError(err).
			WithFields(log.Fields{
				"step.text":     s.Text,
				"scenario.name": sc.Definition.Name,
				"scenario.id":   sc.ID,
			}).Debugf("scenario: step completed")
	})

	s.Step(`^I register the following contract$`, sc.iRegisterTheFollowingContract)
	s.Step(`^I have deployed contract "([^"]*)"$`, sc.iHaveDeployedContract)
	s.Step(`^I send envelopes to topic "([^"]*)"$`, sc.iSendEnvelopesToTopic)
	s.Step(`^Envelopes should be in topic "([^"]*)"$`, sc.envelopeShouldBeInTopic)
	s.Step(`^Envelopes should have payload set$`, sc.envelopesShouldHavePayloadSet)
	s.Step(`^Envelopes should have nonce set$`, sc.envelopesShouldHaveNonceSet)
	s.Step(`^Envelopes should have raw and hash set$`, sc.envelopesShouldHaveRawAndHashSet)
	s.Step(`^Envelopes should have from set$`, sc.envelopesShouldHaveFromSet)
	s.Step(`^Envelopes should have log decoded$`, sc.envelopesShouldHaveLogDecoded)
}
