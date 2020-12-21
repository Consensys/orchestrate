// +build unit

package abi

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	testSuite := []struct {
		sig            string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			sig:            "testMethod(int8)",
			args:           []string{strings.ReplaceAll(strconv.FormatInt(int64(-15), 16), "-", "-0x")},
			expectedOutput: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int8)",
			args:           []string{"0x" + strconv.FormatInt(int64(15), 16)},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000000000000000f",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int16)",
			args:           []string{strings.ReplaceAll(strconv.FormatInt(int64(-777), 16), "-", "-0x")},
			expectedOutput: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffcf7",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int16)",
			args:           []string{"0x" + strconv.FormatInt(int64(777), 16)},
			expectedOutput: "0x0000000000000000000000000000000000000000000000000000000000000309",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int32)",
			args:           []string{strings.ReplaceAll(strconv.FormatInt(int64(-666), 16), "-", "-0x")},
			expectedOutput: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd66",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int32)",
			args:           []string{"0x" + strconv.FormatInt(int64(666), 16)},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000000000000029a",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int64)",
			args:           []string{strings.ReplaceAll(strconv.FormatInt(int64(-9876543234567), 16), "-", "-0x")},
			expectedOutput: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffff70470261df9",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int64)",
			args:           []string{"0x" + strconv.FormatInt(int64(9876543234567), 16)},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000008fb8fd9e207",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"0x" + strconv.FormatUint(uint64(45), 16)},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000000000000002d",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint16)",
			args:           []string{"0x" + strconv.FormatUint(uint64(888), 16)},
			expectedOutput: "0x0000000000000000000000000000000000000000000000000000000000000378",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint32)",
			args:           []string{"0x" + strconv.FormatUint(uint64(888), 16)},
			expectedOutput: "0x0000000000000000000000000000000000000000000000000000000000000378",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint64)",
			args:           []string{"0x" + strconv.FormatUint(uint64(888), 16)},
			expectedOutput: "0x0000000000000000000000000000000000000000000000000000000000000378",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int256)",
			args:           []string{hexutil.EncodeBig(big.NewInt(456744578797645890))},
			expectedOutput: "0x0000000000000000000000000000000000000000000000000656aec64451b042",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int256)",
			args:           []string{fmt.Sprintf("%#x", big.NewInt(-456744578797645890))},
			expectedOutput: "0xfffffffffffffffffffffffffffffffffffffffffffffffff9a95139bbae4fbe",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint256)",
			args:           []string{hexutil.EncodeBig(big.NewInt(898555348797645890))},
			expectedOutput: "0x0000000000000000000000000000000000000000000000000c784f54381c6442",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int40)",
			args:           []string{hexutil.EncodeBig(big.NewInt(8764243))},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000000000085bb53",
			expectedError:  false,
		},
		{
			sig:            "testMethod(int40)",
			args:           []string{fmt.Sprintf("%#x", big.NewInt(-8764243))},
			expectedOutput: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7a44ad",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint24)",
			args:           []string{fmt.Sprintf("%#x", big.NewInt(7565467))},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000000000073709b",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{""},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"a"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"-0xa"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(int8)",
			args:           []string{"0s"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"0x"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"0x07ab"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(int8)",
			args:           []string{"0x7Zab"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"0x7Zab"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"-0xa"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(uint8)",
			args:           []string{"0x9a95139bbae4fbe"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(int8)",
			args:           []string{"0x9a95139bbae4fbe"},
			expectedOutput: "",
			expectedError:  true,
		},
		{
			sig:            "testMethod(address[3])",
			args:           []string{"[\"0xca35b7d915458ef540ade6068dfe2f44e8fa733c\",\"0x14723a09acff6d2a60dcdf7aa4aff308fddc160c\",\"0x4b0897b0513fdc7c541b6d9d7e929c4e5364d2db\"]"},
			expectedOutput: "0x000000000000000000000000ca35b7d915458ef540ade6068dfe2f44e8fa733c00000000000000000000000014723a09acff6d2a60dcdf7aa4aff308fddc160c0000000000000000000000004b0897b0513fdc7c541b6d9d7e929c4e5364d2db",
			expectedError:  false,
		},
		{
			sig:            "testMethod(uint256[3])",
			args:           []string{"[\"0x1\",\"0x2\",\"0x3\"]"},
			expectedOutput: "0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003",
			expectedError:  false,
		},
		{
			sig:            "testMethod(address[])",
			args:           []string{"[\"0xca35b7d915458ef540ade6068dfe2f44e8fa733c\",\"0x14723a09acff6d2a60dcdf7aa4aff308fddc160c\",\"0x4b0897b0513fdc7c541b6d9d7e929c4e5364d2db\"]"},
			expectedOutput: "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000ca35b7d915458ef540ade6068dfe2f44e8fa733c00000000000000000000000014723a09acff6d2a60dcdf7aa4aff308fddc160c0000000000000000000000004b0897b0513fdc7c541b6d9d7e929c4e5364d2db",
			expectedError:  false,
		},
		{
			sig: "testMethod(address,uint256)",
			args: []string{
				"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
				"0x2386f26fc10000",
			},
			expectedOutput: "000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000",
			expectedError:  false,
		},
		{
			sig: "testMethod(address,bytes,uint256,uint17,bool,bytes,bytes16)",
			args: []string{
				"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
				"0x2386f26fc10000",
				"0x6009608a02a7a15fd6689d6dad560c44e9ab61ff",
				"0xdd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af",
				"0x1",
				"0xa1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b5",
				"0xa1b2c3d4e5f67890",
			},
			expectedOutput: "0x000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb68400000000000000000000000000000000000000000000000000000000000000e00000000000000000000000006009608a02a7a15fd6689d6dad560c44e9ab61ff000000000000000000000dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775af000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001200000000000000000a1b2c3d4e5f678900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000072386f26fc10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021a1a45fabb381e6ab02448013f651fa0792c3fa05b38771f161cb8f7ebdbee973b500000000000000000000000000000000000000000000000000000000000000",
			expectedError:  false,
		},
		{
			sig: "testMethod(address,uint256)",
			args: []string{
				"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
			},
			expectedOutput: "",
			expectedError:  true,
		},
	}

	for i, test := range testSuite {
		sig, _ := ParseMethodSignature(test.sig)
		data, err := Pack(sig, test.args...)
		if test.expectedError {
			assert.Error(t, err, "GetNonce (%d/%d): should get an error", i+1, len(testSuite))
			return
		}
		assert.NoError(t, err, "GetNonce (%d/%d):: should not get an error", i+1, len(testSuite))
		assert.Equal(t, test.expectedOutput, hexutil.Encode(data), "GetNonce (%d/%d): expected equal payload", i+1, len(testSuite))
	}
}
