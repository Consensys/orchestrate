package types

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// Log holds Data about a log
type Log struct {
	types.Log
	DecodedData map[string]string
}

// Receipt holds Data about a receipt
type Receipt struct {
	types.Receipt
	Logs []*Log

	BlockHash   common.Hash
	BlockNumber uint64
	TxIndex     uint64
}

// NewReceipt creates a new receipt
func newReceipt(root []byte, failed bool, cumulativeGasUsed uint64) *Receipt {
	return &Receipt{
		Receipt: *types.NewReceipt(root, failed, cumulativeGasUsed),
		Logs:    make([]*Log, 0),
	}
}

func (r *Receipt) reset() {
	r.PostState = r.PostState[0:0]
	r.Status = 0
	r.CumulativeGasUsed = 0
	r.Bloom.SetBytes([]byte{})
	r.Logs = r.Logs[0:0]
	r.TxHash.SetBytes([]byte{})
	r.ContractAddress.SetBytes([]byte{})
	r.GasUsed = 0
	r.BlockHash.SetBytes([]byte{})
	r.BlockNumber = 0
	r.TxIndex = 0
}

// SetDecodedData set DecodedData to log
func (l *Log) SetDecodedData(m map[string]string) {
	l.DecodedData = m
}

// String receipt
func (r *Receipt) String() map[string]interface{} {
	receipt := make(map[string]interface{})

	if !bytes.Equal(r.PostState, []byte{}) {
		receipt["PostState"] = hexutil.Encode(r.PostState[:])
	}
	if !reflect.DeepEqual(r.Status, uint64(0)) {
		receipt["Status"] = fmt.Sprintf("%v", r.Status)
	}
	if !reflect.DeepEqual(r.CumulativeGasUsed, uint64(0)) {
		receipt["CumulativeGasUsed"] = fmt.Sprintf("%v", r.CumulativeGasUsed)
	}
	if !reflect.DeepEqual(r.Bloom, types.Bloom{}) {
		receipt["Bloom"] = hexutil.Encode(r.Bloom.Bytes()[:])
	}
	if len(r.Logs) > 0 {
		receipt["Logs"] = r.LogsString()
	}
	if !reflect.DeepEqual(r.TxHash, common.Hash{}) {
		receipt["TxHash"] = hexutil.Encode(r.TxHash[:])
	}
	if !reflect.DeepEqual(r.ContractAddress, common.Address{}) {
		receipt["ContractAddress"] = r.ContractAddress.Hex()
	}
	if !reflect.DeepEqual(r.GasUsed, uint64(0)) {
		receipt["GasUsed"] = fmt.Sprintf("%v", r.GasUsed)
	}
	if !reflect.DeepEqual(r.BlockHash, common.Hash{}) {
		receipt["BlockHash"] = r.BlockHash.Hex()
	}
	if !reflect.DeepEqual(r.BlockNumber, uint64(0)) {
		receipt["BlockNumber"] = fmt.Sprintf("%v", r.BlockNumber)
	}
	if !reflect.DeepEqual(r.TxIndex, uint64(0)) {
		receipt["TxIndex"] = fmt.Sprintf("%v", r.TxIndex)
	}

	return receipt
}

// String log
func (l *Log) String() map[string]interface{} {
	log := map[string]interface{}{
		"Address":     l.Address.Hex(),
		"Topics":      l.TopicsString(),
		"Data":        hexutil.Encode(l.Data[:]),
		"BlockNumber": fmt.Sprintf("%v", l.BlockNumber),
		"TxHash":      l.TxHash.Hex(),
		"TxIndex":     fmt.Sprintf("%v", l.TxIndex),
		"BlockHash":   l.BlockHash.Hex(),
		"Index":       fmt.Sprintf("%v", l.Index),
		"Removed":     l.Removed,
	}
	if len(l.DecodedData) > 0 {
		log["DecodedData"] = l.DecodedData
	}

	return log
}

// TopicsString log
func (l *Log) TopicsString() []string {
	var topics []string
	for _, v := range l.Topics {
		topics = append(topics, v.Hex())
	}
	return topics
}

// LogsString log
func (r *Receipt) LogsString() []map[string]interface{} {
	logs := make([]map[string]interface{}, 0)

	for _, log := range r.Logs {
		logs = append(logs, log.String())
	}

	return logs
}
