package protobuf

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

// LoadAccount load an Account protobuffer to a Account object
func LoadAccount(pb *tracepb.Account, acc *types.Account) error {
	acc.ID = pb.GetId()
	a := common.HexToAddress(pb.GetAddress())
	acc.Address = &a

	return nil
}

// DumpAccount dump an Account object to a protobuffer object
func DumpAccount(acc *types.Account, pb *tracepb.Account) {
	pb.Id = acc.ID
	pb.Address = acc.Address.Hex()
}

// LoadChain load a Chain protobuffer to a Chain object
func LoadChain(pb *tracepb.Chain, chain *types.Chain) error {
	v, err := hexutil.DecodeBig(pb.GetId())
	if err != nil {
		return err
	}
	chain.ID = v
	chain.IsEIP155 = pb.GetIsEIP155()
	return nil
}

// DumpChain dump Chain object to a protobuffer Chain object
func DumpChain(chain *types.Chain, pb *tracepb.Chain) {
	pb.Id = hexutil.EncodeBig(chain.ID)
	pb.IsEIP155 = chain.IsEIP155
}

// LoadCall load a Call protobuffer to a Call object
func LoadCall(pb *tracepb.Call, c *types.Call) {
	c.MethodID = pb.GetMethodId()
	c.Args = pb.GetArgs()
}

// DumpCall dump Call object to a protobuffer Call object
func DumpCall(c *types.Call, pb *tracepb.Call) {
	pb.MethodId = c.MethodID
	pb.Args = c.Args
}

// LoadError load an Error protobuffer to an Error object
func LoadError(pb *tracepb.Error) *types.Error {
	return &types.Error{
		Err:  fmt.Errorf(pb.GetMessage()),
		Type: pb.GetType(),
	}
}

// DumpError dumpt Error object to protobuffer
func DumpError(err *types.Error) *tracepb.Error {
	return &tracepb.Error{
		Message: err.Error(),
		Type:    err.Type,
	}
}

// LoadTrace load a Trace protobuffer to a Trace object
func LoadTrace(pb *tracepb.Trace, t *types.Trace) {
	LoadAccount(pb.GetSender(), t.Sender())
	LoadChain(pb.GetChain(), t.Chain())
	LoadAccount(pb.GetReceiver(), t.Receiver())
	LoadCall(pb.GetCall(), t.Call())
	LoadTx(pb.GetTransaction(), t.Tx())
	LoadReceipt(pb.GetReceipt(), t.Receipt())
	t.Errors = []*types.Error{}
	for _, err := range pb.GetErrors() {
		t.Errors = append(t.Errors, LoadError(err))
	}
}

// DumpTrace dump Trace object to a transaction protobuffer
func DumpTrace(t *types.Trace, pb *tracepb.Trace) {
	if pb.Sender == nil {
		pb.Sender = &tracepb.Account{}
	}
	DumpAccount(t.Sender(), pb.Sender)

	if pb.Chain == nil {
		pb.Chain = &tracepb.Chain{}
	}
	DumpChain(t.Chain(), pb.Chain)

	if pb.Receiver == nil {
		pb.Receiver = &tracepb.Account{}
	}
	DumpAccount(t.Receiver(), pb.Receiver)

	if pb.Call == nil {
		pb.Call = &tracepb.Call{}
	}
	DumpCall(t.Call(), pb.Call)

	if pb.Transaction == nil {
		pb.Transaction = &ethpb.Transaction{}
	}
	DumpTx(t.Tx(), pb.Transaction)

	pb.Errors = []*tracepb.Error{}
	for _, err := range t.Errors {
		pb.Errors = append(pb.Errors, DumpError(err))
	}

	if pb.Receipt == nil {
		pb.Receipt = &ethpb.Receipt{}
	}
	DumpReceipt(t.Receipt(), pb.Receipt)

}
