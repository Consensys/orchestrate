// +build unit

package tx

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	error1 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/ethereum"
)

func TestEnvelope(t *testing.T) {
	envelope := &TxRequest{
		Id:    uuid.Must(uuid.NewV4()).String(),
		Chain: "testChain",
		Params: &Params{
			From:            "0x7e654d251da770a068413677967f6d3ea2fea9e4",
			To:              "0xdbb881a51cd4023e4400cef3ef73046743f08da3",
			Gas:             "10",
			GasPrice:        "1089",
			Value:           "56757",
			Nonce:           "10",
			Data:            "0xab",
			Contract:        "test",
			MethodSignature: "constructor()",
		},
	}
	_, err := envelope.Envelope()

	assert.NoError(t, err)
}

func TestRequestToBuilder(t *testing.T) {
	testSet := []struct {
		name            string
		txEnvelope      *TxEnvelope
		expectedBuilder *Envelope
		expectedError   error
	}{
		{
			"tx request without error",
			&TxEnvelope{
				Msg: &TxEnvelope_TxRequest{
					TxRequest: &TxRequest{
						Headers: map[string]string{"testHeader1Key": "testHeader1Value"},
						Chain:   "testChainName",
						Params: &Params{
							From:            "0x7e654d251da770a068413677967f6d3ea2fea9e4",
							To:              "0xdbb881a51cd4023e4400cef3ef73046743f08da3",
							Gas:             "11",
							GasPrice:        "12",
							Value:           "13",
							Nonce:           "14",
							Data:            "0xab",
							Contract:        "testContractName[testContractTag]",
							MethodSignature: "testMethodSignature(string,string)",
							Args:            []string{"test1", "test2"},
							Raw:             "0x02",
							PrivateFor:      []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="},
							PrivateFrom:     "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=",
						},
						Id:            "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
						ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
					},
				},
				InternalLabels: map[string]string{
					ChainIDLabel:      "1",
					TxHashLabel:       "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
					ChainUUIDLabel:    "testChainUUID",
					ScheduleUUIDLabel: "scheduleUUID",
					JobUUIDLabel:      "jobUUID",
				},
			},
			&Envelope{
				Headers: map[string]string{"testHeader1Key": "testHeader1Value"},
				Chain: Chain{
					ChainName: "testChainName",
					ChainUUID: "testChainUUID",
					ChainID:   big.NewInt(1),
				},
				Tx: Tx{
					From:     &(&struct{ x ethcommon.Address }{ethcommon.HexToAddress("0x7e654d251da770a068413677967f6d3ea2fea9e4")}).x,
					To:       &(&struct{ x ethcommon.Address }{ethcommon.HexToAddress("0xdbb881a51cd4023e4400cef3ef73046743f08da3")}).x,
					Gas:      &(&struct{ x uint64 }{11}).x,
					GasPrice: big.NewInt(12),
					Value:    big.NewInt(13),
					Nonce:    &(&struct{ x uint64 }{14}).x,
					Data:     "0xab",
					Raw:      "0x02",
					TxHash:   &(&struct{ x ethcommon.Hash }{ethcommon.HexToHash("0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2")}).x,
				},
				Contract: Contract{
					ContractName:    "testContractName",
					ContractTag:     "testContractTag",
					MethodSignature: "testMethodSignature(string,string)",
					Args:            []string{"test1", "test2"},
				},
				Private: Private{
					PrivateFor:  []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="},
					PrivateFrom: "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=",
				},
				ID:            "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
				ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
				InternalLabels: map[string]string{
					ChainIDLabel:      "1",
					TxHashLabel:       "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
					ChainUUIDLabel:    "testChainUUID",
					ScheduleUUIDLabel: "scheduleUUID",
					JobUUIDLabel:      "jobUUID",
				},
				Errors: make([]*error1.Error, 0),
			},
			nil,
		},
		{
			"tx request with validation error",
			&TxEnvelope{
				Msg: &TxEnvelope_TxRequest{
					TxRequest: &TxRequest{
						Headers: map[string]string{"testHeader1Key": "testHeader1Value"},
						Chain:   "testChainName",
						Params: &Params{
							From:            "0x7e654d251da770a068413677967f6d3ea2fea9e4",
							To:              "0xdbb881a51cd4023e4400cef3ef73046743f08da3",
							Gas:             "11",
							GasPrice:        "12",
							Value:           "13",
							Nonce:           "14",
							Data:            "0xab",
							Contract:        "testContractName[testContractTag]",
							MethodSignature: "testMethodSignature(string,string)",
							Args:            []string{"test1", "test2"},
							PrivateFor:      []string{"not a Base64"},
							PrivateFrom:     "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=",
						},
						Id:            "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
						ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
					},
				},
				InternalLabels: map[string]string{
					ChainIDLabel:   "1",
					ChainUUIDLabel: "testChainUUID",
				},
			},
			nil,
			errors.DataError("[42000@: invalid PrivateFor[0] got not a Base64]"),
		},
		{
			"tx request with invalid internal labels error",
			&TxEnvelope{
				Msg: &TxEnvelope_TxRequest{
					TxRequest: &TxRequest{
						Headers: map[string]string{"testHeader1Key": "testHeader1Value"},
						Chain:   "testChainName",
						Params: &Params{
							From:            "0x7e654d251da770a068413677967f6d3ea2fea9e4",
							To:              "0xdbb881a51cd4023e4400cef3ef73046743f08da3",
							Gas:             "11",
							GasPrice:        "12",
							Value:           "13",
							Nonce:           "14",
							Data:            "0xab",
							Contract:        "testContractName[testContractTag]",
							MethodSignature: "testMethodSignature(string,string)",
							Args:            []string{"test1", "test2"},
							PrivateFor:      []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="},
							PrivateFrom:     "ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=",
						},
						Id:            "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
						ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
					},
				},
				InternalLabels: map[string]string{
					ChainIDLabel:   "@",
					ChainUUIDLabel: "@",
				},
			},
			nil,
			errors.DataError("invalid chainID - got @"),
		},
		{
			"tx response without error",
			&TxEnvelope{
				Msg: &TxEnvelope_TxResponse{
					TxResponse: &TxResponse{
						Headers:       map[string]string{"testHeader1Key": "testHeader1Value"},
						Id:            "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
						JobUUID:       "jobUUID",
						ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
						Transaction: &ethereum.Transaction{
							From:     "0x7e654d251da770a068413677967f6d3ea2fea9e4",
							Nonce:    "14",
							To:       "0xdbb881a51cd4023e4400cef3ef73046743f08da3",
							Value:    "13",
							Gas:      "11",
							GasPrice: "12",
							Data:     "0xab",
							Raw:      "0x02",
							TxHash:   "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
						},
						Errors: []*error1.Error{
							{Message: "testError", Code: 10, Component: "testComponent"},
						},
					},
				},
				InternalLabels: map[string]string{
					ChainIDLabel:   "1",
					TxHashLabel:    "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
					ChainUUIDLabel: "testChainUUID",
				},
			},
			&Envelope{
				Headers: map[string]string{"testHeader1Key": "testHeader1Value"},
				Chain: Chain{
					ChainUUID: "testChainUUID",
					ChainID:   big.NewInt(1),
				},
				Tx: Tx{
					From:     &(&struct{ x ethcommon.Address }{ethcommon.HexToAddress("0x7e654d251da770a068413677967f6d3ea2fea9e4")}).x,
					To:       &(&struct{ x ethcommon.Address }{ethcommon.HexToAddress("0xdbb881a51cd4023e4400cef3ef73046743f08da3")}).x,
					Gas:      &(&struct{ x uint64 }{11}).x,
					GasPrice: big.NewInt(12),
					Value:    big.NewInt(13),
					Nonce:    &(&struct{ x uint64 }{14}).x,
					Data:     "0xab",
					Raw:      "0x02",
					TxHash:   &(&struct{ x ethcommon.Hash }{ethcommon.HexToHash("0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2")}).x,
				},
				ID:            "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
				ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
				InternalLabels: map[string]string{
					ChainIDLabel:      "1",
					TxHashLabel:       "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
					ChainUUIDLabel:    "testChainUUID",
					JobUUIDLabel:      "jobUUID",
					ScheduleUUIDLabel: "",
				},
				Errors: []*error1.Error{
					{Message: "testError", Code: 10, Component: "testComponent"},
				},
			},
			nil,
		},
		{
			"tx request with error",
			&TxEnvelope{
				Msg: &TxEnvelope_TxResponse{
					TxResponse: &TxResponse{
						Headers:       map[string]string{"testHeader1Key": "testHeader1Value"},
						Id:            "envelopID",
						JobUUID:       "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
						ContextLabels: map[string]string{"testContextLabelsKey": "testContextLabelsValue"},
						Transaction: &ethereum.Transaction{
							From:     "error",
							Nonce:    "error",
							To:       "error",
							Value:    "error",
							Gas:      "error",
							GasPrice: "error",
							Data:     "0xab",
							Raw:      "0x02",
							TxHash:   "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
						},
						Errors: []*error1.Error{
							{Message: "testError", Code: 10, Component: "testComponent"},
						},
					},
				},
				InternalLabels: map[string]string{
					ChainIDLabel:   "1",
					TxHashLabel:    "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
					ChainUUIDLabel: "testChainUUID",
				},
			},
			nil,
			errors.DataError("[42000@: invalid gas - got error 42000@: invalid nonce - got error 42000@: invalid gasPrice - got error 42000@: invalid value - got error 42000@: invalid from - got error 42000@: invalid to - got error]"),
		},
		{
			"invalid tx envelope",
			&TxEnvelope{
				InternalLabels: map[string]string{
					ChainIDLabel:      "1",
					TxHashLabel:       "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
					ChainUUIDLabel:    "testChainUUID",
					ScheduleUUIDLabel: "scheduleUUID",
				},
			},
			nil,
			errors.DataError("invalid tx envelope"),
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			b, err := test.txEnvelope.Envelope()
			assert.Equal(t, test.expectedError, err, "Expected same error")
			assert.Equal(t, test.expectedBuilder, b, "Expected same builder")
		})
	}
}

func TestGetID(t *testing.T) {
	envelopeReq := &TxEnvelope{
		Msg: &TxEnvelope_TxRequest{
			TxRequest: &TxRequest{
				Id: "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6",
			},
		},
	}
	assert.Equal(t, "14483d15-d3bf-4aa0-a1ba-1244ba9ef2a6", envelopeReq.GetID())

	envelopeRes := &TxEnvelope{
		Msg: &TxEnvelope_TxResponse{
			TxResponse: &TxResponse{
				Id: "22405248-1ebe-4fef-845d-f91082f4485e",
			},
		},
	}
	assert.Equal(t, "22405248-1ebe-4fef-845d-f91082f4485e", envelopeRes.GetID())

	envelope := &TxEnvelope{}
	assert.Empty(t, envelope.GetID())
}

func TestGetParsedContract(t *testing.T) {
	p := &Params{}
	contractName, contractTag, err := p.GetParsedContract()
	assert.Empty(t, contractName)
	assert.Empty(t, contractTag)
	assert.NoError(t, err)
}

func TestChainID(t *testing.T) {
	e := &TxEnvelope{InternalLabels: make(map[string]string)}
	e.SetChainID(big.NewInt(10))
	assert.Equal(t, "10", e.GetChainID(), "Expected same error")
}

func TestChainUUID(t *testing.T) {
	e := &TxEnvelope{InternalLabels: make(map[string]string)}
	e.SetChainUUID("testChainUUID")
	assert.Equal(t, "testChainUUID", e.GetChainUUID(), "Expected same error")
}

func TestTxHash(t *testing.T) {
	e := &TxEnvelope{InternalLabels: make(map[string]string)}
	e.SetTxHash("0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2")
	assert.Equal(t, "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2", e.GetTxHash(), "Expected same error")
	assert.Equal(t, "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2", e.TxHash().Hex(), "Expected same error")
}

func TestMustGetTxRequest(t *testing.T) {
	e := &TxEnvelope{Msg: &TxEnvelope_TxRequest{TxRequest: &TxRequest{Id: "test"}}}
	assert.Equal(t, e.Msg.(*TxEnvelope_TxRequest).TxRequest, e.MustGetTxRequest(), "Expected same tx request")
}

func TestMustGetTxResponse(t *testing.T) {
	e := &TxEnvelope{Msg: &TxEnvelope_TxResponse{TxResponse: &TxResponse{Id: "test"}}}
	assert.Equal(t, e.Msg.(*TxEnvelope_TxResponse).TxResponse, e.MustGetTxResponse(), "Expected same tx response")
}
