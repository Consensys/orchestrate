// +build unit

package steps

// TODO: add tests
// import (
// 	"context"
// 	"io/ioutil"
// 	gohttp "net/http"
// 	"strings"
// 	"testing"
//
// 	"github.com/Shopify/sarama/mocks"
// 	gherkin "github.com/cucumber/messages-go/v10"
// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/jwt/generator"
// 	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
// 	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
// 	crc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/client"
// 	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
// 	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/transaction-scheduler/client"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/cucumber/parser"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/tracker"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/e2e/utils"
// )
//
// type ScenarioTestSuite struct {
// 	suite.Suite
// 	Context    *ScenarioContext
// 	chanReg    *chanregistry.ChanRegistry
// 	producer   *mocks.SyncProducer
// 	httpClient *gohttp.Client
// 	crc        svc.ContractRegistryClient
// 	txs        txscheduler.TransactionSchedulerClient
// 	aliasesReg *parser.AliasRegistry
// }
//
// func (s *ScenarioTestSuite) SetupSuite() {
// 	// Set viper configuration
// 	viper.Set(broker.TxCrafterViperKey, "tx-crafter")
//
// 	// Set channel registry
// 	s.chanReg = chanregistry.NewChanRegistry()
// 	s.producer = mocks.NewSyncProducer(s.T(), nil)
// 	s.httpClient = http.NewClient()
// 	s.crc = crc.GlobalClient()
// 	s.txs = txscheduler.GlobalClient()
// 	generator.Init(context.Background())
//
// }
//
// func (s *ScenarioTestSuite) SetupTest() {
// 	s.Context = NewScenarioContext(
// 		s.chanReg,
// 		s.httpClient,
// 		chainregistry.GlobalClient(),
// 		s.crc,
// 		s.txs,
// 		s.producer,
// 		parser.GlobalParser(),
// 		parser.GlobalAliasRegistry(),
// 		generator.GlobalJWTGenerator(),
// 	)
// 	sc := &gherkin.Pickle{}
// 	sc.Name = "stepTable-scenario"
// 	s.Context.init(sc)
// }
//
// func (s *ScenarioTestSuite) TestInitScenarioContext() {
// 	scenario := &gherkin.Pickle{}
// 	scenario.Name = "stepTable-1"
// 	s.Context.init(scenario)
// 	assert.Equal(s.T(), "stepTable-1", s.Context.Pickle.Name, "Context definition should have been set")
// }
//
// func (s *ScenarioTestSuite) TestParseEnvelopes() {
// 	table := &gherkin.PickleStepArgument_PickleTable{
// 		Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{
// 			{
// 				Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
// 					{Value: "chainID"},
// 					{Value: "from"},
// 				},
// 			},
// 			{
// 				Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
// 					{Value: "888"},
// 					{Value: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"},
// 				},
// 			},
// 		},
// 	}
//
// 	envelopes, err := s.Context.parser.ParseEnvelopes(table)
// 	require.Nil(s.T(), err, "ParseEnvelopes should not error")
//
// 	trackers := s.Context.newTrackers(envelopes)
// 	require.Len(s.T(), trackers, 1, "A tracker should have been created")
// 	assert.Equal(s.T(), "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4", trackers[0].Current.GetFromString())
// }
//
// func (s *ScenarioTestSuite) TestISendEnvelopesToTopic() {
// 	table := &gherkin.PickleStepArgument_PickleTable{
// 		Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{
// 			{
// 				Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
// 					{Value: "chainID"},
// 					{Value: "from"},
// 				},
// 			},
// 			{
// 				Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
// 					{Value: "888"},
// 					{Value: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"},
// 				},
// 			},
// 		},
// 	}
//
// 	s.producer.ExpectSendMessageAndSucceed()
// 	err := s.Context.iSendEnvelopesToTopic("tx.crafter", table)
// 	assert.NoError(s.T(), err)
// }
//
// func (s *ScenarioTestSuite) TestEnvelopeShouldBeInTopic() {
// 	// Prepare trackers
// 	input := tx.NewEnvelope()
// 	s.Context.setMetadata(input)
// 	t := s.Context.newTracker(input)
// 	s.Context.setTrackers([]*tracker.Tracker{t})
//
// 	output := tx.NewEnvelope().
// 		SetID(input.GetID()).
// 		SetContextLabels(input.GetContextLabels())
//
// 	err := s.Context.chanReg.Send(utils.LongKeyOf("tx.crafter", output.GetID()), output)
// 	assert.NoError(s.T(), err, "Send in registry should not error")
//
// 	err = s.Context.envelopeShouldBeInTopic("tx.crafter")
// 	assert.NoError(s.T(), err, "envelopeShouldBeInTopic should not error")
// 	assert.Equal(s.T(), output, s.Context.trackers[0].Current, "Envelope on tracker should have been updated")
// }
//
// func (s *ScenarioTestSuite) TestNavJsonResponse() {
// 	// Prepare trackers
// 	rawResp := `[{"uuid":"a8750fc5-4786-4d24-b9fb-5690f6d7c3ac","name":"besu","tenantID":"_"}, {"uuid":"xxx-xxxx","name":"besu2","tenantID":"tenantId2"}]`
//
// 	val, err := navJSONResponse("0.uuid", []byte(rawResp))
// 	assert.NoError(s.T(), err, "navJSONResponse should not error")
// 	assert.Equal(s.T(), "a8750fc5-4786-4d24-b9fb-5690f6d7c3ac", val)
//
// 	val, err = navJSONResponse("0.tenantID", []byte(rawResp))
// 	assert.NoError(s.T(), err, "navJSONResponse should not error")
// 	assert.Equal(s.T(), "_", val)
//
// 	val, err = navJSONResponse("1.name", []byte(rawResp))
// 	assert.NoError(s.T(), err, "navJSONResponse should not error")
// 	assert.Equal(s.T(), "besu2", val)
//
// 	rawResp = `{"jsonrpc": "2.0","id": 1,"result": "5fn2sNAT11mNYDg9gRFeFD1JHmFhoz6Yqd8jsypeq3k="}`
//
// 	val, err = navJSONResponse("result", []byte(rawResp))
// 	assert.NoError(s.T(), err, "navJSONResponse should not error")
// 	assert.Equal(s.T(), "5fn2sNAT11mNYDg9gRFeFD1JHmFhoz6Yqd8jsypeq3k=", val)
// }
//
// func (s *ScenarioTestSuite) TestTheResponseShouldHave() {
// 	table := &gherkin.PickleStepArgument_PickleTable{
// 		Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{
// 			{
// 				Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
// 					{Value: "idempotencyKey"},
// 					{Value: "params.from"},
// 					{Value: "schedule.uuid"},
// 					{Value: "schedule.jobs.0.uuid"},
// 				},
// 			},
// 			{
// 				Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
// 					{Value: "test6"},
// 					{Value: "0x93f7274c9059e601be4512f656b57b830e019e41"},
// 					{Value: "~"},
// 					{Value: "~"},
// 				},
// 			},
// 		},
// 	}
//
// 	stringReader := strings.NewReader(`{"idempotencyKey":"test6","params":{"from":"0x93f7274c9059e601be4512f656b57b830e019e41","methodSignature":"constructor()","to":"0x93f7274c9059e601be4512f656b57b830e019e23"},"schedule":{"uuid":"e2ecf10a-7244-4307-bdea-e734a88d178c","chainUUID":"69bce69b-261d-4e87-8e7f-170bd3527922","jobs":[{"uuid":"48aee833-5d72-4efc-b994-1ae139286557","transaction":{"from":"0x93f7274c9059e601be4512f656b57b830e019e41","to":"0x93f7274c9059e601be4512f656b57b830e019e23"},"status":"STARTED","createdAt":"2020-05-16T20:22:59.304757Z"}],"createdAt":"2020-05-16T20:22:59.304757Z"},"createdAt":"2020-05-16T20:22:59.304757Z"}`)
// 	s.Context.httpResponse = &gohttp.Response{
// 		Body: ioutil.NopCloser(stringReader),
// 	}
// 	err := s.Context.responseShouldHaveFields(table)
// 	assert.NoError(s.T(), err)
// }
//
// func TestScenarioTestSuite(t *testing.T) {
// 	suite.Run(t, new(ScenarioTestSuite))
// }
