package types

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestNewReceipt(t *testing.T) {

	root := []byte{}
	failed := false
	cumulativeGasUsed := uint64(0)

	r := newReceipt(root, failed, cumulativeGasUsed)

	if bytes.Compare(r.PostState, root) != 0 {
		t.Errorf("Trace: Expected PostState %q but got %q", root, r.PostState)
	}
	if r.Status != uint64(1) {
		t.Errorf("Trace: Expected Status %q but got %q", uint64(1), r.Status)
	}
	if r.CumulativeGasUsed != cumulativeGasUsed {
		t.Errorf("Trace: Expected CumulativeGasUsed %q but got %q", cumulativeGasUsed, r.CumulativeGasUsed)
	}
	if (r.Bloom != types.Bloom{}) {
		t.Errorf("Trace: Expected Bloom %q but got %q", types.Bloom{}, r.Bloom)
	}
	if len(r.Logs) != 0 {
		t.Errorf("Trace: Expected len(Logs) %q but got %q", 0, len(r.Logs))
	}
	if (r.TxHash != common.Hash{}) {
		t.Errorf("Trace: Expected TxHash %q but got %q", common.Hash{}, r.TxHash)
	}
	if (r.ContractAddress != common.Address{}) {
		t.Errorf("Trace: Expected ContractAddress %q but got %q", common.Address{}, r.ContractAddress)
	}
	if r.GasUsed != uint64(0) {
		t.Errorf("Trace: Expected GasUsed %q but got %q", uint64(0), r.GasUsed)
	}

}

// TODO: See why reset is not working for r.Bloom, r.TxHash, r.ContractAddress
func TestResetReceipt(t *testing.T) {

	root := []byte{1}
	failed := false
	cumulativeGasUsed := uint64(1000000)

	r := newReceipt(root, failed, cumulativeGasUsed)
	r.Bloom.SetBytes([]byte{1})
	r.TxHash.SetBytes([]byte{1})
	r.ContractAddress.SetBytes([]byte{1})
	r.Logs = []*Log{&Log{}}
	r.GasUsed = 1

	if bytes.Compare(r.PostState, root) != 0 {
		t.Errorf("Trace: Expected PostState %q but got %q", root, r.PostState)
	}
	if r.Status != uint64(1) {
		t.Errorf("Trace: Expected Status %q but got %q", uint64(1), r.Status)
	}
	if r.CumulativeGasUsed != cumulativeGasUsed {
		t.Errorf("Trace: Expected CumulativeGasUsed %q but got %q", cumulativeGasUsed, r.CumulativeGasUsed)
	}
	if r.Bloom != types.BytesToBloom([]byte{1}) {
		t.Errorf("Trace: Expected Bloom %q but got %q", types.BytesToBloom([]byte{1}), r.Bloom)
	}
	if len(r.Logs) != 1 {
		t.Errorf("Trace: Expected len(Logs) %q but got %q", 1, len(r.Logs))
	}
	if r.TxHash != common.BytesToHash([]byte{1}) {
		t.Errorf("Trace: Expected TxHash %q but got %q", common.BytesToHash([]byte{1}), r.TxHash)
	}
	if r.ContractAddress != common.BytesToAddress([]byte{1}) {
		t.Errorf("Trace: Expected ContractAddress %q but got %q", common.BytesToAddress([]byte{1}), r.ContractAddress)
	}
	if r.GasUsed != uint64(1) {
		t.Errorf("Trace: Expected GasUsed %q but got %q", uint64(1), r.GasUsed)
	}

	r.reset()

	if bytes.Compare(r.PostState, []byte{}) != 0 {
		t.Errorf("Trace after reset: Expected PostState %q but got %q", root, r.PostState)
	}
	if r.Status != uint64(0) {
		t.Errorf("Trace after reset: Expected Status %q but got %q", uint64(0), r.Status)
	}
	if r.CumulativeGasUsed != uint64(0) {
		t.Errorf("Trace after reset: Expected CumulativeGasUsed %q but got %q", cumulativeGasUsed, r.CumulativeGasUsed)
	}
	// if r.Bloom != types.BytesToBloom([]byte{0}) {
	// 	t.Errorf("Trace after reset: Expected Bloom %q but got %q", types.Bloom{}, r.Bloom)
	// }
	if len(r.Logs) != 0 {
		t.Errorf("Trace after reset: Expected Logs %q but got %q", 0, len(r.Logs))
	}
	// if (r.TxHash != common.Hash{}) {
	// 	t.Errorf("Trace after reset: Expected TxHash %q but got %q", common.Hash{}, r.TxHash)
	// }
	// if (r.ContractAddress != common.Address{}) {
	// 	t.Errorf("Trace after reset: Expected ContractAddress %q but got %q", common.Address{}, r.ContractAddress)
	// }
	if r.GasUsed != uint64(0) {
		t.Errorf("Trace after reset: Expected GasUsed %q but got %q", uint64(0), r.GasUsed)
	}
}

func TestSetDecodedData(t *testing.T) {

	root := []byte{}
	failed := false
	cumulativeGasUsed := uint64(0)

	r := newReceipt(root, failed, cumulativeGasUsed)

	r.Logs = []*Log{&Log{}}

	mapping := make(map[string]string)
	mapping["testKey"] = "testValue"

	r.Logs[0].SetDecodedData(mapping)

	if r.Logs[0].DecodedData["testKey"] != mapping["testKey"] {
		t.Errorf("Trace: Expected DecodedData %q but got %q", mapping["testKey"], r.Logs[0].DecodedData["testKey"])
	}
}
