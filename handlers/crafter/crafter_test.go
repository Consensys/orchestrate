package crafter

import (
	"context"
	"fmt"
	"testing"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/args"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/ethereum"
)

const testsNum = 11

type MockCrafter struct{}

var (
	callPayload        = "0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
	constructorPayload = "0xf622a9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000"
)

func (c *MockCrafter) CraftCall(method *ethAbi.Method, methodArgs ...string) ([]byte, error) {
	if len(methodArgs) != 1 {
		return []byte(``), errors.InvalidArgsCountError("could not craft call expected args len to be 1").SetComponent("mock")
	}
	return hexutil.MustDecode(callPayload), nil
}

func (c *MockCrafter) CraftConstructor(bytecBTWode []byte, method *ethAbi.Method, methodArgs ...string) ([]byte, error) {
	if len(methodArgs) != 1 {
		return []byte(``), errors.InvalidArgsCountError("could not craft call expected args len to be 1").SetComponent("mock")
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
		ctx.Set("error.code", errors.InvalidArgsCountError("").GetCode())
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
		ctx.Set("error.code", errors.EncodingError("").GetCode())
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
		ctx.Set("error.code", errors.InvalidSignatureError("").GetCode())
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
		ctx.Set("error.code", errors.NotFoundError("").GetCode())
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
		ctx.Set("error.code", errors.InvalidArgsCountError("").GetCode())
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
	contractRegistry *clientmock.ContractRegistryClient
}

func (s *CrafterTestSuite) SetupSuite() {
	s.contractRegistry = clientmock.New()
	_, err := s.contractRegistry.RegisterContract(context.Background(),
		&contractregistry.RegisterContractRequest{
			Contract: &abi.Contract{
				Id: &abi.ContractId{
					Name: "known",
				},
				Bytecode: hexutil.MustDecode("0x"),
			},
		})
	assert.NoError(s.T(), err)
	s.Handler = Crafter(s.contractRegistry, &MockCrafter{})
}

func (s *CrafterTestSuite) TestCrafter() {
	var txctxs []*engine.TxContext
	for i := 0; i < testsNum; i++ {
		txctxs = append(txctxs, makeCrafterContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors", txctx.Envelope.Args)
		for _, err := range txctx.Envelope.Errors {
			assert.Equal(s.T(), txctx.Get("error.code").(uint64), err.GetCode(), "Error code be correct")
		}
		assert.Equal(s.T(), txctx.Get("result").(string), txctx.Envelope.Tx.TxData.GetData().Hex(), "Expected correct payload", txctx.Envelope.Args)
	}
}

func TestCrafter(t *testing.T) {
	suite.Run(t, new(CrafterTestSuite))
}
