package steps

import (
	gohttp "net/http"

	"github.com/Shopify/sarama"
	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/gofrs/uuid"
	"github.com/mitchellh/copystructure"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/jwt/generator"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/alias"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/tracker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/utils"
)

var TOPICS = [...]string{
	"tx.crafter",
	"tx.signer",
	"tx.sender",
	"tx.decoded",
	"tx.recover",
	"account.generator",
	"account.generated",
}

// ScenarioContext is container for scenario context data
type ScenarioContext struct {
	Pickle *gherkin.Pickle

	// trackers track envelopes that are generated within the stepTable session
	// as they are processed in the system
	trackers []*tracker.Tracker

	// defaultTracker allows to capture envelopes that are generated
	// within the system (to be captured those envelopes should have scenario.id set)
	defaultTracker *tracker.Tracker

	// chanReg to register envelopes channels on trackers
	chanReg *chanregistry.ChanRegistry

	httpClient   *gohttp.Client
	httpResponse *gohttp.Response

	aliases *alias.Registry

	// Chain-Registry
	ChainRegistry chainregistry.ChainRegistryClient

	// RegistryClient
	ContractRegistry registry.ContractRegistryClient

	// Transaction Schedule
	TransactionScheduler txscheduler.TransactionSchedulerClient

	// Producer to producer envelopes in topics
	producer sarama.SyncProducer

	logger *log.Entry

	jwtGenerator *generator.JWTGenerator

	TearDownFunc []func()
}

func NewScenarioContext(
	chanReg *chanregistry.ChanRegistry,
	httpClient *gohttp.Client,
	chainReg chainregistry.ChainRegistryClient,
	contractRegistry registry.ContractRegistryClient,
	txScheduler txscheduler.TransactionSchedulerClient,
	producer sarama.SyncProducer,
	aliasesReg *alias.Registry,
	jwtGenerator *generator.JWTGenerator,
) *ScenarioContext {
	sc := &ScenarioContext{
		chanReg:              chanReg,
		httpClient:           httpClient,
		aliases:              aliasesReg,
		ChainRegistry:        chainReg,
		ContractRegistry:     contractRegistry,
		TransactionScheduler: txScheduler,
		producer:             producer,
		logger:               log.NewEntry(log.StandardLogger()),
		jwtGenerator:         jwtGenerator,
	}

	return sc
}

// initScenarioContext initialize a scenario context - create a random scenario id - initialize a logger enrich with the scenario name - initialize envelope chan
func (sc *ScenarioContext) init(s *gherkin.Pickle) {
	// Hook the Pickle to the scenario context
	sc.Pickle = s
	sc.aliases.Set(sc.Pickle.Id, sc.Pickle.Id, "scenarioID")

	// Prepare default tracker
	sc.defaultTracker = sc.newTracker(nil)

	// Enrich logger
	sc.logger = sc.logger.WithFields(log.Fields{
		"scenario.name": sc.Pickle.Name,
		"scenario.id":   sc.Pickle.Id,
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
			log.WithFields(log.Fields{
				"id":          e.GetID(),
				"scenario.id": sc.Pickle.Id,
				"topic":       topic,
			}).Debugf("registered new envelope")
			sc.chanReg.Register(utils.LongKeyOf(topic, e.GetID()), ch)
		} else {
			sc.chanReg.Register(utils.ShortKeyOf(topic, sc.Pickle.Id), ch)
		}
	}

	return t
}

func (sc *ScenarioContext) setMetadata(e *tx.Envelope) {
	if e.GetID() == "" {
		_ = e.SetID(uuid.Must(uuid.NewV4()).String())
	}
	// Prepare envelope metadata
	_ = e.SetContextLabelsValue("debug", "true").
		SetContextLabelsValue("scenario.id", sc.Pickle.Id).
		SetContextLabelsValue("scenario.name", sc.Pickle.Name)
}

func (sc *ScenarioContext) newTrackers(envelopes []*tx.Envelope) []*tracker.Tracker {
	// Create a tracker for every envelope
	var trackers []*tracker.Tracker
	for _, e := range envelopes {
		// Create a tracker
		sc.setMetadata(e)
		trackers = append(trackers, sc.newTracker(e))
	}

	return trackers
}

func (sc *ScenarioContext) setTrackers(trackers []*tracker.Tracker) {
	sc.trackers = trackers
}

type stepTable func(*gherkin.PickleStepArgument_PickleTable) error

func (sc *ScenarioContext) preProcessTableStep(tableFunc stepTable) stepTable {
	return func(table *gherkin.PickleStepArgument_PickleTable) error {
		err := sc.replaceAliases(table)
		if err != nil {
			return err
		}

		c, _ := copystructure.Copy(table)
		copyTable := c.(*gherkin.PickleStepArgument_PickleTable)

		return tableFunc(copyTable)
	}
}

func InitializeScenario(s *godog.ScenarioContext) {
	sc := NewScenarioContext(
		chanregistry.GlobalChanRegistry(),
		http.NewClient(),
		chainregistry.GlobalClient(),
		contractregistry.GlobalClient(),
		txscheduler.GlobalClient(),
		broker.GlobalSyncProducer(),
		alias.GlobalAliasRegistry(),
		generator.GlobalJWTGenerator(),
	)

	s.BeforeScenario(sc.init)
	s.AfterScenario(sc.tearDown)

	s.BeforeStep(func(s *gherkin.Pickle_PickleStep) {
		log.WithFields(log.Fields{
			"step.text":     s.Text,
			"scenario.name": sc.Pickle.Name,
			"scenario.id":   sc.Pickle.Id,
		}).Debugf("scenario: step starts")
	})
	s.AfterStep(func(s *gherkin.Pickle_PickleStep, err error) {
		log.WithError(err).
			WithFields(log.Fields{
				"step.text":     s.Text,
				"scenario.name": sc.Pickle.Name,
				"scenario.id":   sc.Pickle.Id,
			}).Debugf("scenario: step completed")
	})

	initEnvelopeSteps(s, sc)
	initHTTP(s, sc)
	registerContractRegistrySteps(s, sc)
}
