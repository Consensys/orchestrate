package testutils

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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

// EnvelopeStoreTestSuite is a test suit for EnvelopeStore
type EnvelopeStoreTestSuite struct {
	suite.Suite
	Store svc.EnvelopeStoreServer
}

func AssertError(t *testing.T, expected string, isError func(err error) bool, err error) {
	assert.Error(t, err, "Error should not be nil")
	assert.True(t, isError(err), "Error should be from correct class")
}

func (s *EnvelopeStoreTestSuite) AssertLoadByTxHash(
	ctx context.Context, req *svc.LoadByTxHashRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *svc.StoreResponse),
) {
	resp, err := s.Store.LoadByTxHash(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertLoadByID(
	ctx context.Context, req *svc.LoadByIDRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *svc.StoreResponse),
) {
	resp, err := s.Store.LoadByID(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertSetStatus(
	ctx context.Context, req *svc.SetStatusRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *svc.StatusResponse),
) {
	resp, err := s.Store.SetStatus(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertStore(
	ctx context.Context, req *svc.StoreRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *svc.StoreResponse),
) {
	resp, err := s.Store.Store(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

func (s *EnvelopeStoreTestSuite) AssertLoadPending(
	ctx context.Context, req *svc.LoadPendingRequest,
	assertErr func(t *testing.T, err error),
	assertResp func(t *testing.T, resp *svc.LoadPendingResponse),
) {
	resp, err := s.Store.LoadPending(ctx, req)
	assertErr(s.T(), err)
	assertResp(s.T(), resp)
}

// TestEnvelopeStore test envelope store
func (s *EnvelopeStoreTestSuite) TestStore() {
	ctx := multitenancy.WithTenantID(context.Background(), multitenancy.DefaultTenantIDName)
	// Load envelopes before storing
	s.AssertLoadByTxHash(
		ctx,
		&svc.LoadByTxHashRequest{
			ChainId: big.NewInt(888).String(),
			TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
		func(t *testing.T, err error) { AssertError(t, "envelope-store", errors.IsNotFoundError, err) },
		func(t *testing.T, resp *svc.StoreResponse) {},
	)

	s.AssertLoadByID(
		ctx,
		&svc.LoadByIDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
		func(t *testing.T, err error) { AssertError(t, "envelope-store", errors.IsNotFoundError, err) },
		func(t *testing.T, resp *svc.StoreResponse) {},
	)

	s.AssertSetStatus(
		ctx,
		&svc.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: svc.Status_PENDING,
		},
		func(t *testing.T, err error) { AssertError(t, "envelope-store", errors.IsNotFoundError, err) },
		func(t *testing.T, resp *svc.StatusResponse) {},
	)

	// Store Envelope
	b := tx.NewEnvelope().SetID("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11").SetChainID(big.NewInt(888)).SetNonce(10).SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).SetGasPrice(big.NewInt(2000))
	_ = b.SetDataString("0xabcd")
	_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
	_ = b.SetTxHashString("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210")

	s.AssertStore(
		ctx,
		&svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "Store should not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store default status should be correct")
			assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
		},
	)

	// Load Envelope
	s.AssertLoadByTxHash(
		ctx,
		&svc.LoadByTxHashRequest{
			ChainId: big.NewInt(888).String(),
			TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadByTxHash should not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "LoadByTxHash status should be correct")
			assert.Equal(t, "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", resp.GetEnvelope().GetID(), "LoadByTxHash Envelope UUID should be correct")
			assert.Equal(t, "888", resp.GetEnvelope().GetChainID(), "LoadByTxHash GetBigChainID should be correct")
		},
	)

	// Set Status
	s.AssertSetStatus(
		ctx,
		&svc.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: svc.Status_PENDING,
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "SetStatus should not error") },
		func(t *testing.T, resp *svc.StatusResponse) {
			assert.Equal(t, svc.Status_PENDING, resp.GetStatusInfo().GetStatus(), "SetStatus status should be PENDING")
			assert.True(t, time.Since(resp.GetStatusInfo().SentAtTime()) < 200*time.Millisecond, "Store pending date should be close")
		},
	)

	s.AssertSetStatus(
		ctx,
		&svc.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: svc.Status_ERROR,
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "SetStatus should not error") },
		func(t *testing.T, resp *svc.StatusResponse) {
			assert.Equal(t, svc.Status_ERROR, resp.GetStatusInfo().GetStatus(), "SetStatus status should be ERROR")
			assert.True(t, time.Since(resp.GetStatusInfo().ErrorAtTime()) < 200*time.Millisecond, "Store error date should be close")
		},
	)

	s.AssertSetStatus(
		ctx,
		&svc.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: svc.Status_MINED,
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "SetStatus should not error") },
		func(t *testing.T, resp *svc.StatusResponse) {
			assert.Equal(t, svc.Status_MINED, resp.GetStatusInfo().GetStatus(), "SetStatus status should be MINED")
			assert.True(t, time.Since(resp.GetStatusInfo().MinedAtTime()) < 200*time.Millisecond, "Store mined date should be close")
		},
	)

	// Load by UUID
	s.AssertLoadByID(
		ctx,
		&svc.LoadByIDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadByID should not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, svc.Status_MINED, resp.GetStatusInfo().GetStatus(), "LoadByID status should be MINED")
			assert.True(t, resp.GetStatusInfo().SentAtTime().Sub(resp.GetStatusInfo().StoredAtTime()) > 0, "Stored should be older than sent date")
		},
	)

	// Stores an already existing envelope UUID with new hash
	newHash := "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a21a"
	b = tx.NewEnvelope().SetID("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11").SetChainID(big.NewInt(888)).SetNonce(10).SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).SetGasPrice(big.NewInt(2000))
	_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
	_ = b.SetTxHashString(newHash)

	time.Sleep(300 * time.Millisecond)
	s.AssertStore(
		ctx,
		&svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "Store should update and not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store status should have been reset to stored")
			assert.Equal(t, newHash, resp.GetEnvelope().GetTxHash(), "Store hash should have been updated")
			assert.True(t, time.Since(resp.GetStatusInfo().StoredAtTime()) < 200*time.Millisecond, "Store stored date should be close")
			assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().SentAtTime()), "Store SentAt should have been reset")
			assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().MinedAtTime()), "Store MinedAt should have been reset")
			assert.True(t, time.Time{}.Equal(resp.GetStatusInfo().ErrorAtTime()), "Store ErrorAt should have been reset")
		},
	)

	// Load Envelope by TxHash with new hash
	s.AssertLoadByTxHash(
		ctx,
		&svc.LoadByTxHashRequest{
			ChainId: big.NewInt(888).String(),
			TxHash:  newHash,
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadByTxHash should not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "LoadByTxHash status should be correct")
			assert.Equal(t, "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", resp.GetEnvelope().GetID(), "LoadByTxHash Envelope UUID should be correct")
			assert.Equal(t, "888", resp.GetEnvelope().GetChainID(), "LoadByTxHash GetBigChainID should be correct")
		},
	)

	// Stores an envelope with new UUID but same Chain and TxHash
	newID := "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-9b53"
	b = tx.NewEnvelope().SetID(newID).SetChainID(big.NewInt(888)).SetNonce(10).SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62")).SetGasPrice(big.NewInt(2000))
	_ = b.SetRawString("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
	_ = b.SetTxHashString(newHash)

	s.AssertStore(
		ctx,
		&svc.StoreRequest{
			Envelope: b.TxEnvelopeAsRequest(),
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "Store should not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, svc.Status_STORED, resp.GetStatusInfo().GetStatus(), "Store status should have been reset to stored")
		},
	)

	// Load by UUID
	s.AssertLoadByID(
		ctx,
		&svc.LoadByIDRequest{
			Id: newID,
		},
		func(t *testing.T, err error) { assert.Nil(t, err, "LoadByID should not error") },
		func(t *testing.T, resp *svc.StoreResponse) {
			assert.Equal(t, "888", resp.GetEnvelope().GetChainID(), "LoadByID GetBigChainID should be correct")
			assert.Equal(t, newHash, resp.GetEnvelope().GetTxHash(), "Store hash should have been updated")
		},
	)

}

var letterRunes = []rune("abcdef0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// TestLoadPending test load pending envelopes
func (s *EnvelopeStoreTestSuite) TestLoadPending() {
	ctx := multitenancy.WithTenantID(context.Background(), multitenancy.DefaultTenantIDName)

	for i, chainID := range []int64{1, 2, 3, 12, 42, 888} {
		b := tx.NewEnvelope().SetID(fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i)).
			SetChainID(big.NewInt(chainID))
		_ = b.SetTxHashString("0x" + RandString(64))

		_, _ = s.Store.Store(
			ctx,
			&svc.StoreRequest{
				Envelope: b.TxEnvelopeAsRequest(),
			},
		)

		// We simulate some exec time between each store
		time.Sleep(100 * time.Millisecond)

		if i%2 == 0 {
			// Every 2 transactions we set status to pending
			_, _ = s.Store.SetStatus(
				ctx,
				&svc.SetStatusRequest{
					Id:     fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i),
					Status: svc.Status_PENDING,
				},
			)
		}
	}

	s.AssertLoadPending(
		ctx,
		&svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(0),
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadPending should not error") },
		func(t *testing.T, resp *svc.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 3, "Count of envelope pending incorrect")
		},
	)

	s.AssertLoadPending(
		ctx,
		&svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(300 * time.Millisecond),
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadPending should not error") },
		func(t *testing.T, resp *svc.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 2, "Count of envelope pending incorrect")
		},
	)

	s.AssertLoadPending(
		ctx,
		&svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(500 * time.Millisecond),
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadPending should not error") },
		func(t *testing.T, resp *svc.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 1, "Count of envelope pending incorrect")
		},
	)

	s.AssertLoadPending(
		ctx,
		&svc.LoadPendingRequest{
			Duration: utils.DurationToPDuration(700 * time.Millisecond),
		},
		func(t *testing.T, err error) { assert.NoError(t, err, "LoadPending should not error") },
		func(t *testing.T, resp *svc.LoadPendingResponse) {
			assert.Len(t, resp.GetResponses(), 0, "Count of envelope pending incorrect")
		},
	)
}
