package protobuf

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
	error    bool
	final    string
}

var addressTests = []Test{
	{"", true, EmptyAddress},
	{"0x", true, EmptyAddress},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", false, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"},
	{"0xaf84242d70ae9d268e2be3616ed497ba28a7b62c", false, "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62G", true, EmptyAddress},
}

func testLoadDumpAddress(test Test, t *testing.T) {
	var a common.Address
	err := LoadAddress(test.original, &a)
	if (err != nil) != test.error {
		if !test.error {
			t.Errorf("Address %q is expected to load correctly", test.original)
		} else {
			t.Errorf("Address %q is expected to NOT load correctly", test.original)
		}
	}

	var hex string
	DumpAddress(&a, &hex)
	if hex != test.final {
		t.Errorf("Address expected to dump to %q but go %q", test.final, hex)
	}
}

func TestLoadDumpAddress(t *testing.T) {
	for _, test := range addressTests {
		testLoadDumpAddress(test, t)
	}
}

var quantityTests = []Test{
	{"", true, EmptyQuantity},
	{"0x", true, EmptyQuantity},
	{EmptyQuantity, false, EmptyQuantity},
	{"0x400", false, "0x400"},
	{"0x00400", true, EmptyQuantity},
	{"0xaf84242d70ae9d268e2be3616ed497ba28a7b62c", false, "0xaf84242d70ae9d268e2be3616ed497ba28a7b62c"},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", false, "0xaf84242d70ae9d268e2be3616ed497ba28a7b62c"},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62G", true, EmptyQuantity},
}

func testLoadDumpQuantity(test Test, t *testing.T) {
	q := big.NewInt(0)
	err := LoadQuantity(test.original, q)
	if (err != nil) != test.error {
		if !test.error {
			t.Errorf("Quantity %q is expected to load correctly", test.original)
		} else {
			t.Errorf("Quantity %q is expected to NOT load correctly", test.original)
		}
	}

	var hex string
	DumpQuantity(q, &hex)
	if hex != test.final {
		t.Errorf("Quantity expected to dump to %q but go %q", test.final, hex)
	}
}

func TestLoadDumpQuantity(t *testing.T) {
	for _, test := range quantityTests {
		testLoadDumpQuantity(test, t)
	}
}

var dataTests = []Test{
	{"", true, EmptyData},
	{EmptyData, false, EmptyData},
	{"0xa", true, EmptyData},
	{"0xaf84242d70ae9d268e2be3616ed497ba28a7b62c", false, "0xaf84242d70ae9d268e2be3616ed497ba28a7b62c"},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C", false, "0xaf84242d70ae9d268e2be3616ed497ba28a7b62c"},
	{"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62G", true, EmptyData},
}

func testLoadDumpData(test Test, t *testing.T) {
	var data []byte
	err := LoadData(test.original, &data)
	if (err != nil) != test.error {
		if !test.error {
			t.Errorf("Data %q is expected to load correctly", test.original)
		} else {
			t.Errorf("Data %q is expected to NOT load correctly", test.original)
		}
	}

	var hex string
	DumpData(data, &hex)
	if hex != test.final {
		t.Errorf("Data expected to dump to %q but go %q", test.final, hex)
	}
}

func TestLoadDumpData(t *testing.T) {
	for _, test := range dataTests {
		testLoadDumpData(test, t)
	}
}

var hashTests = []Test{
	{"", false, EmptyHash},
	{"0x", false, EmptyHash},
	{"0xa", false, "0x000000000000000000000000000000000000000000000000000000000000000a"},
	{"0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210", false, "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210"},
}

func testLoadDumpHash(test Test, t *testing.T) {
	var h common.Hash
	err := LoadHash(test.original, &h)
	if (err != nil) != test.error {
		if !test.error {
			t.Errorf("Hash %q is expected to load correctly", test.original)
		} else {
			t.Errorf("Hash %q is expected to NOT load correctly", test.original)
		}
	}

	var hex string
	DumpHash(h, &hex)
	if hex != test.final {
		t.Errorf("Hash expected to dump to %q but go %q", test.final, hex)
	}
}

func TestLoadDumpHash(t *testing.T) {
	for _, test := range hashTests {
		testLoadDumpHash(test, t)
	}
}

type TxDataTest struct {
	nonce    uint64
	to       string
	value    string
	gas      uint64
	gasPrice string
	data     string
}

func testPbTxDataEquality(pb *ethpb.TxData, txDataTest *TxDataTest, t *testing.T) {
	if pb.GetNonce() != txDataTest.nonce {
		t.Errorf("Expected Nonce to be %q but got %q", txDataTest.nonce, pb.GetNonce())
	}

	if pb.GetTo() != txDataTest.to {
		t.Errorf("Expected To to be %q but got %q", txDataTest.to, pb.GetTo())
	}

	if pb.GetValue() != txDataTest.value {
		t.Errorf("Expected Value to be %q but got %q", txDataTest.value, pb.GetValue())
	}

	if pb.GetGas() != txDataTest.gas {
		t.Errorf("Expected Gas to be %q but got %q", txDataTest.gas, pb.GetGas())
	}

	if pb.GetGasPrice() != txDataTest.gasPrice {
		t.Errorf("Expected GasPrice to be %q but got %q", txDataTest.gasPrice, pb.GetGasPrice())
	}

	if pb.GetData() != txDataTest.data {
		t.Errorf("Expected Data to be %q but got %q", txDataTest.data, pb.GetData())
	}
}

func TestLoadDumpTxData(t *testing.T) {
	var txData types.TxData
	LoadTxData(nil, &txData)

	var pb ethpb.TxData
	DumpTxData(&txData, &pb)
	txDataTest := TxDataTest{0, EmptyAddress, EmptyQuantity, 0, EmptyQuantity, EmptyData}
	testPbTxDataEquality(&pb, &txDataTest, t)
}

type TransactionTest struct {
	txData TxDataTest
	raw    string
	hash   string
	from   string
}

func testPbTransactionEquality(pb *ethpb.Transaction, txTest *TransactionTest, t *testing.T) {
	testPbTxDataEquality(pb.GetTxData(), &txTest.txData, t)

	if pb.GetRaw() != txTest.raw {
		t.Errorf("Expected Raw to be %q but got %q", txTest.raw, pb.GetRaw())
	}

	if pb.GetHash() != txTest.hash {
		t.Errorf("Expected Hash to be %q but got %q", txTest.hash, pb.GetHash())
	}

	if pb.GetFrom() != txTest.from {
		t.Errorf("Expected From to be %q but got %q", txTest.from, pb.GetFrom())
	}
}

func TestLoadDumpTransaction(t *testing.T) {
	var tx types.Transaction

	// Test 1: Loading from probuffer nil
	LoadTransaction(nil, &tx)

	var pb ethpb.Transaction
	DumpTransaction(&tx, &pb)
	txTest := TransactionTest{
		TxDataTest{0, EmptyAddress, EmptyQuantity, 0, EmptyQuantity, EmptyData},
		EmptyData,
		EmptyHash,
		EmptyAddress,
	}
	testPbTransactionEquality(&pb, &txTest, t)

	// Test 2: Loading from not nil
	LoadTransaction(
		&ethpb.Transaction{
			TxData: &ethpb.TxData{Nonce: 1, To: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684", Value: "0x2386f26fc10000", Gas: 21136, GasPrice: "0xee6b2800", Data: "0xabcd"},
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
			From:   "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
		},
		&tx,
	)

	DumpTransaction(&tx, &pb)
	txTest = TransactionTest{
		TxDataTest{
			1,
			"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
			"0x2386f26fc10000",
			21136,
			"0xee6b2800",
			"0xabcd",
		},
		"0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
		"0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
		"0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
	}
	testPbTransactionEquality(&pb, &txTest, t)
}
