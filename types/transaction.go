package types

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// TxData is an adapter from proto buffer to go-ethereum
type TxData struct {
	Nonce    uint64
	To       *common.Address
	Value    *big.Int
	GasLimit uint64
	GasPrice *big.Int
	Data     []byte
}

// GetNonce return Tx Nonce
func (txData *TxData) GetNonce() uint64 {
	return txData.Nonce
}

// SetNonce set Tx Nonce
func (txData *TxData) SetNonce(n uint64) {
	txData.Nonce = n
}

// GetTo return Tx recipient
func (txData *TxData) GetTo() *common.Address {
	return txData.To
}

// SetTo set Tx recipient
func (txData *TxData) SetTo(a *common.Address) {
	txData.To = a
}

// GetValue return Tx Value
func (txData *TxData) GetValue() *big.Int {
	return txData.Value
}

// SetValue set Tx Value
func (txData *TxData) SetValue(v *big.Int) {
	txData.Value = v
}

// GetGasLimit return Tx Gas Limit
func (txData *TxData) GetGasLimit() uint64 {
	return txData.GasLimit
}

// SetGasLimit set Tx Gas Limit
func (txData *TxData) SetGasLimit(l uint64) {
	txData.GasLimit = l
}

// GetGasPrice return Tx Gas Price
func (txData *TxData) GetGasPrice() *big.Int {
	return txData.GasPrice
}

// SetGasPrice set Tx Gas Price
func (txData *TxData) SetGasPrice(p *big.Int) {
	txData.GasPrice = p
}

// GetData return Tx Data Input
func (txData *TxData) GetData() []byte {
	return txData.Data
}

// SetData set Tx Data Input
func (txData *TxData) SetData(data []byte) {
	txData.Data = data
}

// Transaction is an adapter from go-ethereum to protobuf
type Transaction struct {
	TxData *TxData
	Raw    []byte
	Hash   *common.Hash
	From   *common.Address
}

// GetTxData Get Tx Data
func (tx *Transaction) GetTxData() *TxData {
	return tx.TxData
}

// SetTxData set Tx Data
func (tx *Transaction) SetTxData(txData *TxData) {
	tx.TxData = txData
}

// GetFrom return Tx sender
func (tx *Transaction) GetFrom() *common.Address {
	return tx.From
}

// SetFrom set Tx sender
func (tx *Transaction) SetFrom(a *common.Address) {
	tx.From = a
}

// GetRaw return Raw Tx
func (tx *Transaction) GetRaw() []byte {
	return tx.Raw
}

// SetRaw set Raw Tx
func (tx *Transaction) SetRaw(raw []byte) {
	tx.Raw = raw
}

// GetHash return Tx Hash
func (tx *Transaction) GetHash() *common.Hash {
	return tx.Hash
}

// SetHash set Tx Hash
func (tx *Transaction) SetHash(hash *common.Hash) {
	tx.Hash = hash
}

// Get returns a go-ethereum transaction
func (tx *Transaction) Get() *types.Transaction {
	txData := tx.GetTxData()
	return types.NewTransaction(txData.Nonce, *txData.To, txData.Value, txData.GasLimit, txData.GasPrice, txData.Data)
}

// Sign signs transaction
func (tx *Transaction) Sign(s types.Signer, prv *ecdsa.PrivateKey) error {
	// Retrieve go ethereum tx
	t := tx.Get()

	// Sign Tx
	t, err := types.SignTx(t, s, prv)
	if err != nil {
		return err
	}

	// Set raw transaction
	raw, err := rlp.EncodeToBytes(t)
	if err != nil {
		return err
	}
	tx.SetRaw(raw)

	// Set hash
	h := t.Hash()
	tx.SetHash(&h)

	return nil
}
