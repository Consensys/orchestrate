package types

import (
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PrivateArgs are transaction arguments to provide to an Ethereum client supporting privacy (such as Quorum)
type PrivateArgs struct {
	// Private Transaction Fields
	PrivateFrom   string   `json:"privateFrom"`
	PrivateFor    []string `json:"privateFor"`
	PrivateTxType string   `json:"restriction"`
}

// Call2PrivateArgs creates PrivateArgs from a call object
func Call2PrivateArgs(req *tx.Envelope) *PrivateArgs {
	var privateArgs PrivateArgs
	privateArgs.PrivateFrom = req.PrivateFrom
	privateArgs.PrivateFor = req.PrivateFor
	privateArgs.PrivateTxType = req.PrivateTxType
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
func Envelope2SendTxArgs(req *tx.Envelope) (*SendTxArgs, error) {
	from, err := req.GetFromAddress()
	if err != nil {
		return nil, err
	}

	args := SendTxArgs{
		From:        from,
		GasPrice:    (*hexutil.Big)(req.GetGasPrice()),
		Value:       (*hexutil.Big)(req.GetValue()),
		Data:        hexutil.Bytes(req.Data),
		Input:       hexutil.Bytes(req.Data),
		PrivateArgs: *(Call2PrivateArgs(req)),
	}

	if req.Gas != nil {
		args.Gas = (*hexutil.Uint64)(req.Gas)
	}

	if req.Nonce != nil {
		args.Nonce = (*hexutil.Uint64)(req.Gas)
	}

	if req.To != nil {
		args.To = req.GetTo()
	}

	return &args, nil
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
