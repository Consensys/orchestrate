package types

import (
	"math/big"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/tx"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PrivateArgs are transaction arguments to provide to an Ethereum client supporting privacy (such as Quorum)
type PrivateArgs struct {
	// Private Transaction Fields
	PrivateFrom    string   `json:"privateFrom"`
	PrivateFor     []string `json:"privateFor"`
	PrivacyGroupID string   `json:"privacyGroupId"`
	PrivateTxType  string   `json:"restriction"`
}

func Call2PrivateArgs(req *tx.Envelope) *PrivateArgs {
	var privateArgs PrivateArgs
	privateArgs.PrivateFrom = req.PrivateFrom
	privateArgs.PrivateFor = req.PrivateFor
	privateArgs.PrivacyGroupID = req.PrivacyGroupID
	privateArgs.PrivateTxType = req.PrivateTxType
	return &privateArgs
}

func CallTx2PrivateArgs(rtx *entities.ETHTransaction) *PrivateArgs {
	var privateArgs PrivateArgs
	privateArgs.PrivateFrom = rtx.PrivateFrom
	privateArgs.PrivateFor = rtx.PrivateFor
	privateArgs.PrivacyGroupID = rtx.PrivacyGroupID
	privateArgs.PrivateTxType = "restricted"
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

// CallArgs contains parameters for contract calls.
type CallArgs struct {
	From     ethcommon.Address  // the sender of the 'transaction'
	To       *ethcommon.Address // the destination contract (nil for contract creation)
	Gas      uint64             // if 0, the call executes with near-infinite gas
	GasPrice *big.Int           // wei <-> gas exchange ratio
	Value    *big.Int           // amount of wei sent along with the call
	Data     []byte             // input data, usually an ABI-encoded contract method invocation
}
