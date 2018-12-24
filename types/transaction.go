package types

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type txData struct {
	nonce    uint64
	to       *common.Address
	value    *big.Int
	gasLimit uint64
	gasPrice *big.Int
	data     []byte
}

func newTxData() txData {
	var (
		a common.Address
		v big.Int
		p big.Int
	)
	return txData{to: &a, value: &v, gasPrice: &p}
}

func (txData *txData) reset() {
	txData.nonce = 0
	txData.to.SetBytes([]byte{})
	txData.value.SetUint64(0)
	txData.gasLimit = 0
	txData.gasPrice.SetUint64(0)
	txData.data = txData.data[0:0]
}

// Tx stores information about the transaction
type Tx struct {
	txData *txData
	raw    []byte
	hash   *common.Hash
}

// NewTx creates a new transaction store
func NewTx() Tx {
	txData := newTxData()
	var h common.Hash
	return Tx{txData: &txData, hash: &h}
}

func (tx *Tx) reset() {
	tx.txData.reset()
	tx.raw = tx.raw[0:0]
	tx.hash.SetBytes([]byte{})
}

// Nonce returns Tx nonce
func (tx *Tx) Nonce() uint64 {
	return tx.txData.nonce
}

// SetNonce set Tx nonce
func (tx *Tx) SetNonce(n uint64) {
	tx.txData.nonce = n
}

// To returns Tx recipient
func (tx *Tx) To() *common.Address {
	return tx.txData.to
}

// SetTo set Tx recipient
func (tx *Tx) SetTo(a *common.Address) {
	tx.txData.to = a
}

// Value returns Tx value
func (tx *Tx) Value() *big.Int {
	return tx.txData.value
}

// SetValue set Tx value
func (tx *Tx) SetValue(b *big.Int) {
	tx.txData.value = b
}

// GasLimit returns Tx gas limit
func (tx *Tx) GasLimit() uint64 {
	return tx.txData.gasLimit
}

// SetGasLimit set Tx gas limit
func (tx *Tx) SetGasLimit(l uint64) {
	tx.txData.gasLimit = l
}

// GasPrice returns Tx gas price
func (tx *Tx) GasPrice() *big.Int {
	return tx.txData.gasPrice
}

// SetGasPrice set Tx gas price
func (tx *Tx) SetGasPrice(b *big.Int) {
	tx.txData.gasPrice = b
}

// Data returns Tx data
func (tx *Tx) Data() []byte {
	return tx.txData.data
}

// SetData set Tx data
func (tx *Tx) SetData(b []byte) {
	tx.txData.data = b
}

// Raw returns raw Tx
func (tx *Tx) Raw() []byte {
	return tx.raw
}

// SetRaw set raw Tx
func (tx *Tx) SetRaw(b []byte) {
	tx.raw = b
}

// Hash returns Tx hash
func (tx *Tx) Hash() *common.Hash {
	return tx.hash
}

// SetHash set Tx hash
func (tx *Tx) SetHash(h *common.Hash) {
	tx.hash = h
}

// Sign signs transaction
func (tx *Tx) Sign(s types.Signer, prv *ecdsa.PrivateKey) error {
	// Create Go-Ethereum transaction object
	txData := tx.txData
	t := types.NewTransaction(txData.nonce, *txData.to, txData.value, txData.gasLimit, txData.gasPrice, txData.data)

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
