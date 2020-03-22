package ethereum

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// FromGethReceipt creates a receipt from a Geth Receipt
func FromGethReceipt(r *ethtypes.Receipt) *Receipt {
	var logs []*Log
	for _, log := range r.Logs {
		logs = append(logs, FromGethLog(log))
	}

	return &Receipt{
		TxHash:            r.TxHash.String(),
		ContractAddress:   r.ContractAddress.String(),
		PostState:         hexutil.Encode(r.PostState),
		Status:            r.Status,
		Bloom:             hexutil.Encode(r.Bloom.Bytes()),
		Logs:              logs,
		GasUsed:           r.GasUsed,
		CumulativeGasUsed: r.CumulativeGasUsed,
	}
}

// SetBlockNumber set block hash
func (r *Receipt) SetBlockNumber(number uint64) *Receipt {
	r.BlockNumber = number
	return r
}

// SetBlockHash set block hash
func (r *Receipt) SetBlockHash(h ethcommon.Hash) *Receipt {
	r.BlockHash = h.String()
	return r
}

// SetTxHash set transaction hash
func (r *Receipt) SetTxHash(h ethcommon.Hash) *Receipt {
	r.TxHash = h.String()
	return r
}

// SetTxIndex set transaction index
func (r *Receipt) SetTxIndex(idx uint64) *Receipt {
	r.TxIndex = idx
	return r
}

// FromGethLog creates a new log from a Geth log
func FromGethLog(log *ethtypes.Log) *Log {
	// Format topics
	var topics []string
	for _, topic := range log.Topics {
		topics = append(topics, topic.String())
	}

	return &Log{
		Address:     log.Address.String(),
		Topics:      topics,
		Data:        hexutil.Encode(log.Data),
		DecodedData: make(map[string]string),
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.String(),
		TxIndex:     uint64(log.TxIndex),
		BlockHash:   log.BlockHash.String(),
		Index:       uint64(log.Index),
		Removed:     log.Removed,
	}
}

func (r *Receipt) GetContractAddr() ethcommon.Address {
	return ethcommon.HexToAddress(r.GetContractAddress())
}

func (r *Receipt) GetTxHashPtr() *ethcommon.Hash {
	hash := ethcommon.HexToHash(r.GetTxHash())
	return &hash
}
