package steps

import (
	"encoding/hex"
	"fmt"
	gohttp "net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/parser"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/tracker"
)

const GenericNamespace = "_"

var TOPICS = [...]string{
	"tx.crafter",
	"tx.signer",
	"tx.sender",
	"tx.decoded",
	"tx.recover",
	"account.generator",
	"account.generated",
}

// NewID creates a 8 character long random id
func NewID() string {
	u := uuid.NewV4()
	buf := make([]byte, 8)
	hex.Encode(buf, u[0:4])
	return string(buf)
}

// ScenarioID generates a random scenario UUID
func ScenarioID(def *gherkin.ScenarioDefinition) string {
	return fmt.Sprintf("|%v|-%v", fmt.Sprintf("%-20v", def.Name)[:20], NewID())
}

// AuthSetup is container for authentication context data
type AuthSetup struct {
	authMethod string
	authData   string
}

// ScenarioContext is container for scenario context data
type ScenarioContext struct {
	ID         string
	Definition *gherkin.ScenarioDefinition

	// Parser to parse cucumber/gherkin entries
	parser *parser.Parser

	// trackers track envelopes that are generated within the test session
	// as they are processed in the system
	trackers []*tracker.Tracker

	// defaultTracker allows to capture envelopes that are generated
	// within the system (to be captured those envelopes should have scenario.id set)
	defaultTracker *tracker.Tracker

	// chanReg to register envelopes channels on trackers
	chanReg *chanregistry.ChanRegistry

	httpClient   *gohttp.Client
	httpResponse *gohttp.Response

	httpAliases *parser.AliasRegistry

	// RegistryClient
	ContractRegistry registry.ContractRegistryClient

	// Producer to producer envelopes in topics
	producer sarama.SyncProducer

	logger *log.Entry

	authSetup AuthSetup
}

func setServiceURL(sc *ScenarioContext) {

	sc.httpAliases.Set(GenericNamespace, "chain-registry", viper.GetString(chainregistry.ChainRegistryURLViperKey))
	sc.httpAliases.Set(GenericNamespace, "contract-registry", "http://contract-registry:8081")
	sc.httpAliases.Set(GenericNamespace, "envelope-store", "http://envelope-store:8081")
}

func NewScenarioContext(
	chanReg *chanregistry.ChanRegistry,
	httpClient *gohttp.Client,
	contractRegistry registry.ContractRegistryClient,
	producer sarama.SyncProducer,
	p *parser.Parser,
) *ScenarioContext {
	sc := &ScenarioContext{
		chanReg:          chanReg,
		httpClient:       httpClient,
		httpAliases:      parser.NewAliasRegistry(),
		ContractRegistry: contractRegistry,
		producer:         producer,
		parser:           p,
		logger:           log.NewEntry(log.StandardLogger()),
		authSetup:        AuthSetup{},
	}

	setServiceURL(sc)

	return sc
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

	// Compute scenario UUID
	sc.ID = ScenarioID(sc.Definition)

	// Prepare default tracker
	sc.defaultTracker = sc.newTracker(nil)

	// Enrich logger
	sc.logger = sc.logger.WithFields(log.Fields{
		"scenario.name": sc.Definition.Name,
		"scenario.id":   sc.ID,
	})
}

func (sc *ScenarioContext) newTracker(e *tx.Envelope) *tracker.Tracker {
	if e != nil {
		sc.setMetadata(e)
	}
	// Set envelope metadata so it can be tracked

	// Create tracker and attach envelope
	t := tracker.NewTracker()
	t.Current = e

	// Initialize output channels on tracker and register channels on channel registry
	for _, topic := range TOPICS {
		// Create channel
		// TODO: make chan size configurable
		ch := make(chan *tx.Envelope, 30)

		// Add channel as a tracker output
		t.AddOutput(topic, ch)

		// Register channel on channel registry
		if e != nil {
			sc.chanReg.Register(LongKeyOf(topic, sc.ID, e.GetID()), ch)
		} else {
			sc.chanReg.Register(ShortKeyOf(topic, sc.ID), ch)
		}
	}

	return t
}

func (sc *ScenarioContext) setMetadata(e *tx.Envelope) {
	// Prepare envelope metadata
	_ = e.SetID(uuid.NewV4().String()).
		SetContextLabelsValue("debug", "true").
		SetContextLabelsValue("scenario.id", sc.ID).
		SetContextLabelsValue("scenario.name", sc.Definition.Name)
}

func (sc *ScenarioContext) newTrackers(envelopes []*tx.Envelope) []*tracker.Tracker {
	// Create a tracker for every envelope
	var trackers []*tracker.Tracker
	for _, e := range envelopes {
		// Create a tracker
		trackers = append(trackers, sc.newTracker(e))
	}

	return trackers
}

func (sc *ScenarioContext) setTrackers(trackers []*tracker.Tracker) {
	sc.trackers = trackers
}

func (sc *ScenarioContext) sendEnvelope(topic string, e *tx.Envelope) error {
	// Prepare message to be sent
	msg := &sarama.ProducerMessage{Topic: viper.GetString(fmt.Sprintf("topic.%v", topic))}

	err := encoding.Marshal(e.TxEnvelopeAsRequest(), msg)
	if err != nil {
		return err
	}

	// Send message
	_, _, err = sc.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id":            e.GetID(),
		"scenario.id":   sc.ID,
		"scenario.name": sc.Definition.Name,
	}).Debugf("scenario: envelope sent")

	return nil
}

func (sc *ScenarioContext) iSendEnvelopesToTopic(topic string, table *gherkin.DataTable) error {
	// Parse table
	envelopes, err := sc.parser.ParseEnvelopes(sc.ID, table)
	if err != nil {
		return errors.DataError("could not parse tx request - got %v", err)
	}

	// Set trackers for each envelope
	sc.setTrackers(sc.newTrackers(envelopes))

	// Send envelopes
	for _, t := range sc.trackers {
		err := sc.sendEnvelope(topic, t.Current)
		if err != nil {
			return errors.InternalError("could not send tx request - got %v", err)
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
	err = sc.sendEnvelope("tx.crafter", trackers[0].Current)
	if err != nil {
		return err
	}

	// Catch envelope after it has been decoded
	err = trackers[0].Load("tx.decoded", 30*time.Second)
	if err != nil {
		return fmt.Errorf("%v: no receipt for contract %q deployment", sc.ID, alias)
	}

	// Alias contract address
	if trackers[0].Current.GetReceipt().GetContractAddress() == "" {
		return fmt.Errorf("%v: contract %q could not be deployed", sc.ID, alias)
	}
	sc.parser.Aliases.Set(sc.ID, alias, trackers[0].Current.GetReceipt().GetContractAddress())

	return nil
}

func (sc *ScenarioContext) envelopeShouldBeInTopic(topic string) error {
	for i, t := range sc.trackers {
		err := t.Load(topic, viper.GetDuration(CucumberTimeoutViperKey))
		if err != nil {
			return fmt.Errorf("%v: envelope n°%v not in topic %q", sc.ID, i, topic)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHavePayloadSet() error {
	for i, t := range sc.trackers {
		if t.Current.GetData() == "" {
			return fmt.Errorf("%v: payload not set envelope n°%v", sc.ID, i)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveNonceSet() error {
	nonces := make(map[string]map[string]map[uint64]bool)
	for _, t := range sc.trackers {
		chain := t.Current.GetChainID().String()
		addr := t.Current.GetFromString()
		nonce, err := t.Current.GetNonceUint64()
		if err != nil {
			return err
		}

		if _, ok := nonces[chain]; !ok {
			nonces[chain] = make(map[string]map[uint64]bool)
		}
		if _, ok := nonces[chain][addr]; !ok {
			nonces[chain][addr] = make(map[uint64]bool)
		}
		if _, ok := nonces[chain][addr][nonce]; ok {
			return fmt.Errorf("%v: nonce %d attributed more than once", sc.ID, t.Current.Nonce)
		}
		nonces[chain][addr][nonce] = true
	}

	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveRawAndHashSet() error {
	for i, t := range sc.trackers {
		if t.Current.Raw == "" {
			return fmt.Errorf("%v: raw not set on envelope n°%v", sc.ID, i)
		}

		if t.Current.TxHash == nil {
			return fmt.Errorf("%v: hash not set on envelope n°%v", sc.ID, i)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveFromSet() error {
	for i, t := range sc.trackers {
		if t.Current.From == nil {
			return fmt.Errorf("%v: from not set on envelope n°%v", sc.ID, i)
		}
	}
	return nil
}

func (sc *ScenarioContext) envelopesShouldHaveLogDecoded() error {
	for i, t := range sc.trackers {
		for _, l := range t.Current.GetReceipt().GetLogs() {
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
		http.NewClient(),
		registryclient.GlobalClient(),
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

	initHTTP(s, sc)
	registerContractRegistrySteps(s, sc)

	s.Step(`^I have deployed contract "([^"]*)"$`, sc.iHaveDeployedContract)
	s.Step(`^I send envelopes to topic "([^"]*)"$`, sc.iSendEnvelopesToTopic)
	s.Step(`^Envelopes should be in topic "([^"]*)"$`, sc.envelopeShouldBeInTopic)
	s.Step(`^Envelopes should have payload set$`, sc.envelopesShouldHavePayloadSet)
	s.Step(`^Envelopes should have nonce set$`, sc.envelopesShouldHaveNonceSet)
	s.Step(`^Envelopes should have raw and hash set$`, sc.envelopesShouldHaveRawAndHashSet)
	s.Step(`^Envelopes should have from set$`, sc.envelopesShouldHaveFromSet)
	s.Step(`^Envelopes should have log decoded$`, sc.envelopesShouldHaveLogDecoded)
}
