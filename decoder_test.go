package ethereum

import (
	"encoding/json"
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
	Address:     common.HexToAddress("0xE41d2489571d322189246DaFA5ebDe1F4699F498"),
	BlockHash:   common.HexToHash("0xea2460a53299f7201d82483d891b26365ff2f49cd9c5c0c7686fd75599fda5b2"),
	BlockNumber: 6383829,
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000001a055690d9db80000"),
	Index:       13,
	TxIndex:     17,
	TxHash:      common.HexToHash("0x7bec5494eddfba3680fb44053c822ffdc24fb5f6ab7e5e9179b897bfac4cf210"),
	Topics: []common.Hash{
		common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		common.HexToHash("0x000000000000000000000000ba826fec90cefdf6706858e5fbafcb27a290fbe0"),
		common.HexToHash("0x0000000000000000000000004aee792a88edda29932254099b9d1e06d537883f"),
	},
}

var LOGFILLABI = []byte(`{
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

var testLogLOGFILLABI = &types.Log{
	Address:     common.HexToAddress("0xE41d2489571d322189246DaFA5ebDe1F4699F498"),
	BlockHash:   common.HexToHash("0xea2460a53299f7201d82483d891b26365ff2f49cd9c5c0c7686fd75599fda5b2"),
	BlockNumber: 6383829,
	Data:        hexutil.MustDecode("0x000000000000000000000000e269e891a2ec8585a378882ffa531141205e92e9000000000000000000000000d7732e3783b0047aa251928960063f863ad022d8000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc20000000000000000000000000000000000000000000032d26d12e980b6000000000000000000000000000000000000000000000000000000602d4ec2c348a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000306a9a7ecbd9446559a2c650b4cfc16d1fb615aa2b3f4f63078da6d021268440"),
	Index:       13,
	TxIndex:     17,
	TxHash:      common.HexToHash("0x7bec5494eddfba3680fb44053c822ffdc24fb5f6ab7e5e9179b897bfac4cf210"),
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

	testEventDecoder := &EventDecoder{
		Inputs: event.Inputs,
	}

	decoded, _ := testEventDecoder.Decode(testLogERC20ABI)

	m := map[string]string{
		"tokens": "30000000000000000000",
		"from":   "0xBA826fEc90CEFdf6706858E5FbaFcb27A290Fbe0",
		"to":     "0x4aEE792A88eDDA29932254099b9d1e06D537883f",
	}
	eq := reflect.DeepEqual(m, decoded)
	if !eq {
		t.Errorf("Decode: expected mapping %q but got %q", m, decoded)
	}
}

func TestDecodeLOGFILLABI(t *testing.T) {
	event := newEvent(LOGFILLABI)

	testEventDecoder := &EventDecoder{
		Inputs: event.Inputs,
	}

	decoded, _ := testEventDecoder.Decode(testLogLOGFILLABI)

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
		t.Errorf("Decode: expected mapping %q but got %q", m, decoded)
	}
}

// func TestFormatIndexedEvent(t *testing.T) {

// }

// func TestFormatNonIndexEvent(t *testing.T) {

// }
