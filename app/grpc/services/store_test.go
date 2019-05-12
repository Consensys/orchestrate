package services

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/context-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

type StoreServiceTestSuite struct {
	suite.Suite
	store *StoreService
}

func (s *StoreServiceTestSuite) SetupTest() {
	s.store = &StoreService{store: mock.NewEnvelopeStore()}
}

func (s *StoreServiceTestSuite) TestStore() {
	req := &store.StoreRequest{Envelope: testEnvelope()}
	_, err := s.store.Store(context.Background(), req)
	assert.Nil(s.T(), err, "Store should not error")
}

func (s *StoreServiceTestSuite) TestLoadByTxHash() {
	// Stores an envelope
	_, err := s.store.Store(context.Background(), &store.StoreRequest{Envelope: testEnvelope()})
	assert.Nil(s.T(), err, "should not error")

	req := &store.TxHashRequest{
		ChainId: "0x3",
		TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
	}
	resp, err := s.store.LoadByTxHash(context.Background(), req)
	assert.Nil(s.T(), err, "LoadByTxHash should not error")
	assert.Equal(s.T(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", resp.GetEnvelope().GetMetadata().GetId(), "Envelope should be correctly loaded")
}

func (s *StoreServiceTestSuite) TestLoadByID() {
	// Stores an envelope
	_, err := s.store.Store(context.Background(), &store.StoreRequest{Envelope: testEnvelope()})
	assert.Nil(s.T(), err, "should not error")

	req := &store.IDRequest{
		Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
	}
	resp, err := s.store.LoadByID(context.Background(), req)
	assert.Nil(s.T(), err, "LoadByID should not error")
	assert.Equal(s.T(), "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210", resp.GetEnvelope().GetTx().GetHash(), "Envelope should be correctly loaded")
}

func (s *StoreServiceTestSuite) TestGetStatus() {
	// Stores an envelope
	_, err := s.store.Store(context.Background(), &store.StoreRequest{Envelope: testEnvelope()})
	assert.Nil(s.T(), err, "should not error")

	req := &store.IDRequest{
		Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
	}
	resp, err := s.store.GetStatus(context.Background(), req)
	assert.Nil(s.T(), err, "GetStatus should not error")
	assert.Equal(s.T(), "stored", resp.GetStatus(), "Status should be correct")
}

func (s *StoreServiceTestSuite) TestSetStatus() {
	// Stores an envelope
	_, err := s.store.Store(context.Background(), &store.StoreRequest{Envelope: testEnvelope()})
	assert.Nil(s.T(), err, "should not error")

	req := &store.SetStatusRequest{
		Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		Status: "pending",
	}

	_, err = s.store.SetStatus(context.Background(), req)
	assert.Nil(s.T(), err, "SetStatus should not error")

	resp, err := s.store.GetStatus(context.Background(), &store.IDRequest{
		Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
	})
	assert.Nil(s.T(), err, "GetStatus should not error")
	assert.Equal(s.T(), "pending", resp.GetStatus(), "Status should be correct")
}

func (s *StoreServiceTestSuite) TestLoadPending() {
	// Stores an envelope and set its status to pending
	_, err := s.store.Store(context.Background(), &store.StoreRequest{Envelope: testEnvelope()})
	assert.Nil(s.T(), err, "should not error")

	_, err = s.store.SetStatus(context.Background(), &store.SetStatusRequest{
		Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		Status: "pending",
	})
	assert.Nil(s.T(), err, "should not error")
	time.Sleep(100 * time.Millisecond)

	req := &store.LoadPendingRequest{
		Duration: (50 * time.Millisecond).Nanoseconds(),
	}

	resp, err := s.store.LoadPending(context.Background(), req)
	assert.Nil(s.T(), err, "LoadPending should not error")
	assert.Len(s.T(), resp.Envelopes, 1, "Expect pending envelopes")

	req = &store.LoadPendingRequest{
		Duration: (200 * time.Millisecond).Nanoseconds(),
	}
	resp, err = s.store.LoadPending(context.Background(), req)
	assert.Nil(s.T(), err, "LoadPending should not error")
	assert.Len(s.T(), resp.Envelopes, 0, "Expect no pending envelopes")
}

func testEnvelope() *envelope.Envelope {
	txData := (&ethereum.TxData{}).
		SetNonce(10).
		SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	return &envelope.Envelope{
		Chain:    &common.Chain{Id: "0x3"},
		Metadata: &envelope.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: txData,
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	}
}

func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(StoreServiceTestSuite))
}
