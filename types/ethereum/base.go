package ethereum

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const addressLength = 20
const hashLength = 32

// IsValidAddress returns true if the byte array has a correct length
func IsValidAddress(address []byte) bool {
	return len(address) == addressLength
}

// NewAccount create a new instance of an Account message and validates the address field
func NewAccount(address []byte) *Account {
	return (&Account{}).SetAddress(address)
}

// HexToAccount create a new instance of an Account message and validates the address field
func HexToAccount(addressHex string) *Account {
	return NewAccount(hexutil.MustDecode(addressHex))
}

// SetAddress sets the value of the address field and validates it
func (a *Account) SetAddress(address []byte) *Account {
	if !IsValidAddress(address) {
		panic("address should be 20 bytes long")
	}

	a.Raw = address
	return a
}

// Address return Account address in Geth types
func (a *Account) Address() ethcommon.Address {
	// TODO: we could have a performance overhead here
	// A solution could be to add an attribute keeping loaded the address cached
	return ethcommon.BytesToAddress(a.GetRaw())
}

// Hex returns Account address in Geth types
func (a *Account) Hex() string {
	if a == nil {
		return "0x"
	}
	return a.Address().Hex()
}

// IsValidHash returns true if the byte array has a correct length
func IsValidHash(hash []byte) bool {
	return len(hash) == hashLength
}

// NewHash creates an instance of Hash message
func NewHash(raw []byte) *Hash {
	return (&Hash{}).SetValue(raw)
}

// HexToHash creates an instance of Hash message by parsing a string first
func HexToHash(hashHex string) *Hash {
	return NewHash(hexutil.MustDecode(hashHex))
}

// Hash return Hash in Geth types
func (h *Hash) Hash() ethcommon.Hash {
	return ethcommon.BytesToHash(h.GetRaw())
}

// SetValue sets a raw value of a hash
func (h *Hash) SetValue(hash []byte) *Hash {
	if !IsValidHash(hash) {
		panic("hash should be 32 bytes long")
	}
	h.Raw = hash
	return h
}

// Hex return Hash in Geth types
func (h *Hash) Hex() string {
	if h == nil {
		return "0x"
	}
	return hexutil.Encode(h.GetRaw())
}

// IntToQuantity creates an instance of Quantity message from a number
func IntToQuantity(value int64) *Quantity {
	return NewQuantity(big.NewInt(value).Bytes())
}

// HexToQuantity creates an instance of Quantity message from a string
func HexToQuantity(hex string) *Quantity {
	return NewQuantity(hexutil.MustDecode(hex))
}

// NewQuantity creates an instance of Quantity message from an array of bytes
func NewQuantity(raw []byte) *Quantity {
	return &Quantity{
		Raw: raw,
	}
}

// SetRawValue allows to set a new value
func (q *Quantity) SetRawValue(value []byte) {
	q.Raw = big.NewInt(0).SetBytes(value).Bytes()
}

// SetIntValue allows to set a new value
func (q *Quantity) SetIntValue(value int64) {
	q.Raw = big.NewInt(value).Bytes()
}

// Value return quantity in big.Int format
func (q *Quantity) Value() *big.Int {
	if q == nil {
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(q.GetRaw())
}

// NewData create an instance of Data message
func NewData(raw []byte) *Data {
	return &Data{
		Raw: raw,
	}
}

// HexToData create an instance of Data message
func HexToData(raw string) *Data {
	return &Data{
		Raw: hexutil.MustDecode(raw),
	}
}

// Hex return Hash in Geth types
func (d *Data) Hex() string {
	if d == nil {
		return "0x"
	}
	return hexutil.Encode(d.GetRaw())
}

// SetRaw sets bytes values for the Data message
func (d *Data) SetRaw(raw []byte) {
	d.Raw = raw
}

// Hash return Hash in Geth types
func (d *Data) Hash() ethcommon.Hash {
	return ethcommon.BytesToHash(d.GetRaw())
}
