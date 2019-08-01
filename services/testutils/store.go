package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// EnvelopeStoreTestSuite is a test suit for EnvelopeStore
type EnvelopeStoreTestSuite struct {
	suite.Suite
	Store evlpstore.EnvelopeStoreServer
}

func AssertError(t *testing.T, expected string, isError func(err error) bool, err error) {
	assert.NotNil(t, err, "Error should not be nil")
	assert.Contains(t, errors.FromError(err).GetComponent(), expected, "Component should be correct")
	assert.True(t, isError(err), "Error should be from correct class")
}

func (s *EnvelopeStoreTestSuite) AssertLoadByTxHash(
	ctx context.Context, req *evlpstore.LoadByTxHashRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *evlpstore.StoreResponse),
) {
	resp, err := s.Store.LoadByTxHash(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertLoadByID(
	ctx context.Context, req *evlpstore.LoadByIDRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *evlpstore.StoreResponse),
) {
	resp, err := s.Store.LoadByID(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertSetStatus(
	ctx context.Context, req *evlpstore.SetStatusRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *evlpstore.StatusResponse),
) {
	resp, err := s.Store.SetStatus(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertStore(
	ctx context.Context, req *evlpstore.StoreRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *evlpstore.StoreResponse),
) {
	resp, err := s.Store.Store(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertLoadPending(
	ctx context.Context, req *evlpstore.LoadPendingRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *evlpstore.LoadPendingResponse),
) {
	resp, err := s.Store.LoadPending(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

// TestEnvelopeStore test envelope store
func (s *EnvelopeStoreTestSuite) TestStore() {
	// Load envelopes before storing
	s.AssertLoadByTxHash(
		context.Background(),
		&evlpstore.LoadByTxHashRequest{
			Chain:  chain.CreateChainInt(888),
			TxHash: ethereum.HexToHash("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"),
		},
		func(t *testing.T, err error) { AssertError(t, "envelope-store", errors.IsNotFoundError, err) },
		func(t *testing.T, resp *evlpstore.StoreResponse) {},
	)

	s.AssertLoadByID(
		context.Background(),
		&evlpstore.LoadByIDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
		func(t *testing.T, err error) { AssertError(t, "envelope-store", errors.IsNotFoundError, err) },
		func(t *testing.T, resp *evlpstore.StoreResponse) {},
	)

	s.AssertSetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: evlpstore.Status_PENDING,
		},
		func(t *testing.T, err error) { AssertError(t, "envelope-store", errors.IsNotFoundError, err) },
		func(t *testing.T, resp *evlpstore.StatusResponse) {},
	)

	// Store Envelope
	evlp := &envelope.Envelope{
		Chain:    chain.CreateChainInt(888),
		Metadata: &envelope.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				Nonce:    10,
				To:       ethereum.HexToAccount("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"),
				GasPrice: ethereum.IntToQuantity(2000),
				Data:     ethereum.HexToData("0xabcd"),
			},
			Raw:  ethereum.HexToData("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"),
			Hash: ethereum.HexToHash("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"),
		},
	}
	s.AssertStore(
		context.Background(),
		&evlpstore.StoreRequest{
			Envelope: evlp,
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "Store should not error") },
		func(t *testing.T, resp *evlpstore.StoreResponse) {
			assert.Equal(t, evlpstore.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store default status should be correct")
			assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
		},
	)

	// Load Envelope
	s.AssertLoadByTxHash(
		context.Background(),
		&evlpstore.LoadByTxHashRequest{
			Chain:  chain.CreateChainInt(888),
			TxHash: ethereum.HexToHash("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadByTxHash should not error") },
		func(t *testing.T, resp *evlpstore.StoreResponse) {
			assert.Equal(t, evlpstore.Status_STORED, resp.GetStatusInfo().GetStatus(), "LoadByTxHash status should be correct")
			assert.Equal(t, "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", resp.GetEnvelope().GetMetadata().GetId(), "LoadByTxHash Envelope ID should be correct")
			assert.Equal(t, "888", resp.GetEnvelope().GetChain().ID().String(), "LoadByTxHash ChainID should be correct")
		},
	)

	// Set Status
	s.AssertSetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: evlpstore.Status_PENDING,
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "SetStatus should not error") },
		func(t *testing.T, resp *evlpstore.StatusResponse) {
			assert.Equal(t, evlpstore.Status_PENDING, resp.GetStatusInfo().GetStatus(), "SetStatus status should be PENDING")
			assert.True(t, time.Since(resp.GetStatusInfo().SentAtTime()) < 200*time.Millisecond, "Store pending date should be close")
		},
	)

	s.AssertSetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: evlpstore.Status_ERROR,
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "SetStatus should not error") },
		func(t *testing.T, resp *evlpstore.StatusResponse) {
			assert.Equal(t, evlpstore.Status_ERROR, resp.GetStatusInfo().GetStatus(), "SetStatus status should be ERROR")
			assert.True(t, time.Since(resp.GetStatusInfo().ErrorAtTime()) < 200*time.Millisecond, "Store error date should be close")
		},
	)

	s.AssertSetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: evlpstore.Status_MINED,
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "SetStatus should not error") },
		func(t *testing.T, resp *evlpstore.StatusResponse) {
			assert.Equal(t, evlpstore.Status_MINED, resp.GetStatusInfo().GetStatus(), "SetStatus status should be MINED")
			assert.True(t, time.Since(resp.GetStatusInfo().MinedAtTime()) < 200*time.Millisecond, "Store mined date should be close")
		},
	)

	// Load by ID
	s.AssertLoadByID(
		context.Background(),
		&evlpstore.LoadByIDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadByID should not error") },
		func(t *testing.T, resp *evlpstore.StoreResponse) {
			assert.Equal(t, evlpstore.Status_MINED, resp.GetStatusInfo().GetStatus(), "LoadByID status should be MINED")
			assert.True(t, resp.GetStatusInfo().SentAtTime().Sub(resp.GetStatusInfo().StoredAtTime()) > 0, "Stored should be older than sent date")
		},
	)

	// Stores an already existing envelope ID with new hash
	newHash := "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a21a"
	evlp = &envelope.Envelope{
		Chain:    chain.CreateChainInt(888),
		Metadata: &envelope.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			Raw:  ethereum.HexToData("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"),
			Hash: ethereum.HexToHash(newHash),
		},
	}

	time.Sleep(300 * time.Millisecond)
	s.AssertStore(
		context.Background(),
		&evlpstore.StoreRequest{
			Envelope: evlp,
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "Store should update and not error") },
		func(t *testing.T, resp *evlpstore.StoreResponse) {
			assert.Equal(t, evlpstore.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store status should have been reset to stored")
			assert.Equal(t, newHash, resp.GetEnvelope().GetTx().GetHash().Hex(), "Store hash should have been updated")
			assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
			assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().SentAtTime()), "Store SentAt should have been reset")
			assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().MinedAtTime()), "Store MinedAt should have been reset")
			assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().ErrorAtTime()), "Store ErrorAt should have been reset")
		},
	)

	// Load Envelope by TxHash with new hash
	s.AssertLoadByTxHash(
		context.Background(),
		&evlpstore.LoadByTxHashRequest{
			Chain:  chain.CreateChainInt(888),
			TxHash: ethereum.HexToHash(newHash),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadByTxHash should not error") },
		func(t *testing.T, resp *evlpstore.StoreResponse) {
			assert.Equal(t, evlpstore.Status_STORED, resp.GetStatusInfo().GetStatus(), "LoadByTxHash status should be correct")
			assert.Equal(t, "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", resp.GetEnvelope().GetMetadata().GetId(), "LoadByTxHash Envelope ID should be correct")
			assert.Equal(t, "888", resp.GetEnvelope().GetChain().ID().String(), "LoadByTxHash ChainID should be correct")
		},
	)
}

// TestLoadPending test load pending envelopes
func (s *EnvelopeStoreTestSuite) TestLoadPending() {
	for i, chainID := range []int64{1, 2, 3, 12, 42, 888} {
		e := &envelope.Envelope{
			Chain:    chain.CreateChainInt(chainID),
			Metadata: &envelope.Metadata{Id: fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i)},
		}

		_, _ = s.Store.Store(
			context.Background(),
			&evlpstore.StoreRequest{
				Envelope: e,
			},
		)

		// We simulate some exec time between each store
		time.Sleep(100 * time.Millisecond)

		if i%2 == 0 {
			// Every 2 transactions we set status to pending
			_, _ = s.Store.SetStatus(
				context.Background(),
				&evlpstore.SetStatusRequest{
					Id:     fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i),
					Status: evlpstore.Status_PENDING,
				},
			)
		}
	}

	s.AssertLoadPending(
		context.Background(),
		&evlpstore.LoadPendingRequest{
			Duration: utils.DurationToPDuration(0),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadPendinf should not error") },
		func(t *testing.T, resp *evlpstore.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 3, "Count of envelope pending incorrect")
		},
	)

	s.AssertLoadPending(
		context.Background(),
		&evlpstore.LoadPendingRequest{
			Duration: utils.DurationToPDuration(300 * time.Millisecond),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadPendinf should not error") },
		func(t *testing.T, resp *evlpstore.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 2, "Count of envelope pending incorrect")
		},
	)

	s.AssertLoadPending(
		context.Background(),
		&evlpstore.LoadPendingRequest{
			Duration: utils.DurationToPDuration(500 * time.Millisecond),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadPendinf should not error") },
		func(t *testing.T, resp *evlpstore.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 1, "Count of envelope pending incorrect")
		},
	)

	s.AssertLoadPending(
		context.Background(),
		&evlpstore.LoadPendingRequest{
			Duration: utils.DurationToPDuration(700 * time.Millisecond),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadPendinf should not error") },
		func(t *testing.T, resp *evlpstore.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 0, "Count of envelope pending incorrect")
		},
	)
}
