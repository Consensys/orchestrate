package steps

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Shopify/sarama/mocks"
	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/service/chanregistry"
)

var testAddress = ethcommon.HexToAddress("0x00").Hex()

func TestFeatureContext(t *testing.T) {
	s := &godog.Suite{}
	FeatureContext(s)
}

type ScenarioTestSuite struct {
	suite.Suite
	Scenario *ScenarioContext
}

func (s *ScenarioTestSuite) SetupTest() {
	s.Scenario = &ScenarioContext{
		Logger: log.StandardLogger().WithFields(log.Fields{
			"Sceneario": "test",
		}),
	}
	s.Scenario.EnvelopesChan = make(map[string]chan *envelope.Envelope)
	s.Scenario.Envelopes = make(map[string]*envelope.Envelope)
	s.Scenario.Value = make(map[string]interface{})
	s.Scenario.Value[testAddress] = testAddress

	viper.Set("cucumber.steps.timeout", 1)
	viper.Set("cucumber.steps.miningtimeout", 1)
}

func (s *ScenarioTestSuite) TestInitScenarioContext() {
	c := chanregistry.NewChanRegistry()
	chanregistry.SetGlobalChanRegistry(c)

	scenario := &gherkin.Scenario{}
	scenario.Name = "test"
	s.Scenario.initScenarioContext(scenario)

	assert.NotEmpty(s.T(), s.Scenario.ScenarioID, "Should not be empty")
	assert.Len(s.T(), s.Scenario.EnvelopesChan, 5, "Should not be empty")

	scenarioOutline := &gherkin.ScenarioOutline{}
	scenarioOutline.Name = "test"
	s.Scenario.initScenarioContext(scenarioOutline)

	assert.NotEmpty(s.T(), s.Scenario.ScenarioID, "Should not be empty")
	assert.Len(s.T(), s.Scenario.EnvelopesChan, 5, "Should not be empty")
}

func (s *ScenarioTestSuite) TestIHaveTheFollowingEnvelope() {

	rawEnvelopes := &gherkin.DataTable{
		Rows: []*gherkin.TableRow{
			&gherkin.TableRow{
				Cells: []*gherkin.TableCell{
					&gherkin.TableCell{
						Value: "chainId",
					},
					&gherkin.TableCell{
						Value: "from",
					},
					&gherkin.TableCell{
						Value: "AliasTo",
					},
				},
			},
			&gherkin.TableRow{
				Cells: []*gherkin.TableCell{
					&gherkin.TableCell{
						Value: "888",
					},
					&gherkin.TableCell{
						Value: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
					},
					&gherkin.TableCell{
						Value: testAddress,
					},
				},
			},
		},
	}

	_ = s.Scenario.iHaveTheFollowingEnvelope(rawEnvelopes)
}

func (s *ScenarioTestSuite) TestISendTheseEnvelopeToCoreStack() {

	producer := mocks.NewSyncProducer(s.T(), nil)
	producer.ExpectSendMessageAndSucceed()
	producer.ExpectSendMessageAndFail(sarama.ErrOutOfBrokers)
	broker.SetGlobalSyncProducer(producer)

	s.Scenario.Envelopes["test"] = &envelope.Envelope{
		Chain: chain.CreateChainInt(888),
	}

	err := s.Scenario.iSendTheseEnvelopeToCoreStack()
	assert.NoError(s.T(), err, "Should not get an error")
	err = s.Scenario.iSendTheseEnvelopeToCoreStack()
	assert.Error(s.T(), err, "Should get an error")
}

func (s *ScenarioTestSuite) TestCoreStackShouldReceiveEnvelopes() {

	mockChan := make(chan *envelope.Envelope)
	s.Scenario.EnvelopesChan[viper.GetString("kafka.topic.crafter")] = mockChan

	testEnvelope := &envelope.Envelope{
		Chain: chain.CreateChainInt(888),
	}
	s.Scenario.Envelopes["test"] = testEnvelope

	var err error

	go func() {
		mockChan <- testEnvelope
	}()

	// Testing the well functioning of the step
	err = s.Scenario.coreStackShouldReceiveEnvelopes()
	assert.NoError(s.T(), err, "Should not get an error")

	// Test for not receiving envelopes before timeout
	err = s.Scenario.coreStackShouldReceiveEnvelopes()
	assert.Error(s.T(), err, "Should get an error")
}

func (s *ScenarioTestSuite) TestTheTxcrafterShouldSetTheData() {

	mockChan := make(chan *envelope.Envelope)
	s.Scenario.EnvelopesChan[viper.GetString("kafka.topic.nonce")] = mockChan

	testEnvelope := &envelope.Envelope{
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Data: ethereum.HexToData("0x00"),
			},
		},
	}
	s.Scenario.Envelopes["test"] = testEnvelope

	var err error

	go func() {
		mockChan <- testEnvelope
	}()

	// Test the well functioning of the step with expected envelopes
	err = s.Scenario.theTxcrafterShouldSetTheData()
	assert.NoError(s.T(), err, "Should not get an error")

	// Test for not receiving envelopes before timeout
	err = s.Scenario.theTxcrafterShouldSetTheData()
	assert.Error(s.T(), err, "Should get an error")

	// Test step with unexpected envelopes
	unexpectedEnvelope := &envelope.Envelope{
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Gas: uint64(10),
			},
		},
	}

	go func() {
		mockChan <- unexpectedEnvelope
	}()

	err = s.Scenario.theTxcrafterShouldSetTheData()
	assert.Error(s.T(), err, "Should not get an error")
}

func (s *ScenarioTestSuite) TestTheTxnonceShouldSetTheNonce() {

	mockChan := make(chan *envelope.Envelope)
	s.Scenario.EnvelopesChan[viper.GetString("kafka.topic.signer")] = mockChan

	addr := [][]byte{
		ethcommon.HexToAddress("0x00").Bytes(),
		ethcommon.HexToAddress("0x01").Bytes(),
	}
	for i := range make([]int, 10) {
		s.Scenario.Envelopes[fmt.Sprintf("%s-%d", "test", i)] = &envelope.Envelope{
			Chain: chain.CreateChainInt(888),
			From:  ethereum.NewAccount(addr[i%len(addr)]),
			Tx: &ethereum.Transaction{
				TxData: &ethereum.TxData{
					Nonce: uint64(i),
				},
			},
		}
	}

	var err error

	// Test the well functioning of the step with expected envelopes
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = s.Scenario.theTxnonceShouldSetTheNonce()
		wg.Done()
	}()

	for _, e := range s.Scenario.Envelopes {
		go func(e *envelope.Envelope) {
			mockChan <- e
		}(e)
	}
	wg.Wait()
	assert.NoError(s.T(), err, "Should not get an error")

	// Test for not receiving envelopes before timeout
	err = s.Scenario.theTxnonceShouldSetTheNonce()
	assert.Error(s.T(), err, "Should get an error")

	// Test step with unexpected envelopes
	s.Scenario.Envelopes = map[string]*envelope.Envelope{
		"unexpected1": &envelope.Envelope{
			From: ethereum.NewAccount(addr[0]),
			Tx: &ethereum.Transaction{
				TxData: &ethereum.TxData{
					Nonce: 10,
				},
			},
		},
		"unexpected2": &envelope.Envelope{
			From: ethereum.NewAccount(addr[0]),
			Tx: &ethereum.Transaction{
				TxData: &ethereum.TxData{
					Nonce: 10,
				},
			},
		},
	}

	wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = s.Scenario.theTxnonceShouldSetTheNonce()
		wg.Done()
	}()

	for _, e := range s.Scenario.Envelopes {
		go func(e *envelope.Envelope) {
			mockChan <- e
		}(e)
	}
	wg.Wait()
	assert.Error(s.T(), err, "Should not get an error")
}

func (s *ScenarioTestSuite) TestTheTxsignerShouldSign() {

	mockChan := make(chan *envelope.Envelope)
	s.Scenario.EnvelopesChan[viper.GetString("kafka.topic.sender")] = mockChan

	testEnvelope := &envelope.Envelope{
		Tx: &ethereum.Transaction{
			Raw: ethereum.HexToData("0x00"),
		},
	}
	s.Scenario.Envelopes["test"] = testEnvelope

	var err error

	go func() {
		mockChan <- testEnvelope
	}()
	// Test the well functioning of the step with expected envelopes
	err = s.Scenario.theTxsignerShouldSign()
	assert.NoError(s.T(), err, "Should not get an error")

	// Test for not receiving envelopes before timeout
	err = s.Scenario.theTxsignerShouldSign()
	assert.Error(s.T(), err, "Should get an error")

	// Test step with unexpected envelopes
	unexpectedEnvelope := &envelope.Envelope{
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Gas: uint64(10),
			},
		},
	}

	go func() {
		mockChan <- unexpectedEnvelope
	}()
	err = s.Scenario.theTxsignerShouldSign()
	assert.Error(s.T(), err, "Should not get an error")
}

func (s *ScenarioTestSuite) TestTheTxlistenerShouldCatchTheTx() {

	chainIds := []int64{1, 2}

	for _, v := range chainIds {
		topic := fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), v)
		s.Scenario.EnvelopesChan[topic] = make(chan *envelope.Envelope)
	}

	for i := range make([]int, 10) {
		s.Scenario.Envelopes[fmt.Sprintf("%s-%d", "test", i)] = &envelope.Envelope{
			Chain: chain.CreateChainInt(chainIds[i%len(chainIds)]),
			Receipt: &ethereum.Receipt{
				TxHash: ethereum.CreateHash(ethcommon.HexToAddress("0x00").Bytes()),
			},
		}
	}

	var err error

	// Test the well functioning of the step with expected envelopes
	common.InParallel(
		func() {
			err = s.Scenario.theTxlistenerShouldCatchTheTx()
		},
		func() {
			for _, e := range s.Scenario.Envelopes {
				go func(e *envelope.Envelope) {
					topic := fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), e.GetChain().ID().Int64())
					s.Scenario.EnvelopesChan[topic] <- e
				}(e)
			}
		},
	)
	assert.NoError(s.T(), err, "Should not get an error")
}

func (s *ScenarioTestSuite) TestTheTxdecoderShouldDecode() {

	mockChan := make(chan *envelope.Envelope)
	s.Scenario.EnvelopesChan[viper.GetString("kafka.topic.decoded")] = mockChan

	decoded := make(map[string]string)
	decoded["test"] = "test"
	testEnvelope := &envelope.Envelope{
		Receipt: &ethereum.Receipt{
			TxHash: ethereum.CreateHash(ethcommon.HexToAddress("0x00").Bytes()),
			Logs: []*ethereum.Log{
				&ethereum.Log{
					Topics: []*ethereum.Hash{
						ethereum.CreateHash(ethcommon.HexToAddress("0x00").Bytes()),
						ethereum.CreateHash(ethcommon.HexToAddress("0x01").Bytes()),
					},
					DecodedData: decoded,
				},
			},
		},
	}
	s.Scenario.Envelopes["test"] = testEnvelope

	var err error

	go func() {
		mockChan <- testEnvelope
	}()
	// Test the well functioning of the step with expected envelopes
	err = s.Scenario.theTxdecoderShouldDecode()
	assert.NoError(s.T(), err, "Should not get an error")

	// Test for not receiving envelopes before timeout
	err = s.Scenario.theTxdecoderShouldDecode()
	assert.Error(s.T(), err, "Should get an error")

	// Test step with unexpected envelopes
	unexpectedEnvelope := &envelope.Envelope{
		Receipt: &ethereum.Receipt{
			TxHash: ethereum.CreateHash(ethcommon.HexToAddress("0x00").Bytes()),
			Logs: []*ethereum.Log{
				&ethereum.Log{
					Topics: []*ethereum.Hash{
						ethereum.CreateHash(ethcommon.HexToAddress("0x00").Bytes()),
						ethereum.CreateHash(ethcommon.HexToAddress("0x01").Bytes()),
					},
					DecodedData: make(map[string]string),
				},
			},
		},
	}

	go func() {
		mockChan <- unexpectedEnvelope
	}()
	err = s.Scenario.theTxdecoderShouldDecode()
	assert.Error(s.T(), err, "Should not get an error")
}

func (s *ScenarioTestSuite) TestBeforeStep() {
	step := &gherkin.Step{
		Text: "Test",
	}
	s.Scenario.beforeStep(step)
}

func TestScenarioTestSuite(t *testing.T) {
	suite.Run(t, new(ScenarioTestSuite))
}
