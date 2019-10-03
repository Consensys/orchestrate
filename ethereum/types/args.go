package types

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/envelope"
)

// PrivateArgs are transaction arguments to provide to an Ethereum client supporting privacy (such as Quorum)
type PrivateArgs struct {
	// Private Transaction Fields
	PrivateFrom   string   `json:"privateFrom"`
	PrivateFor    []string `json:"privateFor"`
	PrivateTxType string   `json:"restriction"`
}

// Call2PrivateArgs creates PrivateArgs from a call object
func Call2PrivateArgs(args *envelope.Args) *PrivateArgs {
	var privateArgs PrivateArgs
	privateArgs.PrivateFrom = args.GetPrivate().GetPrivateFrom()
	privateArgs.PrivateFor = args.GetPrivate().GetPrivateFor()
	privateArgs.PrivateTxType = args.GetPrivate().GetPrivateTxType()
	return &privateArgs
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
	from := e.GetFrom().Address()
	args := SendTxArgs{
		From:        from,
		GasPrice:    (*hexutil.Big)(e.GetTx().GetTxData().GetGasPriceBig()),
		Value:       (*hexutil.Big)(e.GetTx().GetTxData().GetValueBig()),
		Data:        hexutil.Bytes(e.GetTx().GetTxData().GetDataBytes()),
		Input:       hexutil.Bytes(e.GetTx().GetTxData().GetDataBytes()),
		PrivateArgs: *(Call2PrivateArgs(e.GetArgs())),
	}

	if gas := e.GetTx().GetTxData().GetGas(); gas != 0 {
		args.Gas = (*hexutil.Uint64)(&gas)
	}

	if nonce := e.GetTx().GetTxData().GetNonce(); nonce != 0 {
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}

	if e.GetTx().GetTxData().GetTo() != nil {
		to := e.GetTx().GetTxData().GetTo().Address()
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
