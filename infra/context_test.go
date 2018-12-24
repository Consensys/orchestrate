package infra

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	ethpb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/ethereum"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

var testKey = "test"

func newHandler(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), s)
	}
}

var errTest = errors.New("Test Error")

func newErrorHandler(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), s)
		ctx.Error(errTest)
	}
}

func newAborter(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), s)
		ctx.AbortWithError(errTest)
	}
}

func newMiddleware(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		sA := fmt.Sprintf("%v-before", s)
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), sA)

		ctx.Next()

		sB := fmt.Sprintf("%v-after", s)
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), sB)
	}
}

func TestNext(t *testing.T) {
	ctx := NewContext()

	var (
		hA   = newHandler("hA", t)
		mA   = newMiddleware("mA", t)
		hErr = newErrorHandler("err", t)
		mB   = newMiddleware("mB", t)
		hB   = newHandler("hB", t)
		a    = newAborter("abort", t)
		hC   = newHandler("hC", t)
		mC   = newMiddleware("mC", t)
	)
	// Initialize context
	ctx.Init([]HandlerFunc{hA, mA, hErr, mB, hB, a, hC, mC})
	ctx.Keys[testKey] = []string{}

	// Handle context
	ctx.Next()

	res := ctx.Keys[testKey].([]string)
	expected := []string{"hA", "mA-before", "err", "mB-before", "hB", "abort", "mB-after", "mA-after"}
	if len(res) != len(expected) {
		t.Errorf("Context: expected %v execution but got %v", len(expected), len(res))
	}

	for i, s := range expected {
		if s != res[i] {
			t.Errorf("Context: expected %q at %v but got %q", s, i, res[i])
		}
	}

	if len(ctx.T.Errors) != 2 {
		t.Errorf("Context: expected 2 errors but got %v", len(ctx.T.Errors))
	}
}

func TestLoadMessage(t *testing.T) {
	testMsg := &sarama.ConsumerMessage{}
	testMsg.Value, _ = proto.Marshal(
		&tracepb.Trace{
			Sender:   &tracepb.Account{Id: "abc", Address: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
			Chain:    &tracepb.Chain{Id: "abc", IsEIP155: true},
			Receiver: &tracepb.Account{Id: "abc", Address: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
			Call:     &tracepb.Call{MethodId: "abc", Args: []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x1"}},
			Transaction: &ethpb.Transaction{
				TxData: &ethpb.TxData{Nonce: 1, To: "0xfF778b716FC07D98839f48DdB88D8bE583BEB684", Value: "0x2386f26fc10000", Gas: 21136, GasPrice: "0xee6b2800", Data: "0xabcd"},
				Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
				Hash:   "0xbf0b3048242aff8287d1dd9de0d2d100cee25d4ea45b8afa28bdfc1e2a775afd",
			},
			Errors: []*tracepb.Error{&tracepb.Error{Type: 0, Message: "Error 0"}, &tracepb.Error{Type: 1, Message: "Error 1"}},
		},
	)

	ctx := NewContext()
	ctx.loadMessage(testMsg)

	protobuf.DumpTrace(ctx.T, ctx.pb)
	if len(ctx.T.Errors) != 2 {
		t.Errorf("Context: expected 2 errors but got %v", len(ctx.T.Errors))
	}

	if ctx.T.Tx().Nonce() != 1 {
		t.Errorf("Context: expected Nonce to be 1 but got %v", ctx.T.Tx().Nonce())
	}

	ctx.Reset()
	ctx.loadMessage(&sarama.ConsumerMessage{Value: []byte(`>>Error<<`)})

	if len(ctx.T.Errors) != 1 {
		t.Errorf("Context: expected 1 error but got %v", len(ctx.T.Errors))
	}

	ctx.Reset()
	ctx.loadMessage(&sarama.ConsumerMessage{Value: []byte(`>>Error<<`)})

	if len(ctx.T.Errors) != 1 {
		t.Errorf("Context: expected 1 error but got %v", len(ctx.T.Errors))
	}
}
