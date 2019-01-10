package protobuf

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
)

const (
	// EmptyEthereum Address
	EmptyAddress = "0x0000000000000000000000000000000000000000"

	// EmptyQuantity
	EmptyQuantity = "0x0"

	// EmptyData
	EmptyData = "0x"

	// EmptyHash
	EmptyHash = "0x0000000000000000000000000000000000000000000000000000000000000000"

	// EmptyBloom
	EmptyBloom = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
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

type LogTest struct {
	address     string
	topics      []string
	data        string
	blockNumber uint64
	txHash      string
	txIndex     uint
	blockHash   string
	index       uint
	removed     bool
}

var nilLogTest = LogTest{
	EmptyAddress,
	[]string{},
	EmptyData,
	0,
	EmptyHash,
	0,
	EmptyHash,
	0,
	false,
}

func testLogEquality(log *ethtypes.Log, test *LogTest, t *testing.T) {
	if log.Address.Hex() != test.address {
		t.Errorf("Expected Address to be %q but got %q", test.address, log.Address.Hex())
	}

	if len(log.Topics) != len(test.topics) {
		t.Errorf("Expected topics count to be %q but got %q", len(test.topics), len(log.Topics))
	}

	for i, topic := range log.Topics {
		if topic.Hex() != test.topics[i] {
			t.Errorf("Expected topics to be %q but got %q", test.topics[i], topic.Hex())
		}
	}

	if hexutil.Encode(log.Data) != test.data {
		t.Errorf("Expected Data to be %q but got %q", test.data, hexutil.Encode(log.Data))
	}

	if log.BlockNumber != test.blockNumber {
		t.Errorf("Expected BlockNumber to be %q but got %q", test.blockNumber, log.BlockNumber)
	}

	if log.TxHash.Hex() != test.txHash {
		t.Errorf("Expected TxHash to be %q but got %q", test.txHash, log.TxHash.Hex())
	}

	if log.TxIndex != test.txIndex {
		t.Errorf("Expected TxIndex to be %q but got %q", test.txIndex, log.TxIndex)
	}

	if log.BlockHash.Hex() != test.blockHash {
		t.Errorf("Expected BlockHash to be %q but got %q", test.blockHash, log.BlockHash.Hex())
	}

	if log.Index != test.index {
		t.Errorf("Expected Index to be %q but got %q", test.index, log.Index)
	}

	if log.Removed != test.removed {
		t.Errorf("Expected Removed to be %v but got %v", test.removed, log.Removed)
	}
}

func testPbLogEquality(pb *ethpb.Log, test *LogTest, t *testing.T) {

	if pb.GetAddress() != test.address {
		t.Errorf("Expected Address to be %q but got %q", test.address, pb.GetAddress())
	}

	if len(pb.GetTopics()) != len(test.topics) {
		t.Errorf("Expected topics count to be %q but got %q", len(test.topics), len(pb.GetTopics()))

		for i, topic := range pb.GetTopics() {
			if topic != test.topics[i] {
				t.Errorf("Expected topics to be %q but got %q", test.topics[i], topic)
			}
		}
	}

	if pb.GetData() != test.data {
		t.Errorf("Expected Data to be %q but got %q", test.data, pb.GetData())
	}

	if pb.GetBlockNumber() != test.blockNumber {
		t.Errorf("Expected BlockNumber to be %q but got %q", test.blockNumber, pb.GetBlockNumber())
	}

	if pb.GetTxHash() != test.txHash {
		t.Errorf("Expected TxHash to be %q but got %q", test.txHash, pb.GetTxHash())
	}

	if pb.GetTxIndex() != uint64(test.txIndex) {
		t.Errorf("Expected TxIndex to be %q but got %q", test.txIndex, pb.GetTxIndex())
	}

	if pb.GetBlockHash() != test.blockHash {
		t.Errorf("Expected BlockHash to be %q but got %q", test.blockHash, pb.GetBlockHash())
	}

	if pb.GetIndex() != uint64(test.index) {
		t.Errorf("Expected Index to be %q but got %q", test.index, pb.GetIndex())
	}

	if pb.GetRemoved() != test.removed {
		t.Errorf("Expected Removed to be %v but got %v", test.removed, pb.GetRemoved())
	}
}

func newLogPb(test LogTest) *ethpb.Log {
	return &ethpb.Log{
		Address:     test.address,
		Topics:      test.topics,
		Data:        test.data,
		BlockNumber: test.blockNumber,
		TxHash:      test.txHash,
		TxIndex:     uint64(test.txIndex),
		BlockHash:   test.blockHash,
		Index:       uint64(test.index),
		Removed:     test.removed,
	}
}

func TestLoadDumpLog(t *testing.T) {
	l := &ethtypes.Log{}

	// Load nil
	LoadLog(nil, l)
	testLogEquality(l, &nilLogTest, t)

	pb := ethpb.Log{}
	DumpLog(l, &pb)
	testPbLogEquality(&pb, &nilLogTest, t)

	// Load not nil
	test := LogTest{
		address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
		topics: []string{
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615",
		},
		data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
		blockNumber: 2019236,
		txHash:      "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
		txIndex:     3,
		blockHash:   "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056",
		removed:     true,
	}
	err := LoadLog(newLogPb(test), l)
	if err != nil {
		t.Errorf("LoadLog: expected successful load but got %v", err)
	}
	pb = ethpb.Log{}

	DumpLog(l, &pb)
	testPbLogEquality(&pb, &test, t)
}

type ReceiptTest struct {
	logs              []LogTest
	contractAddress   string
	postState         string
	status            uint64
	txHash            string
	bloom             string
	gasUsed           uint64
	cumulativeGasUsed uint64
}

var nilReceiptTest = ReceiptTest{
	[]LogTest{},
	EmptyAddress,
	EmptyData,
	0,
	EmptyHash,
	EmptyBloom,
	0,
	0,
}

func testReceiptEquality(r *ethtypes.Receipt, test *ReceiptTest, t *testing.T) {
	if len(r.Logs) != len(test.logs) {
		t.Errorf("Expected logs count to be %q but got %q", len(test.logs), len(r.Logs))
	}

	for i, log := range r.Logs {
		testLogEquality(log, &test.logs[i], t)
	}

	if r.ContractAddress.Hex() != test.contractAddress {
		t.Errorf("Expected ContractAddress to be %q but got %q", test.contractAddress, r.ContractAddress.Hex())
	}

	if hexutil.Encode(r.PostState) != test.postState {
		t.Errorf("Expected postState to be %q but got %q", test.postState, hexutil.Encode(r.PostState))
	}

	if r.Status != test.status {
		t.Errorf("Expected Status to be %q but got %q", test.status, r.Status)
	}

	if r.TxHash.Hex() != test.txHash {
		t.Errorf("Expected TxHash to be %q but got %q", test.txHash, r.TxHash.Hex())
	}

	if common.ToHex(r.Bloom.Bytes()) != test.bloom {
		t.Errorf("Expected Bloom to be %q but got %q", test.bloom, common.ToHex(r.Bloom.Bytes()))
	}

	if r.GasUsed != test.gasUsed {
		t.Errorf("Expected GasUsed to be %v but got %v", test.gasUsed, r.GasUsed)
	}

	if r.CumulativeGasUsed != test.cumulativeGasUsed {
		t.Errorf("Expected CumulativeGasUsed to be %v but got %v", test.cumulativeGasUsed, r.CumulativeGasUsed)
	}
}

func testPbReceiptEquality(pb *ethpb.Receipt, test *ReceiptTest, t *testing.T) {
	if len(pb.GetLogs()) != len(test.logs) {
		t.Errorf("Expected logs count to be %q but got %q", len(test.logs), pb.GetLogs())
	}

	for i, log := range pb.GetLogs() {
		testPbLogEquality(log, &test.logs[i], t)
	}

	if pb.GetContractAddress() != test.contractAddress {
		t.Errorf("Expected ContractAddress to be %q but got %q", test.contractAddress, pb.GetContractAddress())
	}

	if pb.GetPostState() != test.postState {
		t.Errorf("Expected postState to be %q but got %q", test.postState, pb.GetPostState())
	}

	if pb.GetStatus() != test.status {
		t.Errorf("Expected Status to be %q but got %q", test.status, pb.GetStatus())
	}

	if pb.GetTxHash() != test.txHash {
		t.Errorf("Expected TxHash to be %q but got %q", test.txHash, pb.GetTxHash())
	}

	if pb.GetBloom() != test.bloom {
		t.Errorf("Expected Bloom to be %q but got %q", test.bloom, pb.GetBloom())
	}

	if pb.GetGasUsed() != test.gasUsed {
		t.Errorf("Expected GasUsed to be %v but got %v", test.gasUsed, pb.GetGasUsed())
	}

	if pb.GetCumulativeGasUsed() != test.cumulativeGasUsed {
		t.Errorf("Expected CumulativeGasUsed to be %v but got %v", test.cumulativeGasUsed, pb.GetCumulativeGasUsed())
	}
}

func newPbReceipt(test ReceiptTest) *ethpb.Receipt {
	r := &ethpb.Receipt{
		Logs:              []*ethpb.Log{},
		ContractAddress:   test.contractAddress,
		PostState:         test.postState,
		Status:            test.status,
		TxHash:            test.txHash,
		Bloom:             test.bloom,
		GasUsed:           test.gasUsed,
		CumulativeGasUsed: test.cumulativeGasUsed,
	}
	for _, log := range test.logs {
		r.Logs = append(r.Logs, newLogPb(log))
	}
	return r
}

func TestLoadDumpReceipt(t *testing.T) {
	r := &ethtypes.Receipt{}

	// Load nil
	LoadReceipt(nil, r)
	testReceiptEquality(r, &nilReceiptTest, t)

	pb := ethpb.Receipt{}
	DumpReceipt(r, &pb)
	testPbReceiptEquality(&pb, &nilReceiptTest, t)

	test := ReceiptTest{
		logs: []LogTest{{
			address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			topics: []string{
				"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
				"0x00000000000000000000000080b2c9d7cbbf30a1b0fc8983c647d754c6525615",
			},
			data:        "0x000000000000000000000000000000000000000000000001a055690d9db80000",
			blockNumber: 2019236,
			txHash:      "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
			txIndex:     3,
			blockHash:   "0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056",
			removed:     true,
		}},
		contractAddress:   "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
		postState:         "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
		status:            1,
		txHash:            "0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e",
		bloom:             "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a055690d9db80000",
		gasUsed:           156,
		cumulativeGasUsed: 14567,
	}

	err := LoadReceipt(newPbReceipt(test), r)
	if err != nil {
		t.Errorf("LoadLog: expected successful load but got %v", err)
	}
	pb = ethpb.Receipt{}

	DumpReceipt(r, &pb)
	testPbReceiptEquality(&pb, &test, t)
}
