package types

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
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

// String chain
func (c *Chain) String() map[string]interface{} {
	chain := make(map[string]interface{})

	chain["ID"] = fmt.Sprintf("%v", c.ID)
	if c.IsEIP155 != false {
		chain["IsEIP155"] = c.IsEIP155
	}

	return chain
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

// String account
func (a *Account) String() map[string]interface{} {
	account := make(map[string]interface{})

	if a.ID != "" {
		account["ID"] = a.ID
	}
	if !reflect.DeepEqual(a.Address, &common.Address{}) {
		account["Address"] = a.Address.Hex()
	}

	return account
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

func (c *Call) String() map[string]interface{} {
	call := make(map[string]interface{})

	if c.MethodID != "" {
		call["MethodID"] = c.MethodID
	}
	if len(c.Args) > 0 {
		call["Args"] = c.Args
	}

	return call
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
	receipt *Receipt

	// Errors
	Errors Errors
}

// NewTrace creates a new trace
func NewTrace() *Trace {
	return &Trace{
		chain:    newChain(),
		sender:   newAccount(),
		receiver: newAccount(),
		call:     newCall(),
		tx:       NewTx(),
		receipt:  newReceipt([]byte{}, true, 0),
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
func (t *Trace) Receipt() *Receipt {
	return t.receipt
}

// Reset re-initiliaze all values stored in trace
func (t *Trace) Reset() {
	t.chain.reset()
	t.sender.reset()
	t.receiver.reset()
	t.call.reset()
	t.tx.reset()
	t.receipt.reset()
	t.Errors = t.Errors[0:0]
}

// String returns Trace mapping
func (t *Trace) String() map[string]interface{} {
	trace := make(map[string]interface{})

	if !reflect.DeepEqual(t.Chain(), newChain()) {
		trace["Chain"] = t.Chain().String()
	}
	if !reflect.DeepEqual(t.Sender(), newAccount()) {
		trace["Sender"] = t.Sender().String()
	}
	if !reflect.DeepEqual(t.Receiver(), newAccount()) {
		trace["Receiver"] = t.Receiver().String()
	}
	if !reflect.DeepEqual(t.Call(), newCall()) {
		trace["Call"] = t.Call().String()
	}
	if !reflect.DeepEqual(t.Tx(), NewTx()) {
		trace["Tx"] = t.Tx().String()
	}
	if !reflect.DeepEqual(t.Receipt(), newReceipt([]byte{}, true, 0)) {
		trace["Receipt"] = t.Receipt().String()
	}

	return trace
}
