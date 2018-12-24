package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestTxSet(t *testing.T) {
	// Init a TxData object
	tx := NewTx()

	// Test SetNonce
	nonce := uint64(23)
	tx.SetNonce(nonce)
	if tx.Nonce() != nonce {
		t.Errorf("Tx: expected set Nonce to %q but got %q", nonce, tx.Nonce())
	}

	// Test SetTo
	to := common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
	tx.SetTo(&to)
	if tx.To().Hex() != to.Hex() {
		t.Errorf("Tx: expected set To to %q but got %q", to.Hex(), tx.To().Hex())
	}

	// Test SetValue
	value := hexutil.MustDecodeBig("0xaf4b80a0d")
	tx.SetValue(value)
	if hexutil.EncodeBig(tx.Value()) != hexutil.EncodeBig(value) {
		t.Errorf("Tx: expected set Value to %q but got %q", hexutil.EncodeBig(value), hexutil.EncodeBig(tx.Value()))
	}

	// Test SetGasLimit
	gas := uint64(2321564)
	tx.SetGasLimit(gas)
	if tx.GasLimit() != gas {
		t.Errorf("Tx: expected set Gas to %q but got %q", gas, tx.GasLimit())
	}

	// Test SetGasPrice
	gasPrice := hexutil.MustDecodeBig("0x12ae30fd5c9")
	tx.SetGasPrice(gasPrice)
	if hexutil.EncodeBig(tx.GasPrice()) != hexutil.EncodeBig(gasPrice) {
		t.Errorf("Tx: expected set GasPrice to %q but got %q", hexutil.EncodeBig(gasPrice), hexutil.EncodeBig(tx.GasPrice()))
	}

	// Test SetData
	data := hexutil.MustDecode("0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd")
	tx.SetData(data)
	if hexutil.Encode(tx.Data()) != hexutil.Encode(data) {
		t.Errorf("Tx: expected set Data to %q but got %q", hexutil.Encode(data), hexutil.Encode(tx.Data()))
	}

	// Test SetRaw
	raw := hexutil.MustDecode("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
	tx.SetRaw(raw)
	if hexutil.Encode(tx.Raw()) != hexutil.Encode(raw) {
		t.Errorf("Tx: expected set Raw to %q but got %q", hexutil.Encode(raw), hexutil.Encode(tx.Raw()))
	}

	// Test SetHash
	hash := common.HexToHash("0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210")
	tx.SetHash(&hash)
	if tx.Hash().Hex() != hash.Hex() {
		t.Errorf("Tx: expected set Hash to %q but got %q", hash.Hex(), tx.Hash().Hex())
	}
}

func TestTxSign(t *testing.T) {
	// Prepare Tx
	var (
		nonce    = uint64(1)
		to       = common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
		value    = hexutil.MustDecodeBig("0x2386f26fc10000")
		gas      = uint64(21136)
		gasPrice = hexutil.MustDecodeBig("0xee6b2800")
		data     = hexutil.MustDecode("0xabcd")
	)
	tx := Tx{
		txData: &txData{nonce, &to, value, gas, gasPrice, data},
	}

	// Sign Tx
	var (
		pKey          = "86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC" // Corresponds to Address: 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff
		privateKey, _ = crypto.HexToECDSA(pKey)
		signer        = types.NewEIP155Signer(big.NewInt(3)) // 3 is Ropsten ChainID
	)
	tx.Sign(signer, privateKey)

	raw := hexutil.MustDecode("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
	if hexutil.Encode(tx.Raw()) != hexutil.Encode(raw) {
		t.Errorf("Tx: expected set Raw to %q but got %q", hexutil.Encode(raw), hexutil.Encode(tx.Raw()))
	}

	hash := common.HexToHash("0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd") // Successfully mined on Ropsten: https://ropsten.etherscan.io/tx/0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd
	if tx.Hash().Hex() != hash.Hex() {
		t.Errorf("Tx: expected set Hash to %q but got %q", hash.Hex(), tx.Hash().Hex())
	}
}
