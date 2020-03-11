package abi

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type bindArgTest struct {
	typ *abi.Type
	arg string
	err bool
}

func testBindArg(t *testing.T, test bindArgTest) interface{} {
	boundArg, e := BindArg(test.typ, test.arg)
	assert.Equal(t, test.err, e != nil, e)
	if test.err {
		assert.True(t, errors.IsSolidityError(e), "Should return Solidity error")
	}
	return boundArg
}

func TestBindArg(t *testing.T) {
	a := "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
	addrtype, _ := abi.NewType("address", "", nil)
	addr := testBindArg(t, bindArgTest{&addrtype, a, false}).(common.Address)
	assert.Equal(t, addr.Hex(), a, fmt.Sprintf("Expect bind %q but got %q", a, addr.Hex()))

	testBindArg(t, bindArgTest{&addrtype, "malformed-address", true})

	dectype, _ := abi.NewType("int256", "", nil)
	dec := testBindArg(t, bindArgTest{&dectype, "0x400", false}).(*big.Int)
	assert.Equal(t, dec.Int64(), int64(1024), fmt.Sprintf("Expect bind to %v but got %v", 1024, dec.Int64()))

	testBindArg(t, bindArgTest{&dectype, "0x34fg", true})

	booltype, _ := abi.NewType("bool", "", nil)
	boolean := testBindArg(t, bindArgTest{&booltype, "0x1", false}).(bool)
	assert.Truef(t, boolean, fmt.Sprintf("Expect bind to %v but got %v", true, false))

	testBindArg(t, bindArgTest{&booltype, "0x2", true})

	bytestype, _ := abi.NewType("bytes", "", nil)
	byteSlice := testBindArg(t, bindArgTest{&bytestype, "0xabcd", false}).([]byte)
	assert.Equal(t, hexutil.Encode(byteSlice), "0xabcd", fmt.Sprintf("Expect bind to %v but got %v", "0xabcd", hexutil.Encode(byteSlice)))

	testBindArg(t, bindArgTest{&bytestype, "0xabcg", true})

	bytes1type, _ := abi.NewType("bytes1", "", nil)
	byte1Array := testBindArg(t, bindArgTest{&bytes1type, "0xa1b2c3d4e5f67890", false}).([1]byte)
	expected := "0xa1"
	assert.Equal(t, hexutil.Encode(byte1Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte1Array[:])))

	bytes8type, _ := abi.NewType("bytes8", "", nil)
	byte8Array := testBindArg(t, bindArgTest{&bytes8type, "0xa1b2c3d4e5f67890", false}).([8]byte)
	expected = "0xa1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte8Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte8Array[:])))

	bytes16type, _ := abi.NewType("bytes16", "", nil)
	byte16Array := testBindArg(t, bindArgTest{&bytes16type, "0xa1b2c3d4e5f67890", false}).([16]byte)
	expected = "0x0000000000000000a1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte16Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte16Array[:])))

	bytes17type, _ := abi.NewType("bytes17", "", nil)
	byte17Array := testBindArg(t, bindArgTest{&bytes17type, "0xa1b2c3d4e5f67890", false}).([17]byte)
	expected = "0x000000000000000000a1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte17Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte17Array[:])))

	bytes32type, _ := abi.NewType("bytes32", "", nil)
	byte32Array := testBindArg(t, bindArgTest{&bytes32type, "0xa1b2c3d4e5f67890", false}).([32]byte)
	expected = "0x000000000000000000000000000000000000000000000000a1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte32Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte32Array[:])))
}

func TestBindArgs(t *testing.T) {
	var (
		_to    = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_value = "0x2386f26fc10000"
	)

	method, _ := ParseMethodSignature("transfer(address,uint256)")
	_, e := BindArgs(&method.Inputs, _to, _value)
	assert.NoError(t, e, "BindArgs: should bind args properly args")

	_, e = BindArgs(&method.Inputs, _to)
	assert.Error(t, e, "Parse method signature should fail")
	assert.True(t, errors.IsSolidityError(e), "Should return Solidity error")
}
