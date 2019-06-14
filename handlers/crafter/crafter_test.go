package crafter

import (
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

type MockABIRegistry struct{}

const testsNum = 11

func (r *MockABIRegistry) RegisterContract(contract *abi.Contract) error {
	return nil
}

func (r *MockABIRegistry) GetContractABI(contract *abi.Contract) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockABIRegistry) GetContractBytecode(contract *abi.Contract) ([]byte, error) {
	if contract.GetId().Name == "unknown" {
		return []byte{}, fmt.Errorf("could not retrieve bytecode")
	}
	return []byte{246, 34}, nil
}

func (r *MockABIRegistry) GetContractDeployedBytecode(contract *abi.Contract) ([]byte, error) {
	return []byte{}, nil
}

func (r *MockABIRegistry) GetMethodsBySelector(selector [4]byte, contract common.AccountInstance) (method *ethAbi.Method, methods []*ethAbi.Method, e error) {
	return &ethAbi.Method{}, make([]*ethAbi.Method, 0), nil
}

func (r *MockABIRegistry) GetEventsBySigHash(selector ethCommon.Hash, contract common.AccountInstance, indexedInputCount uint) (event *ethAbi.Event, events []*ethAbi.Event, e error) {
	return &ethAbi.Event{}, make([]*ethAbi.Event, 0), nil
}

func (r *MockABIRegistry) RequestAddressUpdate(contract common.AccountInstance) error {
	return nil
}

type MockCrafter struct{}

var (
	callPayload        = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	constructorPayload = "0xf622a9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
)

func (c *MockCrafter) CraftCall(method ethAbi.Method, methodArgs ...string) ([]byte, error) {
	if len(methodArgs) != 1 {
		return []byte(``), fmt.Errorf("could not craft call expected args len to be 1")
	}
	return hexutil.MustDecode(callPayload), nil
}

func (c *MockCrafter) CraftConstructor(bytecBTWode []byte, method ethAbi.Method, methodArgs ...string) ([]byte, error) {
	if len(methodArgs) != 1 {
		return []byte(``), fmt.Errorf("could not craft constructor expected args len to be 1")
	}
	return hexutil.MustDecode(constructorPayload), nil
}

func makeCrafterContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope.Tx = &ethereum.Transaction{TxData: &ethereum.TxData{}}

	switch i {
	case 0:
		ctx.Set("errors", 0)
		ctx.Set("result", "0x")
	case 1:
		ctx.Envelope.Tx.TxData = (&ethereum.TxData{}).SetData(hexutil.MustDecode("0xa9059cbb"))
		ctx.Set("errors", 0)
		ctx.Set("result", "0xa9059cbb")
	case 2:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Method: &abi.Method{Signature: "known()"},
				Args:   []string{"test"},
			},
		}
		ctx.Set("errors", 0)
		ctx.Set("result", callPayload)
	case 3:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Method: &abi.Method{Signature: "known()"},
			},
		}
		ctx.Set("errors", 1)
		ctx.Set("result", "0x")
	case 4:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "known",
					},
				},
				Method: &abi.Method{Signature: "constructor()"},
				Args:   []string{"0xabcd"},
			},
		}
		ctx.Set("errors", 0)
		ctx.Set("result", constructorPayload)
	case 5:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "known",
					},
				},
				// Invalid ABI
				Method: &abi.Method{
					Signature: "constructor()",
					Abi:       []byte{1, 2, 3},
				},
				Args: []string{"0xabcd"},
			},
		}
		ctx.Set("errors", 1)
		ctx.Set("result", "0x")
	case 6:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "known",
					},
				},
				// Invalid method signature
				Method: &abi.Method{Signature: "constructor)"},
				Args:   []string{"0xabcd"},
			},
		}
		ctx.Set("errors", 1)
		ctx.Set("result", "0x")
	case 7:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "unknown",
					},
				},
				Method: &abi.Method{Signature: "constructor()"},
				// Invalid contract name
				Args: []string{"0xabcd"},
			},
		}
		ctx.Set("errors", 1)
		ctx.Set("result", "0x")
	case 8:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "known",
					},
				},
				Method: &abi.Method{Signature: "constructor()"},
				// Invalid number of arguments for a constructor
				Args: []string{"0xabcd", "123"},
			},
		}
		ctx.Set("errors", 1)
		ctx.Set("result", "0x")
	case 9:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "known",
					},
				},
				Method: &abi.Method{Signature: "constructor()"},
				Args:   []string{"0xabcd"},
			},
		}
		ctx.Envelope.Tx = nil
		ctx.Set("errors", 0)
		ctx.Set("result", constructorPayload)
	case 10:
		ctx.Envelope.Args = &envelope.Args{

			Call: &args.Call{
				Contract: &abi.Contract{
					Id: &abi.ContractId{
						Name: "known",
					},
				},
				Method: &abi.Method{Signature: "constructor()"},
				Args:   []string{"0xabcd"},
			},
		}
		ctx.Envelope.Tx = &ethereum.Transaction{TxData: nil}
		ctx.Set("errors", 0)
		ctx.Set("result", constructorPayload)

	default:
		panic(fmt.Sprintf("No test case with number %d", i))
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
	txctxs := []*engine.TxContext{}
	for i := 0; i < testsNum; i++ {
		txctxs = append(txctxs, makeCrafterContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors", txctx.Envelope.Args)
		assert.Equal(s.T(), txctx.Get("result").(string), txctx.Envelope.Tx.TxData.GetData().Hex(), "Expected correct payload", txctx.Envelope.Args)
	}
}

func TestCrafter(t *testing.T) {
	suite.Run(t, new(CrafterTestSuite))
}
