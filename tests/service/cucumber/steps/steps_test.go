package steps

import (
	"net/http"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	httpclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client"

	"github.com/Shopify/sarama/mocks"
	"github.com/cucumber/godog/gherkin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/chanregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/parser"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/tracker"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
	crc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client"
)

type ScenarioTestSuite struct {
	suite.Suite
	Context    *ScenarioContext
	chanReg    *chanregistry.ChanRegistry
	producer   *mocks.SyncProducer
	httpClient *http.Client
	crc        svc.ContractRegistryClient
}

func (s *ScenarioTestSuite) SetupSuite() {
	// Set viper configuration
	viper.Set(broker.TxCrafterViperKey, "tx-crafter")

	// Set channel registry
	s.chanReg = chanregistry.NewChanRegistry()
	s.producer = mocks.NewSyncProducer(s.T(), nil)
	s.httpClient = httpclient.GlobalClient()
	s.crc = crc.GlobalClient()
}

func (s *ScenarioTestSuite) SetupTest() {
	s.Context = NewScenarioContext(s.chanReg, s.httpClient, s.crc, s.producer, parser.New())
	sc := &gherkin.Scenario{}
	sc.Name = "test-scenario"
	s.Context.init(sc)
}

func (s *ScenarioTestSuite) TestInitScenarioContext() {
	scenario := &gherkin.Scenario{}
	scenario.Name = "test-1"
	s.Context.init(scenario)
	assert.Equal(s.T(), "test-1", s.Context.Definition.Name, "Context definition should have been set")
	assert.NotEqual(s.T(), "", s.Context.ID, "UUID should have been set")

	scenarioOutline := &gherkin.ScenarioOutline{}
	scenarioOutline.Name = "test-2"
	s.Context.init(scenarioOutline)
	assert.Equal(s.T(), "test-2", s.Context.Definition.Name, "Context definition should have been set")
}

func (s *ScenarioTestSuite) TestParseEnvelopes() {
	table := &gherkin.DataTable{
		Rows: []*gherkin.TableRow{
			{
				Cells: []*gherkin.TableCell{
					{Value: "chainID"},
					{Value: "from"},
				},
			},
			{
				Cells: []*gherkin.TableCell{
					{Value: "888"},
					{Value: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"},
				},
			},
		},
	}

	envelopes, err := s.Context.parser.ParseEnvelopes("", table)
	require.Nil(s.T(), err, "ParseEnvelopes should not error")

	trackers := s.Context.newTrackers(envelopes)
	require.Len(s.T(), trackers, 1, "A tracker should have been created")
	assert.Equal(s.T(), "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4", trackers[0].Current.GetFromString())
}

func (s *ScenarioTestSuite) TestISendEnvelopesToTopic() {
	table := &gherkin.DataTable{
		Rows: []*gherkin.TableRow{
			{
				Cells: []*gherkin.TableCell{
					{Value: "chainID"},
					{Value: "from"},
				},
			},
			{
				Cells: []*gherkin.TableCell{
					{Value: "888"},
					{Value: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"},
				},
			},
		},
	}

	s.producer.ExpectSendMessageAndSucceed()
	err := s.Context.iSendEnvelopesToTopic("tx.crafter", table)
	assert.NoError(s.T(), err)
}

func (s *ScenarioTestSuite) TestEnvelopeShouldBeInTopic() {
	// Prepare trackers
	input := tx.NewBuilder()
	s.Context.setMetadata(input)
	t := s.Context.newTracker(input)
	s.Context.setTrackers([]*tracker.Tracker{t})

	output := tx.NewBuilder().
		SetID(input.GetID()).
		SetContextLabels(input.GetContextLabels())

	err := s.Context.chanReg.Send(LongKeyOf("tx.crafter", output.GetContextLabelsValue("scenario.id"), output.GetID()), output)
	assert.NoError(s.T(), err, "Send in registry should not error")

	err = s.Context.envelopeShouldBeInTopic("tx.crafter")
	assert.NoError(s.T(), err, "envelopeShouldBeInTopic should not error")
	assert.Equal(s.T(), output, s.Context.trackers[0].Current, "Builder on tracker should have been updated")
}

func TestScenarioTestSuite(t *testing.T) {
	suite.Run(t, new(ScenarioTestSuite))
}
