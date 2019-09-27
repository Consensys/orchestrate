package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// FromGethReceipt creates a receipt from a Geth Receipt
func FromGethReceipt(r *ethtypes.Receipt) *Receipt {
	logs := []*Log{}
	for _, log := range r.Logs {
		logs = append(logs, FromGethLog(log))
	}

	return &Receipt{
		TxHash:            &Hash{Raw: r.TxHash.Bytes()},
		ContractAddress:   &Account{Raw: r.ContractAddress.Bytes()},
		PostState:         r.PostState,
		Status:            r.Status,
		Bloom:             r.Bloom.Bytes(),
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
func (r *Receipt) SetBlockHash(h common.Hash) *Receipt {
	if r.GetBlockHash() != nil {
		r.GetBlockHash().Raw = h.Bytes()
	} else {
		r.BlockHash = &Hash{Raw: h.Bytes()}
	}

	return r
}

// SetTxHash set transaction hash
func (r *Receipt) SetTxHash(h common.Hash) *Receipt {
	if r.GetTxHash() != nil {
		r.GetTxHash().Raw = h.Bytes()
	} else {
		r.TxHash = &Hash{Raw: h.Bytes()}
	}
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
	topics := []*Hash{}
	for _, topic := range log.Topics {
		topics = append(topics, &Hash{Raw: topic.Bytes()})
	}

	return &Log{
		Address:     &Account{Raw: log.Address.Bytes()},
		Topics:      topics,
		Data:        log.Data,
		DecodedData: make(map[string]string),
		BlockNumber: log.BlockNumber,
		TxHash:      &Hash{Raw: log.TxHash.Bytes()},
		TxIndex:     uint64(log.TxIndex),
		BlockHash:   &Hash{Raw: log.BlockHash.Bytes()},
		Index:       uint64(log.Index),
		Removed:     log.Removed,
	}
}
