package ethereum

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

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

// UnmarshalJSON unmarshal from JSON.
func (r *Receipt) UnmarshalJSON(input []byte) error {
	var dec struct {
		PostState         *hexutil.Bytes     `json:"root"`
		Status            *hexutil.Uint64    `json:"status"`
		CumulativeGasUsed *hexutil.Uint64    `json:"cumulativeGasUsed" gencodec:"required"`
		EffectiveGasPrice string             `json:"effectiveGasPrice,omitempty"`
		Bloom             *ethtypes.Bloom    `json:"logsBloom"         gencodec:"required"`
		Logs              []*ethtypes.Log    `json:"logs"              gencodec:"required"`
		TxHash            *ethcommon.Hash    `json:"transactionHash" gencodec:"required"`
		ContractAddress   *ethcommon.Address `json:"contractAddress"`
		GasUsed           *hexutil.Uint64    `json:"gasUsed" gencodec:"required"`
		BlockHash         *ethcommon.Hash    `json:"blockHash,omitempty"`
		BlockNumber       *hexutil.Big       `json:"blockNumber,omitempty"`
		TransactionIndex  *hexutil.Uint      `json:"transactionIndex"`
		RevertReason      string             `json:"revertReason"`
	}
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.PostState != nil {
		r.PostState = hexutil.Encode(*dec.PostState)
	}
	if dec.Status != nil {
		r.Status = uint64(*dec.Status)
	}
	if dec.CumulativeGasUsed == nil {
		return errors.New("missing required field 'cumulativeGasUsed' for Receipt")
	}
	r.CumulativeGasUsed = uint64(*dec.CumulativeGasUsed)

	if dec.EffectiveGasPrice != "" {
		r.EffectiveGasPrice = dec.EffectiveGasPrice
	}

	if dec.Bloom == nil {
		return errors.New("missing required field 'logsBloom' for Receipt")
	}
	r.Bloom = hexutil.Encode(dec.Bloom.Bytes())
	if dec.Logs == nil {
		return errors.New("missing required field 'logs' for Receipt")
	}
	for _, log := range dec.Logs {
		r.Logs = append(r.Logs, FromGethLog(log))
	}
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionHash' for Receipt")
	}
	r.TxHash = dec.TxHash.String()
	if dec.ContractAddress != nil {
		r.ContractAddress = dec.ContractAddress.String()
	}
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Receipt")
	}
	r.GasUsed = uint64(*dec.GasUsed)
	if dec.BlockHash != nil {
		r.BlockHash = dec.BlockHash.String()
	}
	if dec.BlockNumber != nil {
		r.BlockNumber = (*big.Int)(dec.BlockNumber).Uint64()
	}
	if dec.TransactionIndex != nil {
		r.TxIndex = uint64(*dec.TransactionIndex)
	}

	if dec.RevertReason != "" {
		message, err := unpackRevertReason(ethcommon.FromHex(dec.RevertReason))
		if err == nil {
			r.RevertReason = message
		}
	}
	return nil
}

var (
	errorSig     = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _ = abi.NewType("string", "", nil)
)

func unpackRevertReason(result []byte) (string, error) {
	if len(result) < 4 || !bytes.Equal(result[:4], errorSig) {
		return "", errors.New("TX result not of type Error(string)")
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return "", errors.New("unpacking revert reason")
	}
	return vs[0].(string), nil
}
