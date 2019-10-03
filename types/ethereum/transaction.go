package ethereum

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// SetNonce set nonce
func (txData *TxData) SetNonce(n uint64) *TxData {
	txData.Nonce = n
	return txData
}

// ToAddress return To in common.Address format
func (txData *TxData) Receiver() ethcommon.Address {
	return txData.GetTo().Address()
}

// SetTo set to address
func (txData *TxData) SetTo(a ethcommon.Address) *TxData {
	if txData.To != nil {
		txData.To.Raw = a.Bytes()
	} else {
		txData.To = &Account{Raw: a.Bytes()}
	}
	return txData
}

// SetValue set value
func (txData *TxData) SetValue(v *big.Int) *TxData {
	if txData.Value != nil {
		txData.Value.Raw = v.Bytes()
	} else {
		txData.Value = &Quantity{Raw: v.Bytes()}
	}
	return txData
}

// GetValueBig returns value of a transaction as a Big integer value
func (txData *TxData) GetValueBig() *big.Int {
	return txData.GetValue().Value()
}

// SetGas set gas limit value
func (txData *TxData) SetGas(l uint64) *TxData {
	txData.Gas = l
	return txData
}

// SetGasPrice set Gas price
func (txData *TxData) SetGasPrice(p *big.Int) *TxData {
	if txData.GasPrice != nil {
		txData.GasPrice.Raw = p.Bytes()
	} else {
		txData.GasPrice = &Quantity{Raw: p.Bytes()}
	}
	return txData
}

// GetGasPriceBig returns gas price in a transaction as a Big integer value
func (txData *TxData) GetGasPriceBig() *big.Int {
	return txData.GetGasPrice().Value()
}

// SetData set Data
func (txData *TxData) SetData(d []byte) *TxData {
	if txData.Data != nil {
		txData.Data.Raw = d
	} else {
		txData.Data = &Data{Raw: d}
	}
	return txData
}

// GetDataBytes set Data
func (txData *TxData) GetDataBytes() []byte {
	return txData.GetData().GetRaw()
}

// SetRaw sets raw transaction
func (tx *Transaction) SetRaw(r []byte) *Transaction {
	if tx.Raw != nil {
		tx.Raw.Raw = r
	} else {
		tx.Raw = &Data{Raw: r}
	}
	return tx
}

// TxHash return transaction hash
func (tx *Transaction) TxHash() ethcommon.Hash {
	if tx.GetHash() == nil {
		return [32]byte{}
	}
	return tx.GetHash().Hash()
}

// SetHash sets transaction hash
func (tx *Transaction) SetHash(h ethcommon.Hash) *Transaction {
	if tx.Hash != nil {
		tx.Hash.Raw = h.Bytes()
	} else {
		tx.Hash = &Hash{Raw: h.Bytes()}
	}

	return tx
}

// IsSigned returns true if transaction is signed, false otherwise
func (tx *Transaction) IsSigned() bool {
	if tx == nil {
		return false
	}
	return tx.GetRaw().GetRaw() != nil && tx.GetHash().GetRaw() != nil
}

// NewTx creates a new transaction
func NewTx() *Transaction {
	return &Transaction{
		TxData: &TxData{},
	}
}
