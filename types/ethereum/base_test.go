package ethereum

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var invalidAddress = []byte{1}

var address1 = ethcommon.BytesToAddress([]byte{0}).Bytes()
var address1Str = "0x0000000000000000000000000000000000000000"

var address2 = ethcommon.BytesToAddress([]byte{1}).Bytes()
var address2Str = "0x0000000000000000000000000000000000000001"

var hash1 = ethcommon.BytesToHash([]byte{1}).Bytes()
var hash1Str = "0x0000000000000000000000000000000000000000000000000000000000000001"

var hash2 = ethcommon.BytesToHash([]byte{2}).Bytes()
var hash2Str = "0x0000000000000000000000000000000000000000000000000000000000000002"

func TestAccount(t *testing.T) {
	a := HexToAccount(address1Str)

	assert.Equal(t, a.Address().String(), address1Str, "Address has incorrect value")
	assert.Equal(t, a.Hex(), address1Str, "Address has incorrect value")
}

func TestCreateValidAddress(t *testing.T) {
	a := NewAccount(address1)

	assert.Equal(t, address1Str, a.Address().String(), "Address has incorrect value")
}

func TestCreateInvalidAddress(t *testing.T) {
	assert.PanicsWithValue(t, "address should be 20 bytes long", func() {
		NewAccount(invalidAddress)
	})
}

func TestSetValidAddress(t *testing.T) {
	a := NewAccount(address1)

	a.SetAddress(address2)

	assert.Equal(t, address2Str, a.Address().String(), "value ")
}

func TestSetInvalidAddress(t *testing.T) {
	a := NewAccount(address1)

	assert.PanicsWithValue(t, "address should be 20 bytes long", func() {
		a.SetAddress(invalidAddress)
	})
}

func TestHash(t *testing.T) {
	h := NewHash(hash1)

	assert.Equal(t, hash1Str, h.Hex(), "#1: value should match")

	h.SetValue(hash2)
	assert.Equal(t, hash2Str, h.Hex(), "#2: value should match")
}

func TestCreateInvalidHash(t *testing.T) {
	assert.PanicsWithValue(t, "hash should be 32 bytes long", func() {
		NewHash([]byte{1})
	})
}

func TestHexToHash(t *testing.T) {
	h := HexToHash(hash1Str)

	assert.Equal(t, hash1Str, h.Hex(), "value should match")
}

func TestQuantity(t *testing.T) {
	q := IntToQuantity(0)
	assert.Equal(t, "0", q.Value().Text(10), "#1: value should match")

	q.SetIntValue(10)
	assert.Equal(t, "10", q.Value().Text(10), "#1: value should match")

	q.SetRawValue(hexutil.MustDecode("0xff"))
	assert.Equal(t, "255", q.Value().Text(10), "#1: value should match")

	q.SetRawValue([]byte{0x01, 0x12})
	assert.Equal(t, "274", q.Value().Text(10), "#1: value should match")
}

func TestData(t *testing.T) {
	d := NewData([]byte{1})
	assert.Equal(t, d.Hex(), "0x01", "value should match")

	d.SetRaw([]byte{2})
	assert.Equal(t, d.Hex(), "0x02", "value should match")
}

func TestHexToData(t *testing.T) {
	d := HexToData("0x01")
	assert.Equal(t, d.Hex(), "0x01", "value should match")
	assert.Equal(t, d.Hash(), ethcommon.BytesToHash([]byte{1}), "value should match")
}
