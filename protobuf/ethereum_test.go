package protobuf

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

const (
	// Empty Ethereum Address
	EmptyAddress = "0x0000000000000000000000000000000000000000"

	// Empty Quantity
	EmptyQuantity = "0x0"

	// Empty Data
	EmptyData = "0x"

	// Empty Hash
	EmptyHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

type Test struct {
	original string
	err      bool
	expected string
}

var addressTests = []Test{
	{"", true, EmptyAddress},
	{"0x", true, EmptyAddress},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", false, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"},
	{"0xaf84242d70ae9d268e2be3616ed497ba28a7b62c", false, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62G", true, EmptyAddress},
}

func TestLoadAddress(t *testing.T) {
	for _, test := range addressTests {
		a, err := LoadAddress(test.original)
		if (err != nil) != test.err {
			if !test.err {
				t.Errorf("LoadAddress: %q is expected to load correctly", test.original)
				if a.Hex() != test.expected {
					t.Errorf("LoadAddress: expected %q but got %q", test.expected, a.Hex())
				}
			} else {
				t.Errorf("LoadAddress: %q is expected to NOT load correctly", test.original)
			}
		}
	}
}

type TxTest struct {
	nonce    uint64
	to       string
	value    string
	gas      uint64
	gasPrice string
	data     string
	raw      string
	hash     string
}

var nilTxTest = TxTest{
	0,
	EmptyAddress,
	EmptyQuantity,
	0,
	EmptyQuantity,
	EmptyData,
	EmptyData,
	EmptyHash,
}

func testTxEquality(tx *types.Tx, test *TxTest, t *testing.T) {
	if tx.Nonce() != test.nonce {
		t.Errorf("Expected Nonce to be %q but got %q", test.nonce, tx.Nonce())
	}

	if tx.To().Hex() != test.to {
		t.Errorf("Expected To to be %q but got %q", test.to, tx.To().Hex())
	}

	if hexutil.EncodeBig(tx.Value()) != test.value {
		t.Errorf("Expected Value to be %q but got %q", test.value, hexutil.EncodeBig(tx.Value()))
	}

	if tx.GasLimit() != test.gas {
		t.Errorf("Expected Gas to be %q but got %q", test.gas, tx.GasLimit())
	}

	if hexutil.EncodeBig(tx.GasPrice()) != test.gasPrice {
		t.Errorf("Expected GasPrice to be %q but got %q", test.gasPrice, hexutil.EncodeBig(tx.GasPrice()))
	}

	if hexutil.Encode(tx.Data()) != test.data {
		t.Errorf("Expected Data to be %q but got %q", test.data, hexutil.Encode(tx.Data()))
	}

	if hexutil.Encode(tx.Raw()) != test.raw {
		t.Errorf("Expected Raw to be %q but got %q", test.raw, hexutil.Encode(tx.Raw()))
	}

	if tx.Hash().Hex() != test.hash {
		t.Errorf("Expected Hash to be %q but got %q", test.hash, tx.Hash().Hex())
	}
}

func testPbTxEquality(pb *ethpb.Transaction, test *TxTest, t *testing.T) {

	if pb.GetTxData().GetNonce() != test.nonce {
		t.Errorf("Expected Nonce to be %q but got %q", test.nonce, pb.GetTxData().GetNonce())
	}

	if pb.GetTxData().GetTo() != test.to {
		t.Errorf("Expected To to be %q but got %q", test.to, pb.GetTxData().GetTo())
	}

	if pb.GetTxData().GetValue() != test.value {
		t.Errorf("Expected Value to be %q but got %q", test.value, pb.GetTxData().GetValue())
	}

	if pb.GetTxData().GetGas() != test.gas {
		t.Errorf("Expected Gas to be %q but got %q", test.gas, pb.GetTxData().GetGas())
	}

	if pb.GetTxData().GetGasPrice() != test.gasPrice {
		t.Errorf("Expected GasPrice to be %q but got %q", test.gasPrice, pb.GetTxData().GetGasPrice())
	}

	if pb.GetTxData().GetData() != test.data {
		t.Errorf("Expected Data to be %q but got %q", test.data, pb.GetTxData().GetData())
	}

	if pb.GetRaw() != test.raw {
		t.Errorf("Expected Raw to be %q but got %q", test.raw, pb.GetRaw())
	}

	if pb.GetHash() != test.hash {
		t.Errorf("Expected Hash to be %q but got %q", test.hash, pb.GetHash())
	}
}

func newTxPb(test TxTest) *ethpb.Transaction {
	return &ethpb.Transaction{
		TxData: &ethpb.TxData{
			Nonce:    test.nonce,
			To:       test.to,
			Value:    test.value,
			Gas:      test.gas,
			GasPrice: test.gasPrice,
			Data:     test.data,
		},
		Raw:  test.raw,
		Hash: test.hash,
	}
}

func TestLoadDumpTx(t *testing.T) {
	tx := types.NewTx()
	// Load nil
	LoadTx(nil, tx)
	testTxEquality(tx, &nilTxTest, t)
	pb := ethpb.Transaction{}

	DumpTx(tx, &pb)
	testPbTxEquality(&pb, &nilTxTest, t)

	// Load not nil
	test := TxTest{
		nonce:    11,
		to:       "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
		value:    "0x400",
		gas:      21136,
		gasPrice: "0xaf84242d70ae9d",
		data:     "0x3616ed497ba28a7b62ca",
		hash:     "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		raw:      "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
	}
	err := LoadTx(newTxPb(test), tx)
	if err != nil {
		t.Errorf("LoadTx: expected successful load but got %v", err)
	}
	pb = ethpb.Transaction{}

	DumpTx(tx, &pb)
	testPbTxEquality(&pb, &test, t)
}
