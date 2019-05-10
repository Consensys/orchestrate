package crafter

import (
	"fmt"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

type MockABIRegistry struct{}

var unkwnon = "unknown@"

func (r *MockABIRegistry) GetMethodByID(id string) (*ethabi.Method, error) {
	if id == unkwnon {
		return &ethabi.Method{}, fmt.Errorf("could not retrieve ABI")
	}
	return &ethabi.Method{}, nil
}

func (r *MockABIRegistry) GetBytecodeByID(id string) ([]byte, error) {
	if id == unkwnon {
		return []byte{}, fmt.Errorf("could not retrieve bytecode")
	}
	return []byte{246, 34}, nil
}

func (r *MockABIRegistry) GetMethodBySig(sig string) (*ethabi.Method, error) {
	return &ethabi.Method{}, nil
}

func (r *MockABIRegistry) GetEventByID(id string) (*ethabi.Event, error) {
	if id == "unknown" {
		return &ethabi.Event{}, fmt.Errorf("could not retrieve ABI")
	}
	return &ethabi.Event{}, nil
}

func (r *MockABIRegistry) GetEventBySig(sig string) (*ethabi.Event, error) {
	return &ethabi.Event{}, nil
}

func (r *MockABIRegistry) RegisterContract(contract *abi.Contract) error {
	return nil
}

type MockCrafter struct{}

var (
	callPayload        = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	constructorPayload = "0xf622a9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
)

func (c *MockCrafter) CraftCall(method ethabi.Method, args ...string) ([]byte, error) {
	if len(args) != 1 {
		return []byte(``), fmt.Errorf("could not craft expected args len to be 1")
	}
	return hexutil.MustDecode(callPayload), nil
}

func (c *MockCrafter) CraftConstructor(bytecode []byte, method ethabi.Method, args ...string) ([]byte, error) {
	if len(args) != 1 {
		return []byte(``), fmt.Errorf("could not craft expected args len to be 1")
	}
	return hexutil.MustDecode(constructorPayload), nil
}

func makeCrafterContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}

	switch i % 6 {
	case 0:
		ctx.Set("errors", 0)
		ctx.Set("result", "")
	case 1:
		ctx.Envelope.Tx.TxData = (&ethereum.TxData{}).SetData(hexutil.MustDecode("0xa9059cbb"))
		ctx.Set("errors", 0)
		ctx.Set("result", "0xa9059cbb")
	case 2:
		ctx.Envelope.Call = &common.Call{
			Method: &abi.Method{Name: "unknown"},
		}
		ctx.Envelope.Call.Args = []string{"test"}
		ctx.Set("errors", 1)
		ctx.Set("result", "")
	case 3:
		ctx.Envelope.Call = &common.Call{
			Method: &abi.Method{Name: "known"},
		}
		ctx.Envelope.Call.Args = []string{"test"}
		ctx.Set("errors", 0)
		ctx.Set("result", callPayload)
	case 4:
		ctx.Envelope.Call = &common.Call{
			Method: &abi.Method{Name: "known"},
		}
		ctx.Set("errors", 1)
		ctx.Set("result", "")
	case 5:
		ctx.Envelope.Call = &common.Call{
			Contract: &abi.Contract{Name: "known"},
			Method:   &abi.Method{Name: "constructor"},
			Args:     []string{"0xabcd"},
		}
		ctx.Set("errors", 0)
		ctx.Set("result", constructorPayload)
	}
	return ctx
}

type CrafterTestSuite struct {
	testutils.HandlerTestSuite
}

func (s *CrafterTestSuite) SetupSuite() {
	s.Handler = Crafter(&MockABIRegistry{}, &MockCrafter{})
}

func (s *CrafterTestSuite) TestCrafter() {
	rounds := 100
	txctxs := []*engine.TxContext{}
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeCrafterContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
		assert.Equal(s.T(), txctx.Get("result").(string), txctx.Envelope.Tx.TxData.GetData(), "Expected correct payload")
	}
}

func TestCrafter(t *testing.T) {
	suite.Run(t, new(CrafterTestSuite))
}
