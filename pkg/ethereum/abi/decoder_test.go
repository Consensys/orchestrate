// +build unit

package abi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/go-ethereum/v1_9_12/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ConsenSys/orchestrate/pkg/types/ethereum"
)

func newEvent(eventABI []byte) *abi.Event {
	var event abi.Event
	_ = json.Unmarshal(eventABI, &event)
	return &event
}

// TODO: Test with all types
func TestFormatIndexedArg(t *testing.T) {

	for i, test := range []struct {
		argType        string
		arg            ethcommon.Hash
		expectedOutput string
	}{
		{
			"string",
			ethcommon.HexToHash("0x41e406698d040bb44cf693b3dc50c37cf3c854c422d2645b1101662741fbaa88"),
			"41e406698d040bb44cf693b3dc50c37cf3c854c422d2645b1101662741fbaa88",
		},
		{
			"bool",
			ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			"false",
		},
		{
			"bool",
			ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
			"true",
		},
		{
			"address",
			ethcommon.HexToHash("0x0000000000000000000000008dd688660ec0babd0b8a2f2de3232645f73cc5eb"),
			"0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		},
		{
			"bytes32",
			ethcommon.HexToHash("0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088"),
			"0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088",
		},
		{
			"uint256",
			ethcommon.BigToHash(big.NewInt(1)),
			"1",
		},
	} {
		typeArg, _ := abi.NewType(test.argType, "", nil)
		output, _ := FormatIndexedArg(&typeArg, test.arg)

		if test.expectedOutput != output {
			t.Errorf("TestFormatIndexedArg (input %d): expected %q but got %q", i, test.expectedOutput, output)
		}
	}

}

// TODO: Update int type
func TestFormatNonIndexedArg(t *testing.T) {

	for i, test := range []struct {
		argType        abi.Type
		arg            interface{}
		expectedOutput string
	}{
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			uint8(2),
			"2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]uint8{1, 2},
			"[1 2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			uint16(2),
			"2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]uint16{1, 2},
			"[1 2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			uint32(2),
			"2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]uint32{1, 2},
			"[1 2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			uint64(2),
			"2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]uint64{1, 2},
			"[1 2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			big.NewInt(2),
			"2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]*big.Int{big.NewInt(1), big.NewInt(2)},
			"[1 2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			int8(-2),
			"-2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]int8{-1, -2},
			"[-1 -2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			int16(-2),
			"-2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]int16{-1, -2},
			"[-1 -2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			int32(-2),
			"-2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]int32{-1, -2},
			"[-1 -2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			int64(-2),
			"-2",
		},
		{
			abi.Type{Type: reflect.TypeOf(&big.Int{}), T: abi.IntTy},
			[]int64{-1, -2},
			"[-1 -2]",
		},
		{
			abi.Type{Type: reflect.TypeOf(byte(0)), T: abi.FixedBytesTy},
			[32]byte{1},
			"0x0100000000000000000000000000000000000000000000000000000000000000",
		},
		{
			abi.Type{Type: reflect.TypeOf(ethcommon.Address{}), T: abi.AddressTy},
			ethcommon.HexToAddress("01"),
			"0x0000000000000000000000000000000000000001",
		},
	} {
		output, _ := FormatNonIndexedArg(&test.argType, test.arg)

		if test.expectedOutput != output {
			t.Errorf("FormatNonIndexedArg (input %d): expected mapping %q but got %q", i, test.expectedOutput, output)
		}
	}

}

func TestDecode(t *testing.T) {

	testSet := []struct {
		abi            []byte
		log            *ethereum.Log
		expectedOutput map[string]string
		expectedError  error
	}{
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"}`),
			&ethereum.Log{
				Data: "0x000000000000000000000000000000000000000000000001a055690d9db80000",
				Topics: []string{
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					"0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0",
					"0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f",
				},
			},
			map[string]string{
				"tokens": "30000000000000000000",
				"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
				"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":true,"name":"maker","type":"address"},{"indexed":false,"name":"taker","type":"address"},{"indexed":true,"name":"feeRecipient","type":"address"},{"indexed":false,"name":"makerToken","type":"address"},{"indexed":false,"name":"takerToken","type":"address"},{"indexed":false,"name":"filledMakerTokenAmount","type":"uint256"},{"indexed":false,"name":"filledTakerTokenAmount","type":"uint256"},{"indexed":false,"name":"paidMakerFee","type":"uint256"},{"indexed":false,"name":"paidTakerFee","type":"uint256"},{"indexed":true,"name":"tokens","type":"bytes32"},{"indexed":false,"name":"orderHash","type":"bytes32"}],"name":"LogFill","type":"event"}`),
			&ethereum.Log{
				Data: "0x000000000000000000000000e269e891a2ec8585a378882ffa531141205e92e9000000000000000000000000d7732e3783b0047aa251928960063f863ad022d8000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000032d26d12e980b6000000000000000000000000000000000000000000000000000000602d4ec2c348a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000306a9a7ecbd9446559a2c650b4cfc16d1fb615aa2b3f4f63078da6d021268440",
				Topics: []string{
					"0x0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3",
					"0x0000000000000000000000008dd688660ec0babd0b8a2f2de3232645f73cc5eb",
					"0x000000000000000000000000e269e891a2ec8585a378882ffa531141205e92e9",
					"0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088",
				},
			},
			map[string]string{
				"filledTakerTokenAmount": "6930282000000000000",
				"paidTakerFee":           "0",
				"tokens":                 "0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088",
				"feeRecipient":           "0xe269E891A2Ec8585a378882fFA531141205e92E9",
				"makerToken":             "0xD7732e3783b0047aa251928960063f863AD022D8",
				"takerToken":             "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"filledMakerTokenAmount": "240000000000000000000000",
				"paidMakerFee":           "0",
				"orderHash":              "0x306a9a7ecbd9446559a2c650b4cfc16d1fb615aa2b3f4f63078da6d021268440",
				"maker":                  "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
				"taker":                  "0xe269E891A2Ec8585a378882fFA531141205e92E9",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"TransferWithData","type":"event"}`),
			&ethereum.Log{
				Data: "0x0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000061000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				Topics: []string{
					"0xe8f0a47da72ca43153c7a5693a827aa8456f52633de9870a736e5605bff4af6d",
					"0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236",
					"0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236",
					"0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f",
				},
			},
			map[string]string{
				"operator":     "0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
				"from":         "0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
				"to":           "0xb5747835141b46f7C472393B31F8F5A57F74A44f",
				"value":        "3",
				"data":         "0x000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01",
				"operatorData": "0x",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"index_origin","type":"uint256"},{"indexed":false,"name":"transfer_id","type":"bytes32"},{"indexed":false,"name":"parameters","type":"bytes32[6]"}],"name":"PendingTransfer","type":"event"}`),
			&ethereum.Log{
				Data: "0x00000000000000000000000000000000000000000000000000000000000000015a4f2c3ad66af173634e1cc1e389232788ec41756ec2821b9a231f996c4faad00000000000000000000000008f371daa8a5325f53b754a7017ac3803382bc8470000000000000000000000003404370fddb2b0e79f2571e170b112a66f974fb95265736572766564000000000000000000000000000000000000000000000000000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f0000000000000000000000000000000000000000000000000000000000000005",
				Topics: []string{
					"0xd8589d63a2df3a19b774d092cc22aec68d0be6537da4f37a362fbba9f6296845",
				},
			},
			map[string]string{
				"index_origin": "1",
				"transfer_id":  "0x5a4f2c3ad66af173634e1cc1e389232788ec41756ec2821b9a231f996c4faad0",
				"parameters":   "[\"0x0000000000000000000000008f371daa8a5325f53b754a7017ac3803382bc847\",\"0x0000000000000000000000003404370fddb2b0e79f2571e170b112a66f974fb9\",\"0x5265736572766564000000000000000000000000000000000000000000000000\",\"0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f\",\"0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f\",\"0x0000000000000000000000000000000000000000000000000000000000000005\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"address[2]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
				Topics: []string{
					"0xc69e2952deb9c87385c05617318dac13d6e25ebcf771aa2ae600cc8b1c1c7c73",
				},
			},
			map[string]string{
				"array": "[\"0x0000000000000000000000000000000000000000\",\"0x0000000000000000000000000000000000000001\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"eventType","type":"string"},{"indexed":false,"name":"nomId","type":"uint256"},{"indexed":false,"name":"shipQty","type":"uint256"},{"indexed":false,"name":"shipper","type":"address"},{"indexed":false,"name":"participants","type":"address[]"}],"name":"ShipperIdentified","type":"event"}`),
			&ethereum.Log{
				Data: "0x00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000003e9000000000000000000000000000000000000000000000000000000000000000000000000000000000000000079ff9389f9a574917a0f14f752590dbbb5f0f01700000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000124e6f6d696e6174696f6e52656a6563746564000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000566d44fb1442d206439f23e09885ee365bb1e43f",
				Topics: []string{
					"0xfe5678ff8df37e23934e6875e6f96885f66eb0185376749e5e0c0e9a91d2f181",
				},
			},
			map[string]string{
				"eventType":    "NominationRejected",
				"nomId":        "1001",
				"shipQty":      "0",
				"shipper":      "0x79fF9389F9a574917a0f14F752590DBbB5f0F017",
				"participants": "[\"0x0000000000000000000000000000000000000000\",\"0x566D44FB1442d206439F23E09885Ee365Bb1E43f\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"bool[3]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
				Topics: []string{
					"0xe93cecbaa9dac518aeb3ecb88fe5e11971da9ddc0d88eaa71902ca0ac9d6afdd",
				},
			},
			map[string]string{
				"array": "[\"true\",\"false\",\"true\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"int256[3]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000006f",
				Topics: []string{
					"0x4e9f8c6a9ba6c36d9f69bdaeda487b34fcba584137b513dd615a59ded56c939a",
				},
			},
			map[string]string{
				"array": "[\"1\",\"11\",\"111\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"int256[3]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff91",
				Topics: []string{
					"0x4e9f8c6a9ba6c36d9f69bdaeda487b34fcba584137b513dd615a59ded56c939a",
				},
			},
			map[string]string{
				"array": "[\"-1\",\"-11\",\"-111\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"uint256[3]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x0000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000006300000000000000000000000000000000000000000000000000000000000003e7",
				Topics: []string{
					"0x39cc9b81f311e9bdf9c08720512a61f27e13fd23c9f03938c704e02a6145c45d",
				},
			},
			map[string]string{
				"array": "[\"9\",\"99\",\"999\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"uint256[3][]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000006f0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000de00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000021000000000000000000000000000000000000000000000000000000000000014d",
				Topics: []string{
					"0x85dfc0f9608903e59d1fc7374a0697adb9fcbc35e25888ede73ed16380ac6dbd",
				},
			},
			map[string]string{
				"array": "[\"[\\\"1\\\",\\\"11\\\",\\\"111\\\"]\",\"[\\\"2\\\",\\\"22\\\",\\\"222\\\"]\",\"[\\\"3\\\",\\\"33\\\",\\\"333\\\"]\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"bool[2][2][2]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				Topics: []string{
					"0xa18b22ca18dff5c2b530dcda530785cb179fd62c5b113971ca9014a86d716832",
				},
			},
			map[string]string{
				"array": "[\"[\\\"[\\\\\\\"true\\\\\\\",\\\\\\\"false\\\\\\\"]\\\",\\\"[\\\\\\\"false\\\\\\\",\\\\\\\"false\\\\\\\"]\\\"]\",\"[\\\"[\\\\\\\"false\\\\\\\",\\\\\\\"false\\\\\\\"]\\\",\\\"[\\\\\\\"false\\\\\\\",\\\\\\\"false\\\\\\\"]\\\"]\"]",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"bool[2][2][]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000",
				Topics: []string{
					"0xd5382ec33a47d9fd2a8bc6b79938a702d96b8d9caf99455323318181d9069162",
				},
			},
			map[string]string{
				"array": "[\"[\\\"[\\\\\\\"true\\\\\\\",\\\\\\\"false\\\\\\\"]\\\",\\\"[\\\\\\\"true\\\\\\\",\\\\\\\"false\\\\\\\"]\\\"]\",\"[\\\"[\\\\\\\"false\\\\\\\",\\\\\\\"false\\\\\\\"]\\\",\\\"[\\\\\\\"false\\\\\\\",\\\\\\\"false\\\\\\\"]\\\"]\",\"[\\\"[\\\\\\\"true\\\\\\\",\\\\\\\"false\\\\\\\"]\\\",\\\"[\\\\\\\"true\\\\\\\",\\\\\\\"false\\\\\\\"]\\\"]\"]",
			},
			nil,
		},
		{
			[]byte(`{"type":"event","name":"tuple","inputs":[{"indexed":false,"name":"t","type":"tuple","components":[{"type":"int256","name":"a"},{"type":"int256","name":"b"}]}]}`),
			&ethereum.Log{
				Data: "0x0000000000000000000000000000000000000000000000000000000000000001ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				Topics: []string{
					"0xc270a86bd76188052175715b1997fa46780b24c6892bbd0f6ff98066f83543a4",
				},
			},
			map[string]string{
				"t": "{\"A\":\"1\",\"B\":\"-1\"}",
			},
			nil,
		},
		{
			[]byte(`{"type":"event","name":"tuple","inputs":[
				{"type":"tuple","name":"s","components":[{"type":"uint256","name":"a"},{"type":"uint256[]","name":"b"},{"type":"tuple[]","name":"c","components":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"}]}]},
				{"type":"tuple","name":"t","components":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"}]},
				{"type":"uint256","name":"a"}
			]}`),
			&ethereum.Log{
				Data: "0x00000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001",
				Topics: []string{
					"0xc270a86bd76188052175715b1997fa46780b24c6892bbd0f6ff98066f83543a4",
				},
			},
			map[string]string{
				"s": "{\"A\":\"1\",\"B\":\"[\\\"1\\\",\\\"2\\\"]\",\"C\":\"[\\\"{\\\\\\\"X\\\\\\\":\\\\\\\"1\\\\\\\",\\\\\\\"Y\\\\\\\":\\\\\\\"2\\\\\\\"}\\\",\\\"{\\\\\\\"X\\\\\\\":\\\\\\\"2\\\\\\\",\\\\\\\"Y\\\\\\\":\\\\\\\"1\\\\\\\"}\\\"]\"}",
				"t": "{\"X\":\"0\",\"Y\":\"1\"}",
				"a": "1",
			},
			nil,
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"uint256[3]"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x0000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000006300000000000000000000000000000000000000000000000000000000000003e7",
				Topics: []string{
					"0x39cc9b81f311e9bdf9c08720512a61f27e13fd23c9f03938c704e02a6145c45d",
					"0x39cc9b81f311e9bdf9c08720512a61f27e13fd23c9f03938c704e02a6145c45d",
				},
			},
			nil,
			fmt.Errorf("decoder error: too many topics"),
		},
		{
			[]byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"array","type":"bool"}],"name":"EventTest","type":"event"}`),
			&ethereum.Log{
				Data: "0x0000000000000001000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000006300000000000000000000000000000000000000000000000000000000000003e7",
				Topics: []string{
					"0x39cc9b81f311e9bdf9c08720512a61f27e13fd23c9f03938c704e02a6145c45d",
				},
			},
			nil,
			fmt.Errorf("decoder error: cannot UnpackValues correctly"),
		},
	}

	for i, test := range testSet {
		event := newEvent(test.abi)
		decoded, err := Decode(event, test.log)
		if err != nil && test.expectedError == nil {
			t.Errorf("Decode: Expecting the following error %v but got %v", test.expectedError, err)
		}
		eq := reflect.DeepEqual(test.expectedOutput, decoded)
		if !eq {
			t.Errorf("Decode (%d/%d) %q: expected mapping %q but got %q", i+1, len(testSet), event.Name, test.expectedOutput, decoded)
		}

	}
}
