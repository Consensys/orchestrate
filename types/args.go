package types

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

// PrivateArgs are transaction arguments to provide to an Ethereum client supporting privacy (such as Quorum)
type PrivateArgs struct {
	// Private Transaction Fields
	PrivateFrom   string   `json:"privateFrom"`
	PrivateFor    []string `json:"privateFor"`
	PrivateTxType string   `json:"restriction"`
}

// Call2PrivateArgs creates PrivateArgs from a call object
func Call2PrivateArgs(call *common.Call) *PrivateArgs {
	var args PrivateArgs
	args.PrivateFrom = call.GetQuorum().GetPrivateFrom()
	args.PrivateFor = call.GetQuorum().GetPrivateFor()
	args.PrivateTxType = call.GetQuorum().GetPrivateTxType()
	return &args
}

// SendTxArgs are arguments to provide to jsonRPC call `eth_sendTransaction`
type SendTxArgs struct {
	// From address in case of a non raw transaction
	From ethcommon.Address `json:"from"`

	// Main transaction attributes
	To       *ethcommon.Address `json:"to"`
	Gas      *hexutil.Uint64    `json:"gas"`
	GasPrice *hexutil.Big       `json:"gasPrice"`
	Value    *hexutil.Big       `json:"value"`
	Nonce    *hexutil.Uint64    `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred
	Data  hexutil.Bytes `json:"data"`
	Input hexutil.Bytes `json:"input"`

	// Private field
	PrivateArgs
}

// Envelope2SendTxArgs creates SendTxArgs from an Envelope
func Envelope2SendTxArgs(e *envelope.Envelope) *SendTxArgs {
	From, _ := e.GetSender().Address()
	args := SendTxArgs{
		From:        From,
		GasPrice:    (*hexutil.Big)(e.GetTx().GetTxData().GasPriceBig()),
		Value:       (*hexutil.Big)(e.GetTx().GetTxData().ValueBig()),
		Data:        hexutil.Bytes(e.GetTx().GetTxData().DataBytes()),
		Input:       hexutil.Bytes(e.GetTx().GetTxData().DataBytes()),
		PrivateArgs: *(Call2PrivateArgs(e.GetCall())),
	}

	if gas := e.GetTx().GetTxData().GetGas(); gas != 0 {
		args.Gas = (*hexutil.Uint64)(&gas)
	}

	if nonce := e.GetTx().GetTxData().GetNonce(); nonce != 0 {
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}

	if e.GetTx().GetTxData().GetTo() != "" {
		to, _ := e.GetTx().GetTxData().ToAddress()
		args.To = &to
	}

	return &args
}

// CallArgs contains parameters for contract calls.
type CallArgs struct {
	From     ethcommon.Address  // the sender of the 'transaction'
	To       *ethcommon.Address // the destination contract (nil for contract creation)
	Gas      uint64             // if 0, the call executes with near-infinite gas
	GasPrice *big.Int           // wei <-> gas exchange ratio
	Value    *big.Int           // amount of wei sent along with the call
	Data     []byte             // input data, usually an ABI-encoded contract method invocation
}
