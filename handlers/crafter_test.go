package handlers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
)

var (
	ERC20TransferABI = []byte(`{
		"constant": false,
		"inputs": [
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	}`)
	ERC20TransferMethod abi.Method
	_                   = json.Unmarshal(ERC20TransferABI, &ERC20TransferMethod)
	ERC20Getter         = NewDummyABIGetter(&ERC20TransferMethod)
)

var errGetABI = errors.New("Could not retrieve ABI")

type ErrorABIGetter struct{}

func (g *ErrorABIGetter) GetMethodByID(ID string) (*abi.Method, error) {
	return nil, errGetABI
}

type TestCrafterMsg struct {
	methodID string
	Args     []string
}

func newCrafterTestMessage(i int) *TestCrafterMsg {
	if i%2 == 0 {
		// Valid args
		return &TestCrafterMsg{"abcde", []string{"0xfF778b716FC07D98839f48DdB88D8bE583BEB684", "0x2386f26fc10000"}}
	}
	// invalid args
	return &TestCrafterMsg{"abcde", []string{}}
}

func testCrafterLoader() infra.HandlerFunc {
	return func(ctx *infra.Context) {
		msg := ctx.Msg.(*TestCrafterMsg)
		ctx.Pb.Call = &tracepb.Call{MethodId: msg.methodID, Args: msg.Args}

		// Load Trace from protobuffer
		protobuf.LoadTrace(ctx.Pb, ctx.T)
	}
}

type testCrafterHandler struct {
	mux     *sync.Mutex
	handled []*infra.Context
}

func (h *testCrafterHandler) Handler(maxtime int, t *testing.T) infra.HandlerFunc {
	return func(ctx *infra.Context) {
		// We add some randomness in time execution
		r := rand.Intn(maxtime)
		time.Sleep(time.Duration(r) * time.Millisecond)
		h.mux.Lock()
		defer h.mux.Unlock()
		h.handled = append(h.handled, ctx)
	}
}

func TestCrafter(t *testing.T) {
	// Valid craft
	h := Crafter(ERC20Getter)
	testH := &testCrafterHandler{
		mux:     &sync.Mutex{},
		handled: []*infra.Context{},
	}

	// Create worker
	w := infra.NewWorker(100)
	w.Use(testCrafterLoader())
	w.Use(h)
	w.Use(testH.Handler(50, t))

	// Create a input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed sarama channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newCrafterTestMessage(i)
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	// We expected half of rounds to have aborted
	if len(testH.handled) != rounds/2 {
		t.Errorf("Crafter: expected %v rounds but got %v", rounds, len(testH.handled))
	}

	expected := "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	for _, ctx := range testH.handled {
		if len(ctx.T.Errors) != 0 && hexutil.Encode(ctx.T.Tx().Data()) != "0x" {
			t.Errorf("Crafter: expected no raw tx on error but got %q", hexutil.Encode(ctx.T.Tx().Data()))
		}
		if len(ctx.T.Errors) == 0 && hexutil.Encode(ctx.T.Tx().Data()) != expected {
			t.Errorf("Crafter: expected raw %q but got %q", expected, hexutil.Encode(ctx.T.Tx().Data()))
		}
	}
}
