// +build integration

package integrationtests

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gofrs/uuid"
	"github.com/gogo/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/txscheduler"

	http2 "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

const (
	waitForEnvelopeTimeOut = 2 * time.Second
)

// txSignerEthereumTestSuite is a test suite for Transaction signer Ethereum
type txSignerEthereumTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func (s *txSignerEthereumTestSuite) TestTxSigner_Ethereum_Public() {
	signature := "0xd35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e01"
	
	s.T().Run("should sign a public ethereum transaction successfully and send it to the sender topic", func(t *testing.T) {
		defer gock.Off()
	
		envelope := fakeEnvelope()
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)
	
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.Equal(t, "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e", retrievedEnvelope.GetRaw())
	})
	
	s.T().Run("should sign a public onetimekey ethereum transaction successfully and send it to the sender topic", func(t *testing.T) {
		defer gock.Off()
	
		envelope := fakeEnvelope()
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)
	
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetRaw())
		assert.NotEqual(
			t,
			"0xf851808398968082520880839896808026a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e",
			retrievedEnvelope.GetRaw(),
		)
	})
	
	s.T().Run("should process a raw transaction and send it to the sender topic", func(t *testing.T) {
		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		_ = envelope.SetRawString("0xf851808398968082520880839896808026a0d35c752d3498e6f5ca1631d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e")
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.Equal(t, envelope.GetRaw(), retrievedEnvelope.GetRaw())
	})

	s.T().Run("should send envelope to tx-recover if an not recoverable error occurs", func(t *testing.T) {
		defer gock.Off()
	
		envelope := fakeEnvelope()
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(http2.StatusUnauthorized).JSON(httputil.ErrorResponse{
			Message: "not authorized requests",
			Code:    666,
		})

		gock.New(txSchedulerURL).Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).Reply(200).JSON(txscheduler.JobResponse{})
	
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.RecoverTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	
		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})
}

func (s *txSignerEthereumTestSuite) TestTxSigner_Ethereum_EEA() {
	signature := "0xd35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e01"

	s.T().Run("should sign a eea transaction successfully and send it to the sender topic", func(t *testing.T) {
		defer gock.Off()

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)

		url := fmt.Sprintf("/ethereum/accounts/%s/sign-eea-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.Equal(t, "0xf8c380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1ea0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", retrievedEnvelope.GetRaw())
	})

	s.T().Run("should sign a eea onetimekey transaction successfully and send it to the sender topic", func(t *testing.T) {
		defer gock.Off()

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)

		url := fmt.Sprintf("/ethereum/accounts/%s/sign-eea-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)
		
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetRaw())
		assert.NotEqual(
			t,
			"0xf8c1808398968082520880839896808026a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1ea0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564",
			retrievedEnvelope.GetRaw(),
		)
	})
}

func (s *txSignerEthereumTestSuite) TestTxSigner_Ethereum_Quorum() {
	signature := "0xd35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e01"

	s.T().Run("should sign a quorum private transaction successfully and send it to the sender topic", func(t *testing.T) {
		defer gock.Off()

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-quorum-private-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.Equal(t, "0xf851808398968082520880839896808026a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e", retrievedEnvelope.GetRaw())
	})

	s.T().Run("should sign a onetimekey quorum private successfully and send it to the sender topic", func(t *testing.T) {
		defer gock.Off()

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-quorum-private-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetRaw())
		assert.NotEqual(
			t,
			"0xf851808398968082520880839896808026a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e",
			retrievedEnvelope.GetRaw(),
		)
	})

	s.T().Run("should not process a tessera private transaction", func(t *testing.T) {
		defer gock.Off()

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX)
		
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-quorum-private-transaction", envelope.GetFromString())
		gock.New(keyManagerURL).Post(url).Reply(200).BodyString(signature)

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.Empty(t, retrievedEnvelope.GetRaw())
	})
}

func (s *txSignerEthereumTestSuite) TestTxSigner_Ethereum_Failures() {
	signature := "0xd35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e01"

	s.T().Run("should not retry a message if the message fails to be decoded", func(t *testing.T) {
		defer gock.Off()

		// First envelope is ignored
		err := s.sendEnvelope(&tx.TxEnvelope{})
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}
	})

	s.T().Run("should retry message on connection error", func(t *testing.T) {
		defer gock.Off()

		// First one fails with service connection error
		envelope := fakeEnvelope()
		url := fmt.Sprintf("/ethereum/accounts/%s/sign-transaction", envelope.GetFromString())

		gock.New(keyManagerURL).Post(url).Times(1).Reply(500).JSON(httputil.ErrorResponse{
			Message: "cannot connect to key-manager service",
			Code:    errors.ServiceConnection,
		})

		gock.New(keyManagerURL).Post(url).Times(1).Reply(200).BodyString(signature)

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.txSignerConfig.SenderTopic, waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetRaw())
	})
}

func (s *txSignerEthereumTestSuite) TestTxSigner_ZHealthCheck() {
	type healthRes struct {
		TransactionScheduler string `json:"transaction-scheduler,omitempty"`
		KeyManager           string `json:"key-manager,omitempty"`
		Kafka                string `json:"kafka,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := s.env.ctx
	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(txSchedulerMetricsURL).Get("/live").Reply(200)
		defer gock.Off()

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), 200, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "OK", status.TransactionScheduler)
		assert.Equal(s.T(), "OK", status.Kafka)
	})

	s.T().Run("should retrieve a negative health check over kafka service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(txSchedulerMetricsURL).Get("/live").Reply(200)
		defer gock.Off()

		// Kill Kafka on first call so data is added in DB and status is CREATED but does not get updated to STARTED
		err = s.env.client.Stop(ctx, kafkaContainerID)
		assert.NoError(t, err)

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		err = s.env.client.StartServiceAndWait(ctx, kafkaContainerID, 10*time.Second)
		assert.NoError(t, err)

		assert.Equal(s.T(), 503, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.NotEqual(s.T(), "OK", status.Kafka)
		assert.Equal(s.T(), "OK", status.TransactionScheduler)
	})
}

func fakeEnvelope() *tx.Envelope {
	jobUUID := uuid.Must(uuid.NewV4()).String()

	envelope := tx.NewEnvelope()
	_ = envelope.SetID(jobUUID)
	_ = envelope.SetJobUUID(jobUUID)
	_ = envelope.SetJobType(tx.JobType_ETH_TX)
	_ = envelope.SetNonce(0)
	_ = envelope.SetFromString("0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5")
	_ = envelope.SetGas(21000)
	_ = envelope.SetGasPriceString("10000000")
	_ = envelope.SetValueString("10000000")
	_ = envelope.SetDataString("0x")
	_ = envelope.SetChainIDString("2888")
	_ = envelope.SetHeadersValue(multitenancy.TenantIDMetadata, "tenantID")
	_ = envelope.SetPrivateFrom("A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=")
	_ = envelope.SetPrivateFor([]string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="})

	return envelope
}

func (s *txSignerEthereumTestSuite) sendEnvelope(protoMessage proto.Message) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = s.env.txSignerConfig.ListenerTopic

	b, err := encoding.Marshal(protoMessage)
	if err != nil {
		return err
	}
	msg.Value = sarama.ByteEncoder(b)

	_, _, err = s.env.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
