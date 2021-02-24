// +build unit

package sarama

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ConsenSys/orchestrate/pkg/types/tx"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	err "github.com/ConsenSys/orchestrate/pkg/types/error"
	"github.com/ConsenSys/orchestrate/pkg/types/ethereum"
)

var PostState = "0xabcdef"
var Bloom = "0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"

var envlp = tx.NewEnvelope().
	SetID("0cdac6e9-8836-4280-8d6b-2e01cba7a1ca").
	MustSetFromString("0xdbb881a51CD4023E4400CEF3ef73046743f08da3").
	SetNonce(10).
	MustSetToString("0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff").
	SetValue(big.NewInt(1000)).
	SetGasPrice(big.NewInt(20)).
	SetGas(1234).
	MustSetRawString("0xbeef").
	MustSetTxHashString("0x24f5acae441335ad59220734d1ffd9cc1f6f525d39f2785859298048c25fb814").
	SetContractName("ERC20").
	SetMethodSignature("transfer(address,uint256)").
	SetArgs([]string{
		"0xfF778b716FC07D98839f48DdB88D8bE583BEB684",
		"0x2386f26fc10000",
	}).
	SetReceipt(&ethereum.Receipt{
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
	}).
	AppendErrors([]*err.Error{
		{Code: 0, Message: "Error 0"},
		{Code: 1, Message: "Error 1"},
	})

var expected, _ = proto.Marshal(envlp.TxEnvelopeAsRequest())

func newEnvelope() *tx.TxEnvelope {
	// Create Envelope
	e := &tx.TxEnvelope{}
	_ = proto.Unmarshal(expected, e)

	return e
}

func TestMarshaller(t *testing.T) {
	messages := make([]*sarama.ProducerMessage, 0)
	rounds := 1
	wg := &sync.WaitGroup{}
	for i := 0; i < rounds; i++ {
		message := &sarama.ProducerMessage{}
		messages = append(messages, message)
		wg.Add(1)
		go func(msg *sarama.ProducerMessage) {
			defer wg.Done()
			_ = Marshal(newEnvelope(), msg)
		}(message)
	}
	wg.Wait()

	for _, msg := range messages {
		b, e := msg.Value.Encode()
		if e != nil {
			t.Errorf("SaramaMarshaller: expected valid value")
		}
		if string(b) != string(expected) {
			t.Errorf("SaramaMarshaller: expected %q but got %q", string(expected), string(b))
		}
	}
}
