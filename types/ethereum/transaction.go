package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// SetNonce set nonce
func (txData *TxData) SetNonce(n uint64) *TxData {
	txData.Nonce = n
	return txData
}

// ToAddress return To in common.Address format
func (txData *TxData) Receiver() ethcommon.Address {
	if txData.GetTo() == "" {
		return ethcommon.Address{0}
	}

	return ethcommon.HexToAddress(txData.GetTo())
}

// SetTo set to address
func (txData *TxData) SetTo(to ethcommon.Address) *TxData {
	txData.To = to.String()
	return txData
}

// SetValue set value
func (txData *TxData) SetValue(value *big.Int) *TxData {
	txData.Value = value.String()
	return txData
}

// GetValueBig returns value of a transaction as a Big integer value
func (txData *TxData) GetValueBig() *big.Int {
	if txData.GetValue() == "" {
		return big.NewInt(0)
	}
	b, _ := new(big.Int).SetString(txData.GetValue(), 10)
	return b
}

// SetGas set gas limit value
func (txData *TxData) SetGas(l uint64) *TxData {
	txData.Gas = l
	return txData
}

// SetGasPrice set Gas price
func (txData *TxData) SetGasPrice(gasPrice *big.Int) *TxData {
	txData.GasPrice = gasPrice.String()
	return txData
}

// GetGasPriceBig returns gas price in a transaction as a Big integer value
func (txData *TxData) GetGasPriceBig() *big.Int {
	if txData.GetGasPrice() == "" {
		return big.NewInt(0)
	}
	b, _ := new(big.Int).SetString(txData.GetGasPrice(), 10)
	return b
}

// SetData set Data
func (txData *TxData) SetData(d []byte) *TxData {
	txData.Data = hexutil.Encode(d)
	return txData
}

// GetDataBytes set Data
func (txData *TxData) GetDataBytes() []byte {
	b, _ := hexutil.Decode(txData.GetData())
	return b
}

// SetRaw sets raw transaction
func (tx *Transaction) SetRaw(r []byte) *Transaction {
	tx.Raw = hexutil.Encode(r)
	return tx
}

// TxHash return transaction hash
func (tx *Transaction) TxHash() ethcommon.Hash {
	if tx.GetHash() == "" {
		return [32]byte{}
	}
	return ethcommon.HexToHash(tx.GetHash())
}

// SetHash sets transaction hash
func (tx *Transaction) SetHash(h ethcommon.Hash) *Transaction {
	tx.Hash = h.String()
	return tx
}

// IsSigned returns true if transaction is signed, false otherwise
func (tx *Transaction) IsSigned() bool {
	if tx == nil {
		return false
	}
	return tx.GetRaw() != "" && tx.GetHash() != ""
}

// NewTx creates a new transaction
func NewTx() *Transaction {
	return &Transaction{
		TxData: &TxData{},
	}
}
