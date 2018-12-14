package protobuf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// LoadSender load a Sender protobuffer to a Sender object
func LoadSender(pb *tracepb.Sender, sender *types.Sender) {
	if pb == nil {
		pb = &tracepb.Sender{}
	}

	sender.SetUserID(pb.GetUserId())
	sender.SetPrivateKeyID(pb.GetPrivateKeyId())
}

// DumpSender dump Sender object to a protobuffer object
func DumpSender(sender *types.Sender, pb *tracepb.Sender) {
	pb.UserId = sender.GetUserID()
	pb.PrivateKeyId = sender.GetPrivateKeyID()
}

// LoadChain load a Chain protobuffer to a Chain object
func LoadChain(pb *tracepb.Chain, chain *types.Chain) {
	if pb == nil {
		pb = &tracepb.Chain{}
	}

	chain.SetID(pb.GetId())
	chain.SetEIP155(pb.GetIsEIP155())
}

// DumpChain dump Chain object to a protobuffer Chain object
func DumpChain(chain *types.Chain, pb *tracepb.Chain) {
	pb.Id = chain.GetID()
	pb.IsEIP155 = chain.GetEIP155()
}

// LoadReceiver load a Receiver protobuffer to a Receiver object
func LoadReceiver(pb *tracepb.Receiver, r *types.Receiver) {
	if pb == nil {
		pb = &tracepb.Receiver{}
	}

	r.SetID(pb.GetId())

	if r.Address == nil {
		var a common.Address
		r.Address = &a
	}
	LoadAddress(pb.Address, r.Address)
}

// DumpReceiver dump Receiver object to a protobuffer Receiver object
func DumpReceiver(r *types.Receiver, pb *tracepb.Receiver) {
	pb.Id = r.GetID()
	DumpAddress(r.GetAddress(), &pb.Address)
}

// LoadCall load a Call protobuffer to a Call object
func LoadCall(pb *tracepb.Call, c *types.Call) {
	if pb == nil {
		pb = &tracepb.Call{}
	}
	c.SetMethodID(pb.GetMethodId())

	if c.Value == nil {
		var v big.Int
		c.Value = &v
	}
	LoadQuantity(pb.GetValue(), c.Value)

	c.SetArgs(pb.GetArgs())

}

// DumpCall dump Call object to a protobuffer Call object
func DumpCall(c *types.Call, pb *tracepb.Call) {
	pb.MethodId = c.GetMethodID()
	DumpQuantity(c.GetValue(), &pb.Value)
	pb.Args = c.GetArgs()
}

// LoadTrace load a Trace protobuffer to a Trace object
func LoadTrace(pb *tracepb.Trace, t *types.Trace) {
	if pb == nil {
		pb = &tracepb.Trace{}
	}

	if t.Sender == nil {
		var s types.Sender
		t.Sender = &s
	}
	LoadSender(pb.Sender, t.Sender)

	if t.Chain == nil {
		var c types.Chain
		t.Chain = &c
	}
	LoadChain(pb.Chain, t.Chain)

	if t.Receiver == nil {
		var r types.Receiver
		t.Receiver = &r
	}
	LoadReceiver(pb.Receiver, t.Receiver)

	if t.Call == nil {
		var c types.Call
		t.Call = &c
	}
	LoadCall(pb.Call, t.Call)

	if t.Tx == nil {
		var tx types.Transaction
		t.Tx = &tx
	}
	LoadTransaction(pb.Transaction, t.Tx)
}

// DumpTrace dump Trace object to a transaction protobuffer
func DumpTrace(t *types.Trace, pb *tracepb.Trace) {
	if pb.Sender == nil {
		var sender tracepb.Sender
		pb.Sender = &sender
	}
	DumpSender(t.GetSender(), pb.Sender)

	if pb.Chain == nil {
		var chain tracepb.Chain
		pb.Chain = &chain
	}
	DumpChain(t.GetChain(), pb.Chain)

	if pb.Receiver == nil {
		var r tracepb.Receiver
		pb.Receiver = &r
	}
	DumpReceiver(t.GetReceiver(), pb.Receiver)

	if pb.Call == nil {
		var c tracepb.Call
		pb.Call = &c
	}
	DumpCall(t.GetCall(), pb.Call)

	if pb.Transaction == nil {
		var tx ethpb.Transaction
		pb.Transaction = &tx
	}
	DumpTransaction(t.GetTx(), pb.Transaction)
}
