package types

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestTrace(t *testing.T) {
	tr := NewTrace()

	tr.Chain().ID = big.NewInt(1024)
	if tr.Chain().ID.Text(16) != "400" {
		t.Errorf("Trace: Expected Chain ID %q but got %q", "1024", tr.Chain().ID.Text(400))
	}

	tr.Receiver().ID = "afg"
	if tr.Receiver().ID != "afg" {
		t.Errorf("Trace: Expected Recveiver ID %q but got %q", "afg", tr.Receiver().ID)
	}

	tr.Sender().ID = "fjt"
	if tr.Sender().ID != "fjt" {
		t.Errorf("Trace: Expected Sender ID %q but got %q", "fjt", tr.Sender().ID)
	}

	tr.Call().MethodID = "xyz"
	if tr.Call().MethodID != "xyz" {
		t.Errorf("Trace: Expected Method ID %q but got %q", "xyz", tr.Call().MethodID)
	}

	tr.Tx().SetNonce(10)
	if tr.Tx().Nonce() != 10 {
		t.Errorf("Trace: Expected Nonce %v but got %v", 10, tr.Tx().Nonce())
	}

	tr.Receipt().Status = 1
	if tr.Receipt().Status != 1 {
		t.Errorf("Trace: Expected Status %v but got %v", 1, tr.Receipt().Status)
	}

	tr.Reset()

	if tr.Chain().ID.Text(16) != "0" {
		t.Errorf("Trace: Expected Chain ID %q but got %q", "0", tr.Chain().ID.Text(16))
	}
	if tr.Receiver().ID != "" {
		t.Errorf("Trace: Expected Recveiver ID %q but got %q", "", tr.Receiver().ID)
	}
	if tr.Sender().ID != "" {
		t.Errorf("Trace: Expected Sender ID %q but got %q", "", tr.Sender().ID)
	}
	if tr.Call().MethodID != "" {
		t.Errorf("Trace: Expected Method ID %q but got %q", "", tr.Call().MethodID)
	}
	if tr.Tx().Nonce() != 0 {
		t.Errorf("Trace: Expected Nonce %v but got %v", 0, tr.Tx().Nonce())
	}

	if tr.Receipt().Status != 0 {
		t.Errorf("Trace: Expected Status %v but got %v", 0, tr.Receipt().Status)
	}
}

func TestTraceToString(t *testing.T) {
	testSet := []struct {
		trace         func() *Trace
		expectedTrace map[string]interface{}
	}{
		{
			func() *Trace {
				tr := NewTrace()
				tr.Chain().ID = big.NewInt(1024)
				return tr
			},
			map[string]interface{}{
				"Chain": map[string]interface{}{"ID": "1024"},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Chain().ID = big.NewInt(1024)
				tr.Chain().IsEIP155 = true
				return tr
			},
			map[string]interface{}{
				"Chain": map[string]interface{}{"ID": "1024", "IsEIP155": true},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Sender().ID = "abcd"
				return tr
			},
			map[string]interface{}{
				"Sender": map[string]interface{}{
					"ID": "abcd",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Sender().ID = "abcd"
				tr.Sender().Address = &common.Address{1}
				return tr
			},
			map[string]interface{}{
				"Sender": map[string]interface{}{
					"ID":      "abcd",
					"Address": "0x0100000000000000000000000000000000000000",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receiver().ID = "abcd"
				tr.Receiver().Address = &common.Address{1}
				return tr
			},
			map[string]interface{}{
				"Receiver": map[string]interface{}{
					"ID":      "abcd",
					"Address": "0x0100000000000000000000000000000000000000",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Call().MethodID = "abcd"
				return tr
			},
			map[string]interface{}{
				"Call": map[string]interface{}{
					"MethodID": "abcd",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Call().MethodID = "abcd"
				tr.Call().Args = []string{"a", "b", "c", "d"}
				return tr
			},
			map[string]interface{}{
				"Call": map[string]interface{}{
					"MethodID": "abcd",
					"Args":     []string{"a", "b", "c", "d"},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetNonce(10)
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"txData": map[string]interface{}{
						"Nonce": "10",
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetTo(&common.Address{1})
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"txData": map[string]interface{}{
						"To": "0x0100000000000000000000000000000000000000",
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetValue(big.NewInt(1))
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"txData": map[string]interface{}{
						"Value": "1",
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetGasLimit(uint64(1))
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"txData": map[string]interface{}{
						"GasLimit": "1",
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetGasPrice(big.NewInt(1))
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"txData": map[string]interface{}{
						"GasPrice": "1",
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetData([]byte{1})
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"txData": map[string]interface{}{
						"Data": "0x01",
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetRaw([]byte{1})
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"raw": "0x01",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Tx().SetHash(&common.Hash{1})
				return tr
			},
			map[string]interface{}{
				"Tx": map[string]interface{}{
					"hash": "0x0100000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().PostState = []byte{2}
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"PostState": "0x02",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().Status = uint64(1)
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"Status": "1",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().CumulativeGasUsed = uint64(1)
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"CumulativeGasUsed": "1",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().Bloom.SetBytes([]byte{1})
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"Bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().TxHash.SetBytes([]byte{2})
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"TxHash": "0x0000000000000000000000000000000000000000000000000000000000000002",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().ContractAddress.SetBytes([]byte{2})
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"ContractAddress": "0x0000000000000000000000000000000000000002",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().GasUsed = uint64(4)
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"GasUsed": "4",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().BlockHash.SetBytes([]byte{2})
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"BlockHash": "0x0000000000000000000000000000000000000000000000000000000000000002",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().BlockNumber = uint64(4)
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"BlockNumber": "4",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().TxIndex = uint64(4)
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"TxIndex": "4",
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().Logs = []*Log{
					&Log{
						Log: types.Log{
							Address:     common.Address{1},
							Topics:      []common.Hash{common.Hash{1}},
							Data:        []byte{1},
							BlockNumber: 10,
							TxHash:      common.Hash{1},
							TxIndex:     10,
							BlockHash:   common.Hash{1},
							Index:       10,
							Removed:     false,
						},
					},
				}
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"Logs": []map[string]interface{}{
						map[string]interface{}{
							"Address":     "0x0100000000000000000000000000000000000000",
							"Topics":      []string{"0x0100000000000000000000000000000000000000000000000000000000000000"},
							"Data":        "0x01",
							"BlockNumber": "10",
							"TxHash":      "0x0100000000000000000000000000000000000000000000000000000000000000",
							"TxIndex":     "10",
							"BlockHash":   "0x0100000000000000000000000000000000000000000000000000000000000000",
							"Index":       "10",
							"Removed":     false,
						},
					},
				},
			},
		},
		{
			func() *Trace {
				tr := NewTrace()
				tr.Receipt().Logs = []*Log{
					&Log{
						Log: types.Log{
							Address:     common.Address{1},
							Topics:      []common.Hash{common.Hash{1}},
							Data:        []byte{1},
							BlockNumber: 10,
							TxHash:      common.Hash{1},
							TxIndex:     10,
							BlockHash:   common.Hash{1},
							Index:       10,
							Removed:     false,
						},
						DecodedData: map[string]string{"test": "test"},
					},
				}
				return tr
			},
			map[string]interface{}{
				"Receipt": map[string]interface{}{
					"Logs": []map[string]interface{}{
						map[string]interface{}{
							"Address":     "0x0100000000000000000000000000000000000000",
							"Topics":      []string{"0x0100000000000000000000000000000000000000000000000000000000000000"},
							"Data":        "0x01",
							"BlockNumber": "10",
							"TxHash":      "0x0100000000000000000000000000000000000000000000000000000000000000",
							"TxIndex":     "10",
							"BlockHash":   "0x0100000000000000000000000000000000000000000000000000000000000000",
							"Index":       "10",
							"Removed":     false,
							"DecodedData": map[string]string{"test": "test"},
						},
					},
				},
			},
		},
	}

	for i, test := range testSet {
		if !reflect.DeepEqual(test.expectedTrace, test.trace().String()) {
			t.Errorf("Trace (%d/%d): Expected Trace String %v but got %v", i+1, len(testSet), test.expectedTrace, test.trace().String())
		}
	}

}
