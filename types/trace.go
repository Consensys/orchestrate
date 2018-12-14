package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Sender stores information concerning transaction sender
type Sender struct {
	UserID       string
	PrivateKeyID string
}

// GetUserID return sender unique identifier
func (sender *Sender) GetUserID() string {
	return sender.UserID
}

// SetUserID set sender identifier
func (sender *Sender) SetUserID(id string) {
	sender.UserID = id
}

// GetPrivateKeyID return private key unique identifier
func (sender *Sender) GetPrivateKeyID() string {
	return sender.PrivateKeyID
}

// SetPrivateKeyID set private key identifier
func (sender *Sender) SetPrivateKeyID(id string) {
	sender.PrivateKeyID = id
}

// Chain stores information about a chain
type Chain struct {
	ID       string
	IsEIP155 bool
}

// GetID return chain unique identifier
func (chain *Chain) GetID() string {
	return chain.ID
}

// SetID set chain identifier
func (chain *Chain) SetID(id string) {
	chain.ID = id
}

// GetEIP155 indicates whether chain supports EIP155
func (chain *Chain) GetEIP155() bool {
	return chain.IsEIP155
}

// SetEIP155 set chain EIP155 indicator
func (chain *Chain) SetEIP155(b bool) {
	chain.IsEIP155 = b
}

// Receiver stores information about receiver of the transaction
type Receiver struct {
	ID      string
	Address *common.Address
}

// GetID return receiver unique identifier
func (r *Receiver) GetID() string {
	return r.ID
}

// SetID set receiver identifier
func (r *Receiver) SetID(id string) {
	r.ID = id
}

// GetAddress return receiver address
func (r *Receiver) GetAddress() *common.Address {
	return r.Address
}

// SetAddress set receiver address
func (r *Receiver) SetAddress(a *common.Address) {
	r.Address = a
}

// Call stores information about transaction call
type Call struct {
	MethodID string
	Value    *big.Int
	Args     []string
}

// GetMethodID return call unique identifier
func (c *Call) GetMethodID() string {
	return c.MethodID
}

// SetMethodID set call method id
func (c *Call) SetMethodID(id string) {
	c.MethodID = id
}

// GetValue return call value
func (c *Call) GetValue() *big.Int {
	return c.Value
}

// SetValue set call value
func (c *Call) SetValue(v *big.Int) {
	c.Value = v
}

// GetArgs return call args
func (c *Call) GetArgs() []string {
	return c.Args
}

// SetArgs set call args
func (c *Call) SetArgs(args []string) {
	c.Args = args
}

// Trace stores contextual information about a transaction call
type Trace struct {
	Sender   *Sender
	Chain    *Chain
	Receiver *Receiver
	Call     *Call
	Tx       *Transaction
}

// GetSender return Trace sender
func (t *Trace) GetSender() *Sender {
	return t.Sender
}

// SetSender set trace sender
func (t *Trace) SetSender(s *Sender) {
	t.Sender = s
}

// GetChain return Trace chain
func (t *Trace) GetChain() *Chain {
	return t.Chain
}

// SetChain set trace chain
func (t *Trace) SetChain(c *Chain) {
	t.Chain = c
}

// GetReceiver return Trace receiver
func (t *Trace) GetReceiver() *Receiver {
	return t.Receiver
}

// SetReceiver set trace receiver
func (t *Trace) SetReceiver(r *Receiver) {
	t.Receiver = r
}

// GetCall return Trace call
func (t *Trace) GetCall() *Call {
	return t.Call
}

// SetCall set trace call
func (t *Trace) SetCall(c *Call) {
	t.Call = c
}

// GetTx return Trace transaction
func (t *Trace) GetTx() *Transaction {
	return t.Tx
}

// SetTx set trace call
func (t *Trace) SetTx(tx *Transaction) {
	t.Tx = tx
}
