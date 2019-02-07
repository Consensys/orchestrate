package ethereum

import (
	"encoding/json"
	"math/big"
	"reflect"

	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

var ERC20ABI = []byte(`{
	"anonymous":false,
	"inputs":[
		{"indexed":true,"name":"from","type":"address"},
		{"indexed":true,"name":"to","type":"address"},
		{"indexed":false,"name":"tokens","type":"uint256"}
	],
	"name":"Transfer",
	"type":"event"
}`)

var testLogERC20ABI = &types.Log{
	Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000001a055690d9db80000"),
	Topics: []common.Hash{
		common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		common.HexToHash("0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0"),
		common.HexToHash("0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f"),
	},
}

var EventTestAbi = []byte(`{
	"anonymous":false,
	"inputs":[
	   {
		  "indexed":true,
		  "name":"maker",
		  "type":"address"
	   },
	   {
		  "indexed":false,
		  "name":"taker",
		  "type":"address"
	   },
	   {
		  "indexed":true,
		  "name":"feeRecipient",
		  "type":"address"
	   },
	   {
		  "indexed":false,
		  "name":"makerToken",
		  "type":"address"
	   },
	   {
		  "indexed":false,
		  "name":"takerToken",
		  "type":"address"
	   },
	   {
		  "indexed":false,
		  "name":"filledMakerTokenAmount",
		  "type":"uint256"
	   },
	   {
		  "indexed":false,
		  "name":"filledTakerTokenAmount",
		  "type":"uint256"
	   },
	   {
		  "indexed":false,
		  "name":"paidMakerFee",
		  "type":"uint256"
	   },
	   {
		  "indexed":false,
		  "name":"paidTakerFee",
		  "type":"uint256"
	   },
	   {
		  "indexed":true,
		  "name":"tokens",
		  "type":"bytes32"
	   },
	   {
		  "indexed":false,
		  "name":"orderHash",
		  "type":"bytes32"
	   }
	],
	"name":"LogFill",
	"type":"event"
 }`)

var EventTestTx = &types.Log{
	Data: hexutil.MustDecode("0x000000000000000000000000e269e891a2ec8585a378882ffa531141205e92e9000000000000000000000000d7732e3783b0047aa251928960063f863ad022d8000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000032d26d12e980b6000000000000000000000000000000000000000000000000000000602d4ec2c348a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000306a9a7ecbd9446559a2c650b4cfc16d1fb615aa2b3f4f63078da6d021268440"),
	Topics: []common.Hash{
		common.HexToHash("0x0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3"),
		common.HexToHash("0x0000000000000000000000008dd688660ec0babd0b8a2f2de3232645f73cc5eb"),
		common.HexToHash("0x000000000000000000000000e269e891a2ec8585a378882ffa531141205e92e9"),
		common.HexToHash("0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088"),
	},
}

func newEvent(eventABI []byte) *abi.Event {
	var event abi.Event
	json.Unmarshal(eventABI, &event)
	return &event
}

func TestDecodeERC20ABI(t *testing.T) {
	event := newEvent(ERC20ABI)

	decoded, _ := Decode(event, testLogERC20ABI)

	expected := map[string]string{
		"tokens": "30000000000000000000",
		"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
		"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
	}
	eq := reflect.DeepEqual(expected, decoded)
	if !eq {
		t.Errorf("Decode ERC20ABI: expected mapping %q but got %q", expected, decoded)
	}
}

func TestDecodeEventTest(t *testing.T) {
	event := newEvent(EventTestAbi)

	decoded, _ := Decode(event, EventTestTx)

	m := map[string]string{
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
	}
	eq := reflect.DeepEqual(m, decoded)
	if !eq {
		t.Errorf("Decode EventTest: expected mapping %q but got %q", m, decoded)
	}
}

// TODO: Test with all types
func TestFormatIndexedArg(t *testing.T) {

	for i, test := range []struct {
		argType        string
		arg            common.Hash
		expectedOutput string
	}{
		{
			"address",
			common.HexToHash("0x0000000000000000000000008dd688660ec0babd0b8a2f2de3232645f73cc5eb"),
			"0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
		},
		{
			"bytes32",
			common.HexToHash("0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088"),
			"0xf08499c9e419ea8c08c4b991f88632593fb36baf4124c62758acb21898711088",
		},
		{
			"uint256",
			common.BigToHash(big.NewInt(1)),
			"1",
		},
	} {
		typeArg, _ := abi.NewType(test.argType, nil)
		output, _ := FormatIndexedArg(typeArg, test.arg)

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
			abi.Type{Type: reflect.TypeOf(common.Address{}), T: abi.AddressTy},
			common.HexToAddress("01"),
			"0x0000000000000000000000000000000000000001",
		},
	} {
		output, _ := FormatNonIndexedArg(test.argType, test.arg)

		if test.expectedOutput != output {
			t.Errorf("FormatNonIndexedArg (input %d): expected mapping %q but got %q", i, test.expectedOutput, output)
		}
	}

}

var ERC1400TransferWithDataAbi = []byte(`{
	"anonymous": false,
	"inputs": [
		{
			"indexed": true,
			"name": "operator",
			"type": "address"
		},
		{
			"indexed": true,
			"name": "from",
			"type": "address"
		},
		{
			"indexed": true,
			"name": "to",
			"type": "address"
		},
		{
			"indexed": false,
			"name": "value",
			"type": "uint256"
		},
		{
			"indexed": false,
			"name": "data",
			"type": "bytes"
		},
		{
			"indexed": false,
			"name": "operatorData",
			"type": "bytes"
		}
	],
	"name": "TransferWithData",
	"type": "event"
}`)

var ERC1400TransferWithDataTx = &types.Log{
	Address:     common.HexToAddress("0xE41d2489571d322189246DaFA5ebDe1F4699F498"),
	BlockHash:   common.HexToHash("0xea2460a53299f7201d82483d891b26365ff2f49cd9c5c0c7686fd75599fda5b2"),
	BlockNumber: 6383829,
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000061000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	Index:       13,
	TxIndex:     17,
	TxHash:      common.HexToHash("0x7bec5494eddfba3680fb44053c822ffdc24fb5f6ab7e5e9179b897bfac4cf210"),
	Topics: []common.Hash{
		common.HexToHash("0xe8f0a47da72ca43153c7a5693a827aa8456f52633de9870a736e5605bff4af6d"),
		common.HexToHash("0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236"),
		common.HexToHash("0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236"),
		common.HexToHash("0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f"),
	},
}

func TestDecodeLOGTransferWithData(t *testing.T) {
	event := newEvent(ERC1400TransferWithDataAbi)

	decoded, _ := Decode(event, ERC1400TransferWithDataTx)

	expected := map[string]string{
		"operator":     "0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
		"from":         "0xd71400daD07d70C976D6AAFC241aF1EA183a7236",
		"to":           "0xb5747835141b46f7C472393B31F8F5A57F74A44f",
		"value":        "3",
		"data":         "0x000000000000000000000000000000000000000000000000000000006f0c7f50cd4b7e4466b726279b1506bc89d8e74ab9268a255eeb1c78f163d51a83c7380d54a8b597ee26351c15c83f922fd6b37334970d3f832e5e11e36acbecb460ffdb01",
		"operatorData": "0x",
	}
	eq := reflect.DeepEqual(expected, decoded)
	if !eq {
		t.Errorf("Decode TransferWithDataERC1400ABI: expected mapping %q but got %q", expected, decoded)
	}
}

var ERC1400BridgePendingTransferAbi = []byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"index_origin","type":"uint256"},{"indexed":false,"name":"transfer_id","type":"bytes32"},{"indexed":false,"name":"parameters","type":"bytes32[6]"}],"name":"PendingTransfer","type":"event"}`)

var ERC1400BridgePendingTransferTx = &types.Log{
	Data: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000015a4f2c3ad66af173634e1cc1e389232788ec41756ec2821b9a231f996c4faad00000000000000000000000008f371daa8a5325f53b754a7017ac3803382bc8470000000000000000000000003404370fddb2b0e79f2571e170b112a66f974fb95265736572766564000000000000000000000000000000000000000000000000000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f0000000000000000000000000000000000000000000000000000000000000005"),
	Topics: []common.Hash{
		common.HexToHash("0xd8589d63a2df3a19b774d092cc22aec68d0be6537da4f37a362fbba9f6296845"),
		common.HexToHash("0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236"),
		common.HexToHash("0x000000000000000000000000d71400dad07d70c976d6aafc241af1ea183a7236"),
		common.HexToHash("0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f"),
	},
}

func TestERC1400BridgePendingTransfer(t *testing.T) {
	event := newEvent(ERC1400BridgePendingTransferAbi)

	decoded, _ := Decode(event, ERC1400BridgePendingTransferTx)

	expected := map[string]string{
		"index_origin": "1",
		"transfer_id":  "0x5a4f2c3ad66af173634e1cc1e389232788ec41756ec2821b9a231f996c4faad0",
		"parameters":   "[0x0000000000000000000000008f371daa8a5325f53b754a7017ac3803382bc847,0x0000000000000000000000003404370fddb2b0e79f2571e170b112a66f974fb9,0x5265736572766564000000000000000000000000000000000000000000000000,0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f,0x000000000000000000000000b5747835141b46f7c472393b31f8f5a57f74a44f,0x0000000000000000000000000000000000000000000000000000000000000005]",
	}
	eq := reflect.DeepEqual(expected, decoded)
	if !eq {
		t.Errorf("Decode TransferWithDataERC1400ABI: expected mapping %q but got %q", expected, decoded)
	}
}

var EventTestNonIndexedAbi = []byte(`{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address[2]"}],"name":"EventTEst","type":"event"}`)

var EventTestNonIndexedTx = &types.Log{
	Data: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"),
	Topics: []common.Hash{
		common.HexToHash("0xc69e2952deb9c87385c05617318dac13d6e25ebcf771aa2ae600cc8b1c1c7c73"),
	},
}

func TestEventNonIndexedTx(t *testing.T) {
	event := newEvent(EventTestNonIndexedAbi)

	decoded, _ := Decode(event, EventTestNonIndexedTx)

	expected := map[string]string{
		"from": "[0x0000000000000000000000000000000000000000,0x0000000000000000000000000000000000000001]",
	}
	eq := reflect.DeepEqual(expected, decoded)
	if !eq {
		t.Errorf("Decode TransferWithDataERC1400ABI: expected mapping %q but got %q", expected, decoded)
	}
}
