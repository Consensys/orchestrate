package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// FromGethReceipt creates a receipt from a Geth Receipt
func FromGethReceipt(receipt *types.Receipt) *Receipt {
	logs := []*Log{}
	for _, log := range receipt.Logs {
		logs = append(logs, FromGethLog(log))
	}

	return &Receipt{
		TxHash:            receipt.TxHash.Hex(),
		ContractAddress:   receipt.ContractAddress.Hex(),
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		Bloom:             receipt.Bloom.Bytes(),
		Logs:              logs,
		GasUsed:           receipt.GasUsed,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
	}
}

// SetBlockNumber set block hash
func (receipt *Receipt) SetBlockNumber(number uint64) *Receipt {
	receipt.BlockNumber = number
	return receipt
}

// SetBlockHash set block hash
func (receipt *Receipt) SetBlockHash(h common.Hash) *Receipt {
	receipt.BlockHash = h.Hex()
	return receipt
}

// SetTxHash set transaction hash
func (receipt *Receipt) SetTxHash(h common.Hash) *Receipt {
	receipt.TxHash = h.Hex()
	return receipt
}

// SetTxIndex set transaction index
func (receipt *Receipt) SetTxIndex(idx uint64) *Receipt {
	receipt.TxIndex = idx
	return receipt
}

// TopicsHash return topics in hash format
func (log *Log) TopicsHash() []common.Hash {
	topics := []common.Hash{}
	for _, topic := range log.GetTopics() {
		topics = append(topics, common.HexToHash(topic))
	}
	return topics
}

// FromGethLog creates a new log from a Geth log
func FromGethLog(log *types.Log) *Log {
	// Format topics
	topics := []string{}
	for _, topic := range log.Topics {
		topics = append(topics, topic.Hex())
	}

	return &Log{
		Address:     log.Address.Hex(),
		Topics:      topics,
		Data:        log.Data,
		DecodedData: make(map[string]string),
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.Hex(),
		TxIndex:     uint64(log.TxIndex),
		BlockHash:   log.BlockHash.Hex(),
		Index:       uint64(log.Index),
		Removed:     log.Removed,
	}
}
