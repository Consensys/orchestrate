package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// EnvelopeStoreTestSuite is a test suit for EnvelopeStore
type EnvelopeStoreTestSuite struct {
	suite.Suite
	Store evlpstore.StoreServer
}

// TestEnvelopeStore test envelope store
func (s *EnvelopeStoreTestSuite) TestStore() {
	// Read / write before storing
	storeResp, err := s.Store.LoadByTxHash(
		context.Background(),
		&evlpstore.TxHashRequest{
			ChainId: chain.CreateChainInt(888),
			TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	)
	assert.NotNil(s.T(), err, "Should error on find envelope by hash")
	assert.True(s.T(), errors.IsNotFoundError(err), "Data should be not found")
	e := errors.FromError(err)
	assert.Contains(s.T(), e.GetComponent(), "envelope-store", "Component should be correct")

	storeResp, err = s.Store.LoadByID(
		context.Background(),
		&evlpstore.IDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
	)
	assert.NotNil(s.T(), err, "Should error on find envelope by ID")
	assert.True(s.T(), errors.IsNotFoundError(err), "LoadByID Data should be not found")
	e = errors.FromError(err)
	assert.Contains(s.T(), e.GetComponent(), "envelope-store", "LoadByID Component should be correct")

	_, err = s.Store.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: "pending",
		},
	)
	assert.NotNil(s.T(), err, "Should error on setStatus")
	assert.True(s.T(), errors.IsNotFoundError(err), "SetStatus Data should be not found")
	e = errors.FromError(err)
	assert.Contains(s.T(), e.GetComponent(), "envelope-store", "SetStatus Component should be correct")

	getStatusResp, err := s.Store.GetStatus(
		context.Background(),
		&evlpstore.IDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
	)
	assert.NotNil(s.T(), err, "Should error on GetStatus")
	assert.True(s.T(), errors.IsNotFoundError(err), "GetStatus Data should be not found")
	e = errors.FromError(err)
	assert.Contains(s.T(), e.GetComponent(), "envelope-store", "GetStatus Component should be correct")

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

	// Store Envelope
	storeResp, err = s.Store.Store(
		context.Background(),
		&evlpstore.StoreRequest{
			Envelope: evlp,
		},
	)
	assert.Nil(s.T(), err, "Should properly store envelope")
	assert.Equal(s.T(), "stored", storeResp.GetStatus(), "Default status should be correct")
	storedAt, _ := ptypes.Timestamp(storeResp.GetLastUpdated())
	assert.True(s.T(), time.Since(storedAt) < 5*time.Second, "Stored date should be close")

	// Load Envelope
	storeResp, err = s.Store.LoadByTxHash(
		context.Background(),
		&evlpstore.TxHashRequest{
			ChainId: chain.CreateChainInt(888),
			TxHash:  "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	)
	assert.Nil(s.T(), err, "Should properly store envelope")
	assert.Equal(s.T(), "stored", storeResp.GetStatus(), "Status should be correct")
	assert.Equal(s.T(), "888", storeResp.GetEnvelope().GetChain().ID().String(), "ChainID should be correct")
	assert.Equal(s.T(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", storeResp.GetEnvelope().GetMetadata().GetId(), "MetadataID should be correct")

	// Set Status
	_, err = s.Store.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: "stored",
		},
	)
	assert.Nil(s.T(), err, "Setting status to %q", "stored")
	_, err = s.Store.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: "error",
		},
	)
	assert.Nil(s.T(), err, "Setting status to %q", "error")
	_, err = s.Store.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: "mined",
		},
	)
	assert.Nil(s.T(), err, "Setting status to %q", "mined")
	_, err = s.Store.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: "pending",
		},
	)
	assert.Nil(s.T(), err, "Setting status to %q", "pending")

	getStatusResp, err = s.Store.GetStatus(
		context.Background(),
		&evlpstore.IDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
	)
	assert.Nil(s.T(), err, "Should not error")
	assert.Equal(s.T(), "pending", getStatusResp.GetStatus(), "Status should be correct")
	sentAt, _ := ptypes.Timestamp(getStatusResp.GetLastUpdated())
	assert.True(s.T(), sentAt.Sub(storedAt) > 0, "Stored should be older than sent date")

	// Stores an already existing
	evlp = &envelope.Envelope{
		Chain:    chain.CreateChainInt(888),
		Metadata: &envelope.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			Raw:  ethereum.HexToData("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"),
			Hash: ethereum.HexToHash("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"),
		},
	}

	// Store Envelope
	storeResp, err = s.Store.Store(
		context.Background(),
		&evlpstore.StoreRequest{
			Envelope: evlp,
		},
	)
	assert.Nil(s.T(), err, "Should update")
	assert.Equal(s.T(), "pending", storeResp.GetStatus(), "Status should be correct")

	// Set status to error
	assert.Nil(s.T(), err, "Setting status to %q", "mined")
	_, err = s.Store.SetStatus(
		context.Background(),
		&evlpstore.SetStatusRequest{
			Id:     "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
			Status: "error",
		},
	)
	assert.Nil(s.T(), err, "Setting status to %q", "error")

	getStatusResp, err = s.Store.GetStatus(
		context.Background(),
		&evlpstore.IDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
	)
	assert.Nil(s.T(), err, "Should not error")
	assert.Equal(s.T(), "error", getStatusResp.GetStatus(), "Status should be correct")
	errorAt, _ := ptypes.Timestamp(getStatusResp.GetLastUpdated())
	assert.True(s.T(), errorAt.Sub(sentAt) > 0, "Stored date should be close")

	// Test to Load By ID
	getStatusResp, err = s.Store.LoadByID(
		context.Background(),
		&evlpstore.IDRequest{
			Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11",
		},
	)
	assert.Nil(s.T(), err, "Should not error")
	assert.Equal(s.T(), "error", getStatusResp.GetStatus(), "Status should be correct")
}

// TestLoadPending test load pending envelopes
func (s *EnvelopeStoreTestSuite) TestLoadPending() {
	for i, chainID := range []int64{1, 2, 3, 12, 42, 888} {
		e := &envelope.Envelope{
			Chain:    chain.CreateChainInt(chainID),
			Metadata: &envelope.Metadata{Id: fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i)},
		}

		_, err := s.Store.Store(
			context.Background(),
			&evlpstore.StoreRequest{
				Envelope: e,
			},
		)
		assert.Nil(s.T(), err, "No error expected")
		// We simulate some exec time between each store
		time.Sleep(100 * time.Millisecond)

		if i%2 == 0 {
			// Every 2 transactions we set status to pending
			_, err := s.Store.SetStatus(
				context.Background(),
				&evlpstore.SetStatusRequest{
					Id:     fmt.Sprintf("a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a1%v", i),
					Status: "pending",
				},
			)
			assert.Nil(s.T(), err, "No error expected")
		}
	}

	loadPendingResp, err := s.Store.LoadPending(context.Background(), &evlpstore.LoadPendingRequest{Duration: 0})
	assert.Nil(s.T(), err, "No error expected on LoadPending")
	assert.Len(s.T(), loadPendingResp.GetEnvelopes(), 3, "Count of envelope pending incorrect")

	loadPendingResp, err = s.Store.LoadPending(context.Background(), &evlpstore.LoadPendingRequest{Duration: (300 * time.Millisecond).Nanoseconds()})
	assert.Nil(s.T(), err, "No error expected on LoadPending")
	assert.Len(s.T(), loadPendingResp.GetEnvelopes(), 2, "Count of envelope pending incorrect")

	loadPendingResp, err = s.Store.LoadPending(context.Background(), &evlpstore.LoadPendingRequest{Duration: (500 * time.Millisecond).Nanoseconds()})
	assert.Nil(s.T(), err, "No error expected on LoadPending")
	assert.Len(s.T(), loadPendingResp.GetEnvelopes(), 1, "Count of envelope pending incorrect")

	loadPendingResp, err = s.Store.LoadPending(context.Background(), &evlpstore.LoadPendingRequest{Duration: (700 * time.Millisecond).Nanoseconds()})
	assert.Nil(s.T(), err, "No error expected on LoadPending")
	assert.Len(s.T(), loadPendingResp.GetEnvelopes(), 0, "Count of envelope pending incorrect")
}
