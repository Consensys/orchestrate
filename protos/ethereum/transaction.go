package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SetNonce set nonce
func (txData *TxData) SetNonce(n uint64) *TxData {
	txData.Nonce = n
	return txData
}

// ToAddress return To in common.Address format
func (txData *TxData) ToAddress() common.Address {
	if txData.GetTo() == "" {
		return common.HexToAddress("")
	}
	if !common.IsHexAddress(txData.GetTo()) {
		panic(fmt.Sprintf("%q is an invalid Ethereum address", txData.GetTo()))
	}
	return common.HexToAddress(txData.GetTo())
}

// SetTo set to address
func (txData *TxData) SetTo(a common.Address) *TxData {
	txData.To = a.Hex()
	return txData
}

// ValueBig return value in big Int format
func (txData *TxData) ValueBig() *big.Int {
	if txData.GetValue() == "" {
		return big.NewInt(0)
	}
	return hexutil.MustDecodeBig(txData.GetValue())
}

// SetValue set value
func (txData *TxData) SetValue(v *big.Int) *TxData {
	txData.Value = hexutil.EncodeBig(v)
	return txData
}

// SetGas set gas limit value
func (txData *TxData) SetGas(l uint64) *TxData {
	txData.Gas = l
	return txData
}

// GasPriceBig return gas price in big.Int format
func (txData *TxData) GasPriceBig() *big.Int {
	if txData.GetGasPrice() == "" {
		return big.NewInt(0)
	}
	return hexutil.MustDecodeBig(txData.GetGasPrice())
}

// SetGasPrice set Gas price
func (txData *TxData) SetGasPrice(p *big.Int) *TxData {
	txData.GasPrice = hexutil.EncodeBig(p)
	return txData
}

// DataBytes return data in byte slice format
func (txData *TxData) DataBytes() []byte {
	if txData.GetData() == "" {
		return []byte{}
	}
	return hexutil.MustDecode(txData.GetData())
}

// SetData set Data
func (txData *TxData) SetData(d []byte) *TxData {
	txData.Data = hexutil.Encode(d)
	return txData
}

// SetRaw sets raw transaction
func (tx *Transaction) SetRaw(r string) *Transaction {
	tx.Raw = r
	return tx
}

// TxHash return transaction hash
func (tx *Transaction) TxHash() common.Hash {
	if tx.GetHash() == "" {
		return common.Hash([32]byte{})
	}
	return common.HexToHash(tx.GetHash())
}

// SetHash sets transaction hash
func (tx *Transaction) SetHash(h common.Hash) *Transaction {
	tx.Hash = h.Hex()
	return tx
}

// NewTx creates a new transaction
func NewTx() *Transaction {
	return &Transaction{
		TxData: &TxData{},
	}
}
