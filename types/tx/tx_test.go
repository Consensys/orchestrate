package tx

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnvelope(t *testing.T) {
	envelope := &TxRequest{
		Id:     uuid.NewV4().String(),
		Method: Method_ETH_SENDRAWTRANSACTION,
		Chain:  "testChain",
		Params: &Params{
			From:            "0x7e654d251da770a068413677967f6d3ea2fea9e4",
			To:              "0xdbb881a51cd4023e4400cef3ef73046743f08da3",
			Gas:             "10",
			GasPrice:        "1089",
			Value:           "56757",
			Nonce:           "10",
			Data:            "0xab",
			Contract:        "test",
			MethodSignature: "constructor()",
		},
	}
	_, err := envelope.Builder()

	assert.NoError(t, err)
}
