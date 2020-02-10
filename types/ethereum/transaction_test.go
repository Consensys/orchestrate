package ethereum

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

const (
	// EmptyEthereum Address
	EmptyAddress = "0x0000000000000000000000000000000000000000"

	// EmptyHash
	EmptyHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

func TestTxData(t *testing.T) {
	// Test on empty TxData
	var txData *TxData

	assert.Equal(t, EmptyAddress, txData.Receiver().Hex(), "Address should be empty")
	assert.Equal(t, int64(0), txData.GetValueBig().Int64(), "Value should be 0")
	assert.Equal(t, int64(0), txData.GetGasPriceBig().Int64(), "Gas price should be 0")
	assert.Equal(t, "", txData.GetData(), "Data should be empty")

	// // TxData information
	txData = &TxData{}
	txData = txData.
		SetNonce(10).
		SetTo(common.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	assert.Equal(t, uint64(10), txData.GetNonce(), "Nonce should be set")
	assert.Equal(t, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", txData.Receiver().Hex(), "To Address should be set")
	assert.Equal(t, int64(100000), txData.GetValueBig().Int64(), "Value should be set")
	assert.Equal(t, int64(200000), txData.GetGasPriceBig().Int64(), "Gas price should be set")
	assert.Equal(t, uint64(2000), txData.GetGas(), "Gas should be set")
	assert.Equal(t, "0xabcd", txData.GetData(), "Data should be set")
}

func TestTransaction(t *testing.T) {
	var tx *Transaction
	assert.Equal(t, EmptyHash, tx.TxHash().Hex(), "Hash should be empty")

	tx = NewTx().
		SetRaw(hexutil.MustDecode("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")).
		SetHash(common.HexToHash("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"))

	tx.TxData.SetNonce(10)

	assert.Equal(
		t,
		"0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		tx.TxHash().Hex(),
		"Hash should be empty",
	)
	assert.Equal(t,
		"0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
		tx.GetRaw(),
		"Raw should be set",
	)
	assert.Equal(t,
		uint64(10),
		tx.GetTxData().GetNonce(),
		"Nonce should be set",
	)
	assert.Truef(t,
		tx.IsSigned(),
		"Transaction is signed",
	)
}

func TestTransactionWithoutRawDataAndHashAreUnsigned(t *testing.T) {
	var tx *Transaction
	assert.Equal(t, EmptyHash, tx.TxHash().Hex(), "Hash should be empty")

	tx = NewTx()

	assert.Falsef(t,
		tx.IsSigned(),
		"Transaction is unsigned",
	)
}

func TestTransaction_IsSigned(t *testing.T) {
	tx := &Transaction{}
	assert.False(t, tx.IsSigned(), "should not be signed")
}

func TestTransaction_SetHash(t *testing.T) {
	h := common.BigToHash(big.NewInt(1))
	tx := (&Transaction{}).SetHash(h)
	assert.Equal(t, h.Hex(), tx.GetHash(), "should not be equal")
}

func TestTransaction_SetRaw(t *testing.T) {
	r := []byte{1}
	tx := (&Transaction{}).SetRaw(r)
	assert.Equal(t, hexutil.Encode(r), tx.GetRaw(), "should not be equal")
}

func TestTxData_SetTo(t *testing.T) {
	to := common.HexToAddress("0x0")
	txData := (&TxData{}).SetTo(to)
	assert.Equal(t, to.Hex(), txData.GetTo(), "should not be equal")
}

func TestTxData_SetValue(t *testing.T) {
	v := big.NewInt(1)
	txData := (&TxData{}).SetValue(v)
	assert.Equal(t, v, txData.GetValueBig(), "should not be equal")
}
