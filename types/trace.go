package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Chain stores information about an Ethereum chain
type Chain struct {
	// ID chain unique identifier
	ID *big.Int
	// EIP155 indicates whether chain supports EIP155
	IsEIP155 bool
}

func newChain() *Chain {
	return &Chain{ID: big.NewInt(0)}
}

func (c *Chain) reset() {
	c.ID.SetInt64(0)
	c.IsEIP155 = false
}

// Account stores information about an Ethereum account
type Account struct {
	// ID account unique identifier
	ID string
	// Address of account
	Address *common.Address
}

func newAccount() *Account {
	var a common.Address
	return &Account{Address: &a}
}

func (a *Account) reset() {
	a.ID = ""
	a.Address.SetBytes([]byte{})
}

// Call stores information about transaction call
type Call struct {
	// Method method unique identifier
	MethodID string
	// Args arguments to send in the call
	Args []string
}

func newCall() *Call {
	return &Call{}
}

func (c *Call) reset() {
	c.MethodID = ""
	c.Args = c.Args[0:0]
}

// Trace stores contextual information about a transaction call
type Trace struct {
	// Chain chain to execute TX on
	chain *Chain
	// Sender of the transaction
	sender *Account
	// Receiver of the transaction (usually a contract)
	receiver *Account
	// Call information about TX call
	call *Call
	// Tx Transaction being executed
	tx *Tx

	// Tx receipt
	receipt *types.Receipt

	// Errors
	Errors []*Error
}

// NewTrace creates a new trace
func NewTrace() *Trace {
	return &Trace{
		chain:    newChain(),
		sender:   newAccount(),
		receiver: newAccount(),
		call:     newCall(),
		tx:       NewTx(),
		receipt:  types.NewReceipt([]byte{}, true, 0),
	}
}

// Chain returns trace chain
func (t *Trace) Chain() *Chain {
	return t.chain
}

// Sender returns trace sender
func (t *Trace) Sender() *Account {
	return t.sender
}

// Receiver returns trace receiver
func (t *Trace) Receiver() *Account {
	return t.receiver
}

// Call returns trace call
func (t *Trace) Call() *Call {
	return t.call
}

// Tx returns trace Tx
func (t *Trace) Tx() *Tx {
	return t.tx
}

// Receipt returns Tx receipt
func (t *Trace) Receipt() *types.Receipt {
	return t.receipt
}

func (t *Trace) resetReceipt() {
	t.receipt.PostState = t.receipt.PostState[0:0]
	t.receipt.Status = 0
	t.receipt.CumulativeGasUsed = 0
	t.receipt.Bloom.SetBytes([]byte{})
	t.receipt.Logs = t.receipt.Logs[0:0]
	t.receipt.TxHash.SetBytes([]byte{})
	t.receipt.ContractAddress.SetBytes([]byte{})
	t.receipt.GasUsed = 0
}

// Reset re-initiliaze all values stored in trace
func (t *Trace) Reset() {
	t.chain.reset()
	t.sender.reset()
	t.receiver.reset()
	t.call.reset()
	t.tx.reset()
	t.Errors = t.Errors[0:0]

	// Reset receipt
	t.resetReceipt()
}
