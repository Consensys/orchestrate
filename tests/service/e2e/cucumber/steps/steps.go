package steps

import (
	gohttp "net/http"

	"github.com/Shopify/sarama"
	"github.com/cucumber/godog"
	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/gofrs/uuid"
	"github.com/mitchellh/copystructure"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt/generator"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/redis"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	rpcClient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/client"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber/alias"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/utils"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/tracker"
)

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

	// API
	client orchestrateclient.OrchestrateClient

	// Producer to producer envelopes in topics
	producer sarama.SyncProducer

	logger *log.Entry

	jwtGenerator *generator.JWTGenerator

	ec ethclient.Client

	nonceManager nonce.Manager

	TearDownFunc []func()
}

func NewScenarioContext(
	chanReg *chanregistry.ChanRegistry,
	httpClient *gohttp.Client,
	chainReg chainregistry.ChainRegistryClient,
	contractRegistry registry.ContractRegistryClient,
	client orchestrateclient.OrchestrateClient,
	producer sarama.SyncProducer,
	aliasesReg *alias.Registry,
	jwtGenerator *generator.JWTGenerator,
	ec ethclient.Client,
	nonceManager nonce.Manager,
) *ScenarioContext {
	sc := &ScenarioContext{
		chanReg:          chanReg,
		httpClient:       httpClient,
		aliases:          aliasesReg,
		ChainRegistry:    chainReg,
		ContractRegistry: contractRegistry,
		client:           client,
		producer:         producer,
		logger:           log.NewEntry(log.StandardLogger()),
		jwtGenerator:     jwtGenerator,
		ec:               ec,
		nonceManager:     nonceManager,
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
	for _, topic := range utils.TOPICS {
		var ckey string
		if e != nil {
			ckey = utils2.LongKeyOf(topic, e.GetID())
		} else {
			ckey = utils2.ShortKeyOf(topic, sc.Pickle.Id)
		}

		// Create channel
		// TODO: make chan size configurable
		var ch = make(chan *tx.Envelope, 30)
		// Register channel on channel registry
		log.WithFields(log.Fields{
			"id":          ckey,
			"scenario.id": sc.Pickle.Id,
			"topic":       topic,
		}).Debugf("tracker: registered new envelope channel")
		sc.chanReg.Register(ckey, ch)

		// Add channel as a tracker output
		t.AddOutput(topic, ch)
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
		http.NewClient(http.NewDefaultConfig()),
		chainregistry.GlobalClient(),
		contractregistry.GlobalClient(),
		orchestrateclient.GlobalClient(),
		broker.GlobalSyncProducer(),
		alias.GlobalAliasRegistry(),
		generator.GlobalJWTGenerator(),
		rpcClient.GlobalClient(),
		redis.GlobalNonceManager(),
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
