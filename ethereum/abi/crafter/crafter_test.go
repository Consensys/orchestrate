package crafter

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/error"
)

const ERC20Payload = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"

type bindArgTest struct {
	typ *abi.Type
	arg string
	err bool
}

func testBindArg(t *testing.T, test bindArgTest) interface{} {
	boundArg, e := bindArg(test.typ, test.arg)
	assert.Equal(t, test.err, e != nil, e)
	if test.err {
		ie, ok := e.(*ierror.Error)
		assert.True(t, ok, "Error should cast to internal error")
		assert.Equal(t, "abi.crafter", ie.GetComponent(), "Component should be correct")
	}
	return boundArg
}

func TestBindArg(t *testing.T) {
	a := "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
	addrtype, _ := abi.NewType("address", "", nil)
	addr := testBindArg(t, bindArgTest{&addrtype, a, false}).(common.Address)
	assert.Equal(t, addr.Hex(), a, fmt.Sprintf("Expect bind %q but got %q", a, addr.Hex()))

	testBindArg(t, bindArgTest{&addrtype, "malformed-address", true})

	dectype, _ := abi.NewType("int", "", nil)
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
	ERC20TransferMethod, e := SignatureToMethod("transfer(address,uint256)")
	assert.Nil(t, e, "Parse method sig: should parse method signature")

	CustomMethod, e := SignatureToMethod("custom(address,bytes,uint256,uint17,bool,bytes,bytes16)")
	assert.Nil(t, e, "Parse method sig: should parse method signature")

	var (
		_to    = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_value = "0x2386f26fc10000"
	)
	_, e = bindArgs(ERC20TransferMethod, _to, _value)
	assert.Nil(t, e, "Prepare Args: should prepare args")

	_, e = bindArgs(ERC20TransferMethod, _to)
	assert.NotNil(t, e, "Parse method signature should fail")
	ie, ok := e.(*ierror.Error)
	assert.True(t, ok, "Error should cast to internal error")
	assert.Equal(t, "abi.crafter", ie.GetComponent(), "Component should be correct")

	var (
		_address = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_bytesA  = "0x2386f26fc10000"
		_uint256 = "0x6009608a02a7a15fd6689d6dad560c44e9ab61ff"
		_uint17  = "0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af"
		_bool    = "0x1"
		_bytesB  = "0xa1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b5"
		_bytes16 = "0xa1b2c3d4e5f67890"
	)
	_, e = bindArgs(CustomMethod, _address, _bytesA, _uint256, _uint17, _bool, _bytesB, _bytes16)
	assert.Nil(t, e, "Prepare Args: should prepare args")
}

func TestPayloadCrafter(t *testing.T) {
	ERC20TransferMethod, e := SignatureToMethod("transfer(address,uint256)")
	assert.Nil(t, e, "Parse method signature: received error")

	CustomMethod, e := SignatureToMethod("custom(address,bytes,uint256,uint17,bool,bytes,bytes16)")
	assert.Nil(t, e, "Parse method signature: received error")

	c := PayloadCrafter{}
	var (
		_to    = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_value = "0x2386f26fc10000"
	)
	data, e := c.CraftCall(ERC20TransferMethod, _to, _value)
	assert.Nil(t, e, "Craft: received error")

	assert.Equal(t, hexutil.Encode(data), ERC20Payload, "Craft: expected equal payload")

	var (
		_address = "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"
		_bytesA  = "0x2386f26fc10000"
		_uint256 = "0x6009608a02a7a15fd6689d6dad560c44e9ab61ff"
		_uint17  = "0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af"
		_bool    = "0x1"
		_bytesB  = "0xa1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b5"
		_bytes16 = "0xa1b2c3d4e5f67890"
	)

	data, e = c.CraftCall(CustomMethod, _address, _bytesA, _uint256, _uint17, _bool, _bytesB, _bytes16)
	assert.Nil(t, e, "Craft: received error")

	expected := "0x1db71ad9000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb68400000000000000000000000000000000000000000000000000000000000000e00000000000000000000000006009608a02a7a15fd6689d6dad560c44e9ab61ff000000000000000000000dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001200000000000000000a1b2c3d4e5f678900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000072386f26fc10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021a1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b500000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, hexutil.Encode(data), expected, "Craft: expected equal payload")
}

var testCrafterData = []struct {
	to    string
	value string
}{
	{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"},
	{},
}

func TestPayloadCrafterConcurrent(t *testing.T) {
	ERC20TransferMethod, e := SignatureToMethod("transfer(address,uint256)")
	assert.Nil(t, e, "Parse method signature: received error")

	c := PayloadCrafter{}
	rounds := 1000
	raws := make(chan []byte, rounds)
	wg := &sync.WaitGroup{}
	for i := 1; i <= rounds; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			raw, e := c.CraftCall(ERC20TransferMethod, testCrafterData[i%2].to, testCrafterData[i%2].value)
			// Test as been designed such as 1 out of 6 entry are valid for a credit
			if e == nil {
				raws <- raw
			}
		}(i)

	}
	wg.Wait()
	close(raws)
	assert.Equal(t, len(raws), rounds/2, "PayloadCrafter: expected specific crafts number")

	payload := "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	for data := range raws {
		assert.Equal(t, hexutil.Encode(data), payload, "Craft: expected equal payload")
	}
}

func TestPayloadCrafterArray(t *testing.T) {
	ArrayInput, e := SignatureToMethod("FunctionTest(uint256[3])")
	assert.Nil(t, e, "Parse method signature: received error")

	c := PayloadCrafter{}
	var (
		_array = "[\"0x1\",\"0x2\",\"0x3\"]"
	)
	data, e := c.CraftCall(ArrayInput, _array)
	assert.Nil(t, e, "Craft: received error")

	expected := "0x71cc037a000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003"
	assert.Equal(t, hexutil.Encode(data), expected, "Craft: expected equal payload")
}

func TestPayloadCrafterArrayAddress(t *testing.T) {
	ArrayAddressInput, e := SignatureToMethod("FunctionTest(address[3])")
	assert.Nil(t, e, "Parse method signature: received error")

	c := PayloadCrafter{}
	var (
		_array = "[\"0xca35b7d915458ef540ade6068dfe2f44e8fa733c\",\"0x14723a09acff6d2a60dcdf7aa4aff308fddc160c\",\"0x4b0897b0513fdc7c541b6d9d7e929c4e5364d2db\"]"
	)
	data, e := c.CraftCall(ArrayAddressInput, _array)
	assert.Nil(t, e, "Craft: received error")

	var expected = "0x620a6a89000000000000000000000000ca35b7d915458ef540ade6068dfe2f44e8fa733c00000000000000000000000014723a09acff6d2a60dcdf7aa4aff308fddc160c0000000000000000000000004b0897b0513fdc7c541b6d9d7e929c4e5364d2db"
	assert.Equal(t, hexutil.Encode(data), expected, "Craft: expected equal payload")
}

func TestPayloadCrafterSliceAddress(t *testing.T) {
	ArrayAddressInput, err := SignatureToMethod("FunctionTest(address[])")
	assert.Nil(t, err, "Parse method signature: received error")

	c := PayloadCrafter{}
	var (
		_array = "[\"0xca35b7d915458ef540ade6068dfe2f44e8fa733c\",\"0x14723a09acff6d2a60dcdf7aa4aff308fddc160c\",\"0x4b0897b0513fdc7c541b6d9d7e929c4e5364d2db\"]"
	)
	data, err := c.CraftCall(ArrayAddressInput, _array)
	assert.Nil(t, err, "Craft: received error")

	var expected = "0x8f2df58300000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000ca35b7d915458ef540ade6068dfe2f44e8fa733c00000000000000000000000014723a09acff6d2a60dcdf7aa4aff308fddc160c0000000000000000000000004b0897b0513fdc7c541b6d9d7e929c4e5364d2db"
	assert.Equal(t, hexutil.Encode(data), expected, "Craft: expected equal payload")
}

func newMethod(methodABI []byte) abi.Method {
	var method abi.Method
	_ = json.Unmarshal(methodABI, &method)
	return method
}

func TestSignatureToMethod(t *testing.T) {
	var EmptyMethod = newMethod([]byte(`{
		"inputs": [
		],
		"rawName": "empty"
	}`))

	var ERC20TransferMethod = newMethod([]byte(`{
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "uint256"
			}
		],
		"rawName": "transfer"
	}`))

	tests := []struct {
		sig    string
		result *abi.Method
		err    bool
	}{
		// {sig: "FunctionTest(address[3])", result: abi.Method{}, err: nil},
		{sig: "Malformed", result: nil, err: true},
		{sig: "", result: nil, err: true},
		{sig: "()Malformed", result: nil, err: true},
		{sig: "Malformed(,)", result: nil, err: true},
		{sig: "Malformed(address,uint)", result: nil, err: true},
		{sig: "Malformed(unknown)", result: nil, err: true},
		// {sig: "Malformed(address,", result: nil, err: true}, TODO: this test should error
		{sig: "empty()", result: &EmptyMethod, err: false},
		{sig: "transfer(address,uint256)", result: &ERC20TransferMethod, err: false},
	}

	for _, test := range tests {
		result, e := SignatureToMethod(test.sig)
		assert.Equal(t, test.err, e != nil, test.sig, result, e)
		if e != nil {
			ie, ok := e.(*ierror.Error)
			assert.True(t, ok, "Error should cast to internal error")
			assert.Equal(t, "abi.crafter", ie.GetComponent(), "Component should be correct")
		}
		assert.Equal(t, test.result, result)
	}
}
