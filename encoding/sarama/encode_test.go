package sarama

import (
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

var PostState, _ = hexutil.Decode("0xabcdef")
var Bloom, _ = hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80")

var testEnvelope = &envelope.Envelope{
	Chain:    &common.Chain{Id: "0x1", IsEIP155: true},
	Sender:   &common.Account{Id: "abcde", Addr: "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"},
	Receiver: &common.Account{Id: "toto", Addr: "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"},
	Tx: &ethereum.Transaction{
		TxData: &ethereum.TxData{
			Nonce:    10,
			To:       "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff",
			Value:    "0xa2bfe3",
			Data:     "0xbabe",
			GasPrice: "0x0",
			Gas:      1234,
		},
		Raw:  "0xbeef",
		Hash: "0x24f5acae441335ad59220734d1ffd9cc1f6f525d39f2785859298048c25fb814",
	},
	Call: &common.Call{
		Contract: &abi.Contract{
			Name: "ERC20",
		},
		Method: &abi.Method{
			Name: "transfer",
		},
		Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"},
	},
	Receipt: &ethereum.Receipt{
		Logs:              []*ethereum.Log{},
		ContractAddress:   "0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C",
		PostState:         PostState,
		Status:            1,
		TxHash:            "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
		Bloom:             Bloom,
		GasUsed:           13456,
		CumulativeGasUsed: 19304777,
		BlockHash:         "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
		BlockNumber:       1234,
		TxIndex:           1,
	},
	Errors: []*common.Error{
		&common.Error{Type: 0, Message: "Error 0"},
		&common.Error{Type: 1, Message: "Error 1"},
	},
	Metadata: &envelope.Metadata{
		Id: "test",
	},
}

var expected, _ = proto.Marshal(testEnvelope)

func newEnvelope() *envelope.Envelope {
	// Create Envelope
	e := &envelope.Envelope{}
	proto.Unmarshal(expected, e)

	return e
}

func TestMarshaller(t *testing.T) {
	m := NewMarshaller()
	messages := make([]*sarama.ProducerMessage, 0)
	rounds := 1
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		message := &sarama.ProducerMessage{}
		messages = append(messages, message)
		wg.Add(1)
		go func(msg *sarama.ProducerMessage) {
			defer wg.Done()
			m.Marshal(newEnvelope(), msg)
		}(message)
	}
	wg.Wait()

	for _, msg := range messages {
		b, err := msg.Value.Encode()
		if err != nil {
			t.Errorf("SaramaMarshaller: expected valid value")
		}
		if string(b) != string(expected) {
			t.Errorf("SaramaMarshaller: expected %q but got %q", string(expected), string(b))
		}
	}
}

