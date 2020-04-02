// +build unit

package abi

import (
	"fmt"
	"math/big"
	"reflect"
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
	err error
}

func testBindArg(t *testing.T, test bindArgTest) interface{} {
	boundArg, e := BindArg(test.typ, test.arg)
	assert.Equal(t, test.err, e)
	return boundArg
}

func TestBindBytesArg(t *testing.T) {
	bytes1type, _ := abi.NewType("bytes1", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes1type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes1")})
	bytes2type, _ := abi.NewType("bytes2", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes2type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes2")})
	bytes3type, _ := abi.NewType("bytes3", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes3type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes3")})
	bytes4type, _ := abi.NewType("bytes4", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes4type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes4")})
	bytes5type, _ := abi.NewType("bytes5", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes5type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes5")})
	bytes6type, _ := abi.NewType("bytes6", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes6type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes6")})
	bytes7type, _ := abi.NewType("bytes7", "", nil)
	_ = testBindArg(t, bindArgTest{&bytes7type, "0xa1b2c3d4e5f67890", errors.InvalidArgError("invalid fixed bytes 0xa1b2c3d4e5f67890 of size 8 - too big for bytes7")})
	_ = testBindArg(t, bindArgTest{&bytes7type, "0xInvalid", errors.InvalidArgError("invalid fixed bytes \"0xInvalid\"")})

	byte1Array := testBindArg(t, bindArgTest{&bytes1type, "0xa1", nil}).([1]byte)
	expected := "0xa1"
	assert.Equal(t, hexutil.Encode(byte1Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte1Array[:])))

	bytes8type, _ := abi.NewType("bytes8", "", nil)
	byte8Array := testBindArg(t, bindArgTest{&bytes8type, "0xa1b2c3d4e5f67890", nil}).([8]byte)
	expected = "0xa1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte8Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte8Array[:])))

	bytes16type, _ := abi.NewType("bytes16", "", nil)
	byte16Array := testBindArg(t, bindArgTest{&bytes16type, "0xa1b2c3d4e5f67890", nil}).([16]byte)
	expected = "0x0000000000000000a1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte16Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte16Array[:])))

	bytes17type, _ := abi.NewType("bytes17", "", nil)
	byte17Array := testBindArg(t, bindArgTest{&bytes17type, "0xa1b2c3d4e5f67890", nil}).([17]byte)
	expected = "0x000000000000000000a1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte17Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte17Array[:])))

	bytes32type, _ := abi.NewType("bytes32", "", nil)
	byte32Array := testBindArg(t, bindArgTest{&bytes32type, "0xa1b2c3d4e5f67890", nil}).([32]byte)
	expected = "0x000000000000000000000000000000000000000000000000a1b2c3d4e5f67890"
	assert.Equal(t, hexutil.Encode(byte32Array[:]), expected, fmt.Sprintf("Expect bind to %v but got %v", expected, hexutil.Encode(byte32Array[:])))
}

func TestBindArg(t *testing.T) {
	testSet := []struct {
		name           string
		typ            string
		arg            string
		expectedOutput func() interface{}
		err            error
	}{
		{
			"",
			"int24",
			"-42",
			func() interface{} { return big.NewInt(-42) },
			nil,
		},
		{
			"",
			"int24",
			"-0x2A",
			func() interface{} { return big.NewInt(-42) },
			nil,
		},
		{
			"",
			"int24",
			"42",
			func() interface{} { return big.NewInt(42) },
			nil,
		},
		{
			"",
			"int24",
			"0x2A",
			func() interface{} { return big.NewInt(42) },
			nil,
		},
		{
			"",
			"uint64",
			"0x1A00300789A",
			func() interface{} { return uint64(1786756757658) },
			nil,
		},
		{
			"",
			"int64",
			"-0x1A00300789A",
			func() interface{} { return int64(-1786756757658) },
			nil,
		},
		{
			"",
			"uint64",
			"1786756757658",
			func() interface{} { return uint64(1786756757658) },
			nil,
		},
		{
			"",
			"int64",
			"-1786756757658",
			func() interface{} { return int64(-1786756757658) },
			nil,
		},
		{
			"",
			"uint8",
			"0x1A00300789A",
			func() interface{} { return nil },
			errors.InvalidArgError("could not parse number 0x1A00300789A - got \"strconv.ParseUint: parsing \\\"1A00300789A\\\": value out of range\""),
		},
		{
			"",
			"int8",
			"-0x1A00300789A",
			func() interface{} { return nil },
			errors.InvalidArgError("could not parse number -0x1A00300789A - got \"strconv.ParseInt: parsing \\\"-1A00300789A\\\": value out of range\""),
		},
		{
			"",
			"uint8",
			"1786756757658",
			func() interface{} { return nil },
			errors.InvalidArgError("could not parse number 1786756757658 - got \"strconv.ParseUint: parsing \\\"1786756757658\\\": value out of range\""),
		},
		{
			"",
			"int8",
			"-1786756757658",
			func() interface{} { return nil },
			errors.InvalidArgError("could not parse number -1786756757658 - got \"strconv.ParseInt: parsing \\\"-1786756757658\\\": value out of range\""),
		},
		{
			"",
			"uint8",
			"18",
			func() interface{} { return uint8(18) },
			nil,
		},
		{
			"",
			"uint8",
			"-18",
			func() interface{} { return nil },
			errors.InvalidArgError("did not expected negative value -18 for uint8"),
		},
		{
			"",
			"uint8",
			"",
			func() interface{} { return nil },
			errors.InvalidArgError("did not expected empty uint/int value"),
		},
		{
			"",
			"uint8",
			"0x",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid number - no value after 0x prefix"),
		},
		{
			"",
			"uint8",
			"0x01",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid number - got: hex number with leading zero digits"),
		},
		{
			"",
			"int24",
			"-0x01",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid number - got: hex number with leading zero digits"),
		},
		{
			"",
			"address",
			"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
			func() interface{} { return common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684") },
			nil,
		},
		{
			"",
			"address",
			"malformed-address",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid ethereum address \"malformed-address\""),
		},
		{
			"",
			"int256",
			"0x400",
			func() interface{} { return big.NewInt(1024) },
			nil,
		},
		{
			"",
			"int256",
			"1024",
			func() interface{} { return big.NewInt(1024) },
			nil,
		},
		{
			"",
			"int256",
			"0x34fg",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid int256 0x34fg"),
		},
		{
			"",
			"bool",
			"0x1",
			func() interface{} { return true },
			nil,
		},
		{
			"",
			"bool",
			"true",
			func() interface{} { return true },
			nil,
		},
		{
			"",
			"bool",
			"1",
			func() interface{} { return true },
			nil,
		},
		{
			"",
			"bool",
			"0x0",
			func() interface{} { return false },
			nil,
		},
		{
			"",
			"bool",
			"false",
			func() interface{} { return false },
			nil,
		},
		{
			"",
			"bool",
			"0",
			func() interface{} { return false },
			nil,
		},
		{
			"",
			"bool",
			"0x2",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid boolean \"0x2\" (expected one of [\"0x0\" \"false\" \"0\" \"0x1\" \"true\" \"1\"])"),
		},
		{
			"",
			"bytes",
			"0xabcd",
			func() interface{} { b, _ := hexutil.Decode("0xabcd"); return b },
			nil,
		},
		{
			"",
			"bytes",
			"0xabcg",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid bytes \"0xabcg\""),
		},
		{
			"",
			"string",
			"test",
			func() interface{} { return "test" },
			nil,
		},
		{
			"",
			"bool[]",
			"[\"true\", \"0x0\"]",
			func() interface{} { return []bool{true, false} },
			nil,
		},
		{
			"",
			"bool[2]",
			"[\"false\", \"0x1\"]",
			func() interface{} { return []bool{false, true} },
			nil,
		},
		{
			"",
			"bool[3]",
			"[\"false\", \"0x1\"]",
			func() interface{} { return nil },
			errors.InvalidArgError("invalid size array \"[\\\"false\\\", \\\"0x1\\\"]\" (expected length 3 but got 2)"),
		},
		{
			"",
			"bool[3]",
			"\"false\", \"0x1\"",
			func() interface{} { return nil },
			errors.InvalidArgError("could not parse array \"false\", \"0x1\" for bool[3] - got invalid character ',' after top-level value"),
		},
		{
			"",
			"address[3]",
			"[\"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae\", \"0x5ed8cee6b63b1c6afce3ad7c92f4fd7e1b8fad9f\", \"0x9ee457023bb3de16d51a003a247baead7fce313d\"]",
			func() interface{} {
				return []common.Address{
					common.HexToAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"),
					common.HexToAddress("0x5ed8cee6b63b1c6afce3ad7c92f4fd7e1b8fad9f"),
					common.HexToAddress("0x9ee457023bb3de16d51a003a247baead7fce313d"),
				}
			},
			nil,
		},
		{
			"",
			"address[3]",
			"[\"0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae\", \"invalidAddress\", \"0x9ee457023bb3de16d51a003a247baead7fce313d\"]",
			func() interface{} {
				return nil
			},
			errors.InvalidArgError("invalid ethereum address \"invalidAddress\""),
		},
		{
			"",
			"tuple",
			"",
			func() interface{} { return nil },
			errors.FeatureNotSupportedError("solidity tuple not supported yet"),
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			abiType, err := abi.NewType(test.typ, "", nil)
			assert.NoError(t, err, "%s should be a valid type", test.typ)
			boundArg, err := BindArg(&abiType, test.arg)
			if err != nil || test.err != nil {
				assert.Equal(t, test.err, err, "Should return Solidity error")
				return
			}
			assert.True(t, reflect.DeepEqual(boundArg, test.expectedOutput()), "%v and %v should be equal", boundArg, test.expectedOutput)
		})
	}
}

func TestBindArgs(t *testing.T) {
	testSet := []struct {
		name           string
		methodSig      string
		args           []string
		expectedOutput func() []interface{}
		err            error
	}{
		{
			"",
			"constructor(string,string,uint8)",
			[]string{"FungibleToken", "FTK", "18"},
			func() []interface{} {
				return []interface{}{
					"FungibleToken",
					"FTK",
					uint8(18),
				}
			},
			nil,
		},
		{
			"",
			"transfer(address,uint256)",
			[]string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"},
			func() []interface{} {
				b, _ := new(big.Int).SetString("2386f26fc10000", 16)
				return []interface{}{
					common.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"),
					b,
				}
			},
			nil,
		},
		{
			"",
			"transfer(address,uint256)",
			[]string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
			func() []interface{} { return nil },
			errors.InvalidArgsCountError("invalid arguments count (expected 2 but got 1)"),
		},
		{
			"",
			"test(uint8)",
			[]string{"-19"},
			func() []interface{} { return nil },
			errors.InvalidArgError("did not expected negative value -19 for uint8"),
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			method, _ := ParseMethodSignature(test.methodSig)
			output, err := BindArgs(&method.Inputs, test.args...)
			if err != nil || test.err != nil {
				assert.Equal(t, err, test.err, "Should return Solidity error")
				return
			}
			assert.True(t, reflect.DeepEqual(output, test.expectedOutput()), "%v and %v should be equal", output, test.expectedOutput)
		})
	}

}
