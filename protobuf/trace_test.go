package protobuf

import (
	"math/big"
	"testing"

	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

type AccountTest struct {
	id      string
	address string
}

func testPbAccountEquality(pb *tracepb.Account, test *AccountTest, t *testing.T) {
	if pb.GetId() != test.id {
		t.Errorf("Expected UserId to be %q but got %q", test.id, pb.GetId())
	}

	if pb.GetAddress() != test.address {
		t.Errorf("Expected Address to be %q but got %q", test.address, pb.GetAddress())
	}
}

func TestLoadAccount(t *testing.T) {
	acc := types.Account{}
	LoadAccount(nil, &acc)

	pb := tracepb.Account{}
	DumpAccount(&acc, &pb)
	test := AccountTest{"", EmptyAddress}
	testPbAccountEquality(&pb, &test, t)

	LoadAccount(&tracepb.Account{Id: "abc", Address: "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"}, &acc)
	DumpAccount(&acc, &pb)
	test = AccountTest{"abc", "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C"}
	testPbAccountEquality(&pb, &test, t)
}

type ChainTest struct {
	ID       string
	isEIP155 bool
}

func testPbChainEquality(pb *tracepb.Chain, chainTest *ChainTest, t *testing.T) {
	if pb.GetId() != chainTest.ID {
		t.Errorf("Expected ID to be %q but got %q", chainTest.ID, pb.GetId())
	}

	if pb.GetIsEIP155() != chainTest.isEIP155 {
		t.Errorf("Expected IsEIP155 to be %v but got %v", chainTest.isEIP155, pb.GetIsEIP155())
	}
}

func TestLoadDumpChain(t *testing.T) {
	chain := types.Chain{ID: big.NewInt(0)}
	LoadChain(nil, &chain)

	pb := tracepb.Chain{}
	DumpChain(&chain, &pb)
	chainTest := ChainTest{"0x0", false}
	testPbChainEquality(&pb, &chainTest, t)

	LoadChain(&tracepb.Chain{Id: "0xabc", IsEIP155: true}, &chain)
	DumpChain(&chain, &pb)
	chainTest = ChainTest{"0xabc", true}
	testPbChainEquality(&pb, &chainTest, t)
}

type CallTest struct {
	methodID string
	args     []string
}

func testPbCallEquality(pb *tracepb.Call, callTest *CallTest, t *testing.T) {
	if pb.GetMethodId() != callTest.methodID {
		t.Errorf("Expected MethodID to be %q but got %q", callTest.methodID, pb.GetMethodId())
	}

	for i, arg := range pb.GetArgs() {
		if arg != callTest.args[i] {
			t.Errorf("Expected Arg to be %v but got %v", callTest.args[i], arg)
		}
	}
}

func TestLoadDumpCall(t *testing.T) {
	var call types.Call
	LoadCall(nil, &call)

	var pb tracepb.Call
	DumpCall(&call, &pb)
	callTest := CallTest{"", []string{}}
	testPbCallEquality(&pb, &callTest, t)

	LoadCall(&tracepb.Call{MethodId: "abc", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}}, &call)
	DumpCall(&call, &pb)
	callTest = CallTest{"abc", []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}}
	testPbCallEquality(&pb, &callTest, t)
}

type ErrorTest struct {
	typ uint64
	msg string
}

func testPbErrorEquality(pb *tracepb.Error, test *ErrorTest, t *testing.T) {
	if pb.GetType() != test.typ {
		t.Errorf("Expected Type to be %v but got %v", test.typ, pb.GetType())
	}

	if pb.GetMessage() != test.msg {
		t.Errorf("Expected Message to be %q but got %q", test.msg, pb.GetMessage())
	}
}

type TraceTest struct {
	sender   AccountTest
	chain    ChainTest
	receiver AccountTest
	call     CallTest
	tx       TxTest
	receipt  ReceiptTest
	errors   []ErrorTest
}

func testPbTraceEquality(pb *tracepb.Trace, traceTest *TraceTest, t *testing.T) {
	testPbAccountEquality(pb.GetSender(), &traceTest.sender, t)
	testPbChainEquality(pb.GetChain(), &traceTest.chain, t)
	testPbAccountEquality(pb.GetReceiver(), &traceTest.receiver, t)
	testPbCallEquality(pb.GetCall(), &traceTest.call, t)
	testPbTxEquality(pb.GetTransaction(), &traceTest.tx, t)
	testPbReceiptEquality(pb.GetReceipt(), &traceTest.receipt, t)

	if len(pb.GetErrors()) != len(traceTest.errors) {
		t.Errorf("Expected %v errors but got %v", len(traceTest.errors), len(pb.GetErrors()))
	}

	for i, err := range pb.GetErrors() {
		testPbErrorEquality(err, &traceTest.errors[i], t)
	}
}

func TestLoadDumpTrace(t *testing.T) {
	trace := types.NewTrace()
	LoadTrace(nil, trace)

	pb := tracepb.Trace{}
	DumpTrace(trace, &pb)
	traceTest := TraceTest{
		AccountTest{"", EmptyAddress},
		ChainTest{"0x0", false},
		AccountTest{"", EmptyAddress},
		CallTest{"", []string{}},
		TxTest{
			0, EmptyAddress, EmptyQuantity, 0, EmptyQuantity, EmptyData,
			EmptyData,
			EmptyHash,
		},
		nilReceiptTest,
		[]ErrorTest{},
	}
	testPbTraceEquality(&pb, &traceTest, t)

	LoadTrace(
		&tracepb.Trace{
			Sender:   &tracepb.Account{Id: "abc", Address: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
			Chain:    &tracepb.Chain{Id: "0x1afc", IsEIP155: true},
			Receiver: &tracepb.Account{Id: "abc", Address: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
			Call:     &tracepb.Call{MethodId: "abc", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}},
			Transaction: &ethpb.Transaction{
				TxData: &ethpb.TxData{Nonce: 1, To: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684", Value: "0x2386f26fc10000", Gas: 21136, GasPrice: "0xee6b2800", Data: "0xabcd"},
				Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
				Hash:   "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
			},
			Receipt: &ethpb.Receipt{
				Logs:              []*ethpb.Log{},
				ContractAddress:   "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
				PostState:         "0xabcdef",
				Status:            1,
				TxHash:            "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
				Bloom:             "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
				GasUsed:           13456,
				CumulativeGasUsed: 19304777,
			},
			Errors: []*tracepb.Error{&tracepb.Error{Type: 0, Message: "Error 0"}, &tracepb.Error{Type: 1, Message: "Error 1"}},
		},
		trace,
	)
	DumpTrace(trace, &pb)
	traceTest = TraceTest{
		AccountTest{"abc", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
		ChainTest{"0x1afc", true},
		AccountTest{"abc", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
		CallTest{"abc", []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}},
		TxTest{
			1,
			"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
			"0x2386f26fc10000",
			21136,
			"0xee6b2800",
			"0xabcd",
			"0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			"0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
		},
		ReceiptTest{
			[]LogTest{},
			"0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
			"0xabcdef",
			1,
			"0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
			"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			13456,
			19304777,
		},
		[]ErrorTest{{0, "Error 0"}, {1, "Error 1"}},
	}
	testPbTraceEquality(&pb, &traceTest, t)
}
