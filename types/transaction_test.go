package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestTxData(t *testing.T) {
	// Init a TxData object
	txData := TxData{}

	// Test SetNonce
	nonce := uint64(23)
	txData.SetNonce(nonce)
	if txData.GetNonce() != nonce {
		t.Errorf("TxData: expected set Nonce to %q but got %q", nonce, txData.GetNonce())
	}

	// Test SetTo
	to := common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
	txData.SetTo(&to)
	if txData.GetTo().Hex() != to.Hex() {
		t.Errorf("TxData: expected set To to %q but got %q", to.Hex(), txData.GetTo().Hex())
	}

	// Test SetValue
	value := hexutil.MustDecodeBig("0xaf4b80a0d")
	txData.SetValue(value)
	if txData.GetValue() != value {
		t.Errorf("TxData: expected set Value to %q but got %q", value, txData.GetValue())
	}

	// Test SetGasLimit
	gas := uint64(2321564)
	txData.SetGasLimit(gas)
	if txData.GetGasLimit() != gas {
		t.Errorf("TxData: expected set Gas to %q but got %q", gas, txData.GetGasLimit())
	}

	// Test SetGasPrice
	gasPrice := hexutil.MustDecodeBig("0x12ae30fd5c9")
	txData.SetGasPrice(gasPrice)
	if txData.GetGasPrice() != gasPrice {
		t.Errorf("TxData: expected set GasPrice to %q but got %q", gasPrice, txData.GetGasPrice())
	}

	// Test SetData
	data := hexutil.MustDecode("0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd")
	txData.SetData(data)
	if hexutil.Encode(txData.GetData()) != hexutil.Encode(data) {
		t.Errorf("TxData: expected set Data to %q but got %q", hexutil.Encode(data), hexutil.Encode(txData.GetData()))
	}
}

func TestTransaction(t *testing.T) {
	// Init transaction object
	tx := Transaction{}

	var (
		nonce    = uint64(1)
		to       = common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684")
		value    = hexutil.MustDecodeBig("0x2386f26fc10000")
		gas      = uint64(21136)
		gasPrice = hexutil.MustDecodeBig("0xee6b2800")
		data     = hexutil.MustDecode("0xabcd")
	)

	// Test SetTxData
	txData := TxData{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		GasLimit: gas,
		GasPrice: gasPrice,
		Data:     data,
	}
	tx.SetTxData(&txData)
	if tx.GetTxData().GetNonce() != nonce {
		t.Errorf("Transaction: expected set Nonce to %q but got %q", txData.Nonce, nonce)
	}

	var (
		from          = common.HexToAddress("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff")
		pKey          = "86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC" // Corresponds to Address: 0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff
		privateKey, _ = crypto.HexToECDSA(pKey)
		signer        = types.NewEIP155Signer(big.NewInt(3)) // 3 is Ropsten ChainID
	)

	// Test SetFrom
	tx.SetFrom(&from)
	if tx.GetFrom().Hex() != from.Hex() {
		t.Errorf("Transaction: expected set From to %q but got %q", from.Hex(), tx.GetFrom().Hex())
	}

	// Test Sign
	tx.Sign(signer, privateKey)

	raw := hexutil.MustDecode("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")
	if hexutil.Encode(tx.GetRaw()) != hexutil.Encode(raw) {
		t.Errorf("Transaction: expected set Raw to %q but got %q", hexutil.Encode(raw), hexutil.Encode(tx.GetRaw()))
	}

	hash := common.HexToHash("0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd") // Successfully mined on Ropsten: https://ropsten.etherscan.io/tx/0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd
	if tx.GetHash().Hex() != hash.Hex() {
		t.Errorf("Transaction: expected set Hash to %q but got %q", hash.Hex(), tx.GetHash().Hex())
	}
}
