// +build integration

package integrationtests

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

// EnvelopeStoreTestSuite is a test suit for EnvelopeStore
type EnvelopeStoreTestSuite struct {
	suite.Suite
	baseURL string
	client  svc.EnvelopeStoreClient
	env     *IntegrationEnvironment
}

func (s *EnvelopeStoreTestSuite) SetupSuite() {
	var err error
	s.client, err = client.DialContextWithDefaultOptions(context.Background(), s.baseURL)
	if err != nil {
		panic(err)
	}
}

func (s *EnvelopeStoreTestSuite) TestEnvelopeStore_StoreEnvelope() {
	ctx := context.Background()
	testEnvelopeUUID := "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"
	var testChainID int64 = 888
	testTxHash := "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"

	s.T().Run("should not load envelope by TxHashes", func(t *testing.T) {
		resp, err := s.client.LoadByTxHashes(ctx, &svc.LoadByTxHashesRequest{
			ChainId:  big.NewInt(testChainID).String(),
			TxHashes: []string{testTxHash},
		})

		assert.NoError(t, err)
		if resp != nil {
			assert.Len(t, resp.Responses, 0)
		}
	})

	s.T().Run("should store a new envelope", func(t *testing.T) {
		// Store Envelope
		b := tx.NewEnvelope().
			SetID(testEnvelopeUUID).
			SetChainID(big.NewInt(testChainID)).
			SetNonce(10).
			SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).
			SetGasPrice(big.NewInt(2000))

		_ = b.SetDataString("0xabcd")
		_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
		_ = b.SetTxHashString(testTxHash)

		resp, err := s.client.Store(ctx, &svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store default status should be correct")
		assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
	})

	s.T().Run("should load envelope by UUID", func(t *testing.T) {
		resp, err := s.client.LoadByID(ctx, &svc.LoadByIDRequest{
			Id: testEnvelopeUUID,
		})

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(testChainID).String(), resp.GetEnvelope().GetChainID(), "ChainID should be the expected")
		assert.Equal(t, testTxHash, resp.GetEnvelope().GetTxHash(), "TxHash should be the expected")
	})

	s.T().Run("should load envelope by TxHash", func(t *testing.T) {
		resp, err := s.client.LoadByTxHash(ctx, &svc.LoadByTxHashRequest{
			ChainId: big.NewInt(testChainID).String(),
			TxHash:  testTxHash,
		})

		assert.NoError(t, err)
		assert.Equal(t, testEnvelopeUUID, resp.GetEnvelope().GetID(), "UUID should be the expected")
	})

	s.T().Run("should load envelope by TxHashes", func(t *testing.T) {
		resp, err := s.client.LoadByTxHashes(ctx, &svc.LoadByTxHashesRequest{
			ChainId:  big.NewInt(testChainID).String(),
			TxHashes: []string{testTxHash},
		})

		assert.NoError(t, err)
		if resp != nil {
			assert.Len(t, resp.Responses, 1)
			assert.Equal(t, testEnvelopeUUID, resp.Responses[0].GetEnvelope().GetID(), "UUID should be the expected")
		}
	})

	s.T().Run("should not load envelope by TxHashes", func(t *testing.T) {
		_, err := s.client.LoadByTxHashes(ctx, &svc.LoadByTxHashesRequest{
			ChainId:  big.NewInt(testChainID).String(),
			TxHashes: []string{},
		})

		assert.Error(t, err)
	})
}

func (s *EnvelopeStoreTestSuite) TestEnvelopeStore_NotFoundAssertion() {
	ctx := context.Background()

	s.T().Run("should assert request to LoadByTxHash of not exiting envelope", func(t *testing.T) {
		_, err := s.client.LoadByTxHash(ctx, &svc.LoadByTxHashRequest{
			ChainId: "888",
			TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a88e534136104b5447a211",
		})

		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err), "should be NotFoundError")
	})

	s.T().Run("should assert LoadByID of not exiting envelope", func(t *testing.T) {
		_, err := s.client.LoadByID(ctx, &svc.LoadByIDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-acde-bd38-0a12",
		})

		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err), "should be NotFoundError")
	})

	s.T().Run("should assert SetStatus of not exiting envelope", func(t *testing.T) {
		_, err := s.client.SetStatus(ctx, &svc.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-acde-6bb9-bd38-0a12",
			Status: svc.Status_PENDING,
		})

		assert.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err), "should be NotFoundError")
	})

	s.T().Run("should assert store a envelope with invalid UUID", func(t *testing.T) {
		// Store Envelope
		b := tx.NewEnvelope().
			SetID("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1-").
			SetChainID(big.NewInt(888)).
			SetNonce(10).
			SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).
			SetGasPrice(big.NewInt(2000))

		_ = b.SetDataString("0xabcd")
		_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
		_ = b.SetTxHashString("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a211")

		_, err := s.client.Store(ctx, &svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		})

		assert.Error(t, err)
		assert.True(t, errors.IsInternalError(err), "should be InternalError")
	})
}

func (s *EnvelopeStoreTestSuite) TestEnvelopeStore_SetStatus() {
	ctx := context.Background()
	testEnvelopeUUID := "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a13"
	var testChainID int64 = 444
	testTxHash := "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a219"

	s.T().Run("should store a new envelope", func(t *testing.T) {
		// Store Envelope
		b := tx.NewEnvelope().
			SetID(testEnvelopeUUID).
			SetChainID(big.NewInt(testChainID)).
			SetNonce(10).
			SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).
			SetGasPrice(big.NewInt(2000))

		_ = b.SetDataString("0xabcd")
		_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
		_ = b.SetTxHashString(testTxHash)

		resp, err := s.client.Store(ctx, &svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store default status should be correct")
		assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
	})

	s.T().Run("should set envelope status to PENDING", func(t *testing.T) {
		resp, err := s.client.SetStatus(ctx, &svc.SetStatusRequest{
			Id:     testEnvelopeUUID,
			Status: svc.Status_PENDING,
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_PENDING, resp.GetStatusInfo().GetStatus(), "SetStatus status should be PENDING")
		assert.True(t, time.Since(resp.GetStatusInfo().SentAtTime()) < 200*time.Millisecond, "Store pending date should be close")
	})

	s.T().Run("should set envelope status to ERROR", func(t *testing.T) {
		resp, err := s.client.SetStatus(ctx, &svc.SetStatusRequest{
			Id:     testEnvelopeUUID,
			Status: svc.Status_ERROR,
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_ERROR, resp.GetStatusInfo().GetStatus(), "SetStatus status should be PENDING")
		assert.True(t, time.Since(resp.GetStatusInfo().SentAtTime()) < 200*time.Millisecond, "Store pending date should be close")
	})

	s.T().Run("should set envelope status to MINED", func(t *testing.T) {
		resp, err := s.client.SetStatus(ctx, &svc.SetStatusRequest{
			Id:     testEnvelopeUUID,
			Status: svc.Status_MINED,
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_MINED, resp.GetStatusInfo().GetStatus(), "SetStatus status should be PENDING")
		assert.True(t, time.Since(resp.GetStatusInfo().SentAtTime()) < 200*time.Millisecond, "Store pending date should be close")
	})

	s.T().Run("should load envelope by UUID", func(t *testing.T) {
		resp, err := s.client.LoadByID(ctx, &svc.LoadByIDRequest{
			Id: testEnvelopeUUID,
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_MINED, resp.GetStatusInfo().GetStatus(), "LoadByID status should be MINED")
		assert.True(t, resp.GetStatusInfo().SentAtTime().Sub(resp.GetStatusInfo().StoredAtTime()) > 0, "Stored should be older than sent date")
	})
}

func (s *EnvelopeStoreTestSuite) TestEnvelopeStore_EnvelopeStoreUpdate() {
	ctx := context.Background()
	testEnvelopeUUID := "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a19"
	testEnvelopeUUIDNew := "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a29"
	var testChainID int64 = 444
	testTxHash := "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a201"
	testTxHashNew := "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a202"

	var b *tx.Envelope
	s.T().Run("should store a new envelope", func(t *testing.T) {
		// Store Envelope
		b = tx.NewEnvelope().
			SetID(testEnvelopeUUID).
			SetChainID(big.NewInt(testChainID)).
			SetNonce(10).
			SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).
			SetGasPrice(big.NewInt(2000))

		_ = b.SetDataString("0xabcd")
		_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
		_ = b.SetTxHashString(testTxHash)

		resp, err := s.client.Store(ctx, &svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store default status should be correct")
		assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
	})

	s.T().Run("should store an already existing envelope UUID with new hash", func(t *testing.T) {
		// Update Envelope
		_ = b.SetTxHashString(testTxHashNew)

		resp, err := s.client.Store(ctx, &svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store status should have been reset to stored")
		assert.Equal(t, testTxHashNew, resp.GetEnvelope().GetTxHash(), "Store hash should have been updated")
		assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
		assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().SentAtTime()), "Store SentAt should have been reset")
		assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().MinedAtTime()), "Store MinedAt should have been reset")
		assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().ErrorAtTime()), "Store ErrorAt should have been reset")
	})

	s.T().Run("should set envelope status to PENDING", func(t *testing.T) {
		resp, err := s.client.SetStatus(ctx, &svc.SetStatusRequest{
			Id:     testEnvelopeUUID,
			Status: svc.Status_PENDING,
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_PENDING, resp.GetStatusInfo().GetStatus(), "SetStatus status should be PENDING")
		assert.True(t, time.Since(resp.GetStatusInfo().SentAtTime()) < 200*time.Millisecond, "Store pending date should be close")
	})

	s.T().Run("should load envelope by new TxHash", func(t *testing.T) {
		resp, err := s.client.LoadByTxHash(ctx, &svc.LoadByTxHashRequest{
			ChainId: big.NewInt(testChainID).String(),
			TxHash:  testTxHashNew,
		})

		assert.NoError(t, err)
		assert.Equal(t, testEnvelopeUUID, resp.GetEnvelope().GetID(), "UUID should be the expected")
	})

	s.T().Run("should store an envelope with new UUID but same Chain and TxHash", func(t *testing.T) {
		// Update Envelope
		_ = b.SetID(testEnvelopeUUIDNew)

		resp, err := s.client.Store(ctx, &svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		})

		assert.NoError(t, err)
		assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store status should have been reset to stored")
		assert.Equal(t, testTxHashNew, resp.GetEnvelope().GetTxHash(), "Store hash should have been updated")
	})

	s.T().Run("should load envelope by new UUID", func(t *testing.T) {
		resp, err := s.client.LoadByID(ctx, &svc.LoadByIDRequest{
			Id: testEnvelopeUUIDNew,
		})

		assert.NoError(t, err)
		assert.Equal(t, testTxHashNew, resp.GetEnvelope().GetTxHash(), "TxHash should be the expected")
	})
}

func (s *EnvelopeStoreTestSuite) TestEnvelopeStore_GetPending() {
	ctx := context.Background()
	testEnvelopeIDTemplete := "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v"
	testStoreInterval := 100 * time.Millisecond

	s.T().Run("should store a bunch of envelopes, set half to PENDING", func(t *testing.T) {
		for i, chainID := range []int64{1, 2, 3, 12, 42, 888} {
			envelopeID := fmt.Sprintf(testEnvelopeIDTemplete, i)
			b := tx.NewEnvelope().SetID(envelopeID).
				SetChainID(big.NewInt(chainID))
			_ = b.SetTxHashString("0x" + RandString(64))

			_, err := s.client.Store(ctx, &svc.StoreRequest{
				Envelope: b.TxEnvelopeAsRequest(),
			})
			assert.NoError(t, err)

			// We simulate some exec time between each store
			time.Sleep(testStoreInterval)

			if i%2 == 0 {
				// Every 2 transactions we set status to pending
				_, err := s.client.SetStatus(ctx, &svc.SetStatusRequest{
					Id:     envelopeID,
					Status: svc.Status_PENDING,
				})
				assert.NoError(t, err)
			}
		}
	})

	s.T().Run("should retrieve three pending envelopes", func(t *testing.T) {
		resp, err := s.client.LoadPending(ctx, &svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(0),
		})
		assert.NoError(t, err)
		assert.Len(t, resp.GetResponses(), 3, "Count of envelope pending incorrect")
	})

	s.T().Run("should retrieve two pending envelopes", func(t *testing.T) {
		resp, err := s.client.LoadPending(ctx, &svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(testStoreInterval * 3),
		})
		assert.NoError(t, err)
		assert.Len(t, resp.GetResponses(), 2, "Count of envelope pending incorrect")
	})

	s.T().Run("should retrieve one pending envelopes", func(t *testing.T) {
		resp, err := s.client.LoadPending(ctx, &svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(testStoreInterval * 5),
		})
		assert.NoError(t, err)
		assert.Len(t, resp.GetResponses(), 1, "Count of envelope pending incorrect")
	})

	s.T().Run("should retrieve none pending envelopes", func(t *testing.T) {
		resp, err := s.client.LoadPending(ctx, &svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(testStoreInterval * 7),
		})
		assert.NoError(t, err)
		assert.Len(t, resp.GetResponses(), 0, "Count of envelope pending incorrect")
	})
}

func RandString(n int) string {
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
