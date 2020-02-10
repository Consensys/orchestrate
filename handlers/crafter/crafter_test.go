package crafter

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client/mocks"
)

const testsNum = 8

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
	ctx.Builder = tx.NewBuilder()

	switch i {
	case 0:
		ctx.Set("errors", 0)
		ctx.Set("result", "")
	case 1:
		_ = ctx.Builder.SetData(hexutil.MustDecode("0xa9059cbb"))
		ctx.Set("errors", 0)
		ctx.Set("result", "0xa9059cbb")
	case 2:
		_ = ctx.Builder.SetMethodSignature("known()").SetArgs([]string{"test"})
		ctx.Set("errors", 0)
		ctx.Set("result", callPayload)
	case 3:
		_ = ctx.Builder.SetMethodSignature("known()")
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.InvalidArgsCountError("").GetCode())
		ctx.Set("result", "")
	case 4:
		_ = ctx.Builder.SetContractName("known").SetMethodSignature("constructor()").SetArgs([]string{"test"})
		ctx.Set("errors", 0)
		ctx.Set("result", constructorPayload)
	case 5:
		// Invalid method signature
		_ = ctx.Builder.SetContractName("known").SetMethodSignature("constructor)").SetArgs([]string{"0xabcd"})
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.InvalidSignatureError("").GetCode())
		ctx.Set("result", "")
	case 6:
		// Invalid contract name
		_ = ctx.Builder.SetContractName("unknown").SetMethodSignature("constructor()").SetArgs([]string{"0xabcd"})
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.NotFoundError("").GetCode())
		ctx.Set("result", "")
	case 7:
		// Invalid number of arguments for a constructor
		_ = ctx.Builder.SetContractName("known").SetMethodSignature("constructor()").SetArgs([]string{"0xabcd", "123"})
		ctx.Set("errors", 1)
		ctx.Set("error.code", errors.InvalidArgsCountError("").GetCode())
		ctx.Set("result", "")
	case 8:
		_ = ctx.Builder.SetContractName("known").SetMethodSignature("constructor()").SetArgs([]string{"0xabcd"})
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
				Abi:              `[]`,
				Bytecode:         hexutil.Encode([]byte{1, 2, 3}),
				DeployedBytecode: hexutil.Encode([]byte{1, 2}),
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

	for i, txctx := range txctxs {
		assert.Len(s.T(), txctx.Builder.Errors, txctx.Get("errors").(int), "%d/%d - Expected right count of errors", i, len(txctxs), txctx.Builder.Args)
		for _, err := range txctx.Builder.Errors {
			assert.Equal(s.T(), txctx.Get("error.code").(uint64), err.GetCode(), "%d/%d - Error code be correct", i, len(txctxs))
		}
		assert.Equal(s.T(), txctx.Get("result").(string), txctx.Builder.GetData(), "%d/%d - Expected correct payload", i, len(txctxs), txctx.Builder.Args)
	}
}

func TestCrafter(t *testing.T) {
	suite.Run(t, new(CrafterTestSuite))
}
