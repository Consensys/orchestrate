package protobuf

import (
	"testing"

	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

type SenderTest struct {
	userID       string
	privateKeyID string
}

func testPbSenderEquality(pb *tracepb.Sender, senderTest *SenderTest, t *testing.T) {
	if pb.GetUserId() != senderTest.userID {
		t.Errorf("Expected UserId to be %q but got %q", senderTest.userID, pb.GetUserId())
	}

	if pb.GetPrivateKeyId() != senderTest.privateKeyID {
		t.Errorf("Expected PrivateKeyId to be %q but got %q", senderTest.privateKeyID, pb.GetPrivateKeyId())
	}
}

func TestLoadDumpSender(t *testing.T) {
	var sender types.Sender
	LoadSender(nil, &sender)

	var pb tracepb.Sender
	DumpSender(&sender, &pb)
	senderTest := SenderTest{"", ""}
	testPbSenderEquality(&pb, &senderTest, t)

	LoadSender(&tracepb.Sender{UserId: "abc", PrivateKeyId: "def"}, &sender)
	DumpSender(&sender, &pb)
	senderTest = SenderTest{"abc", "def"}
	testPbSenderEquality(&pb, &senderTest, t)
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
	var chain types.Chain
	LoadChain(nil, &chain)

	var pb tracepb.Chain
	DumpChain(&chain, &pb)
	chainTest := ChainTest{"", false}
	testPbChainEquality(&pb, &chainTest, t)

	LoadChain(&tracepb.Chain{Id: "abc", IsEIP155: true}, &chain)
	DumpChain(&chain, &pb)
	chainTest = ChainTest{"abc", true}
	testPbChainEquality(&pb, &chainTest, t)
}

type ReceiverTest struct {
	ID      string
	address string
}

func testPbReceiverEquality(pb *tracepb.Receiver, receiverTest *ReceiverTest, t *testing.T) {
	if pb.GetId() != receiverTest.ID {
		t.Errorf("Expected ID to be %q but got %q", receiverTest.ID, pb.GetId())
	}

	if pb.GetAddress() != receiverTest.address {
		t.Errorf("Expected IsEIP155 to be %v but got %v", receiverTest.address, pb.GetAddress())
	}
}

func TestLoadDumpReceiver(t *testing.T) {
	var receiver types.Receiver
	LoadReceiver(nil, &receiver)

	var pb tracepb.Receiver
	DumpReceiver(&receiver, &pb)
	receiverTest := ReceiverTest{"", EmptyAddress}
	testPbReceiverEquality(&pb, &receiverTest, t)

	LoadReceiver(&tracepb.Receiver{Id: "abc", Address: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"}, &receiver)
	DumpReceiver(&receiver, &pb)
	receiverTest = ReceiverTest{"abc", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"}
	testPbReceiverEquality(&pb, &receiverTest, t)
}

type CallTest struct {
	methodID string
	value    string
	args     []string
}

func testPbCallEquality(pb *tracepb.Call, callTest *CallTest, t *testing.T) {
	if pb.GetMethodId() != callTest.methodID {
		t.Errorf("Expected MethodID to be %q but got %q", callTest.methodID, pb.GetMethodId())
	}

	if pb.GetValue() != callTest.value {
		t.Errorf("Expected Value to be %v but got %v", callTest.value, pb.GetValue())
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
	callTest := CallTest{"", EmptyQuantity, []string{}}
	testPbCallEquality(&pb, &callTest, t)

	LoadCall(&tracepb.Call{MethodId: "abc", Value: "0xabc", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}}, &call)
	DumpCall(&call, &pb)
	callTest = CallTest{"abc", "0xabc", []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}}
	testPbCallEquality(&pb, &callTest, t)
}

type TraceTest struct {
	sender   SenderTest
	chain    ChainTest
	receiver ReceiverTest
	call     CallTest
	tx       TransactionTest
}

func testPbTraceEquality(pb *tracepb.Trace, traceTest *TraceTest, t *testing.T) {
	testPbSenderEquality(pb.GetSender(), &traceTest.sender, t)
	testPbChainEquality(pb.GetChain(), &traceTest.chain, t)
	testPbReceiverEquality(pb.GetReceiver(), &traceTest.receiver, t)
	testPbCallEquality(pb.GetCall(), &traceTest.call, t)
	testPbTransactionEquality(pb.GetTransaction(), &traceTest.tx, t)
}

func TestLoadDumpTrace(t *testing.T) {
	var trace types.Trace
	LoadTrace(nil, &trace)

	var pb tracepb.Trace
	DumpTrace(&trace, &pb)
	traceTest := TraceTest{
		SenderTest{"", ""},
		ChainTest{"", false},
		ReceiverTest{"", EmptyAddress},
		CallTest{"", EmptyQuantity, []string{}},
		TransactionTest{
			TxDataTest{0, EmptyAddress, EmptyQuantity, 0, EmptyQuantity, EmptyData},
			EmptyData,
			EmptyHash,
			EmptyAddress,
		},
	}
	testPbTraceEquality(&pb, &traceTest, t)

	LoadTrace(
		&tracepb.Trace{
			Sender:   &tracepb.Sender{UserId: "abc", PrivateKeyId: "def"},
			Chain:    &tracepb.Chain{Id: "abc", IsEIP155: true},
			Receiver: &tracepb.Receiver{Id: "abc", Address: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
			Call:     &tracepb.Call{MethodId: "abc", Value: "0xabc", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}},
			Transaction: &ethpb.Transaction{
				TxData: &ethpb.TxData{Nonce: 1, To: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684", Value: "0x2386f26fc10000", Gas: 21136, GasPrice: "0xee6b2800", Data: "0xabcd"},
				Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
				Hash:   "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
				From:   "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
			},
		},
		&trace,
	)
	DumpTrace(&trace, &pb)
	traceTest = TraceTest{
		SenderTest{"abc", "def"},
		ChainTest{"abc", true},
		ReceiverTest{"abc", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
		CallTest{"abc", "0xabc", []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}},
		TransactionTest{
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
		},
	}
	testPbTraceEquality(&pb, &traceTest, t)
}
