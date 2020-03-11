package crafter

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	typesabi "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"
)

type MockCrafter struct {
	crafter *abi.BaseCrafter
}

func (c *MockCrafter) CraftCall(sig string, args ...string) ([]byte, error) {
	if len(args) != 1 {
		return []byte(``), fmt.Errorf("could not craft call expected args len to be 1")
	}
	return c.crafter.CraftCall(sig, args...)
}

func (c *MockCrafter) CraftConstructor(bytecode []byte, sig string, args ...string) ([]byte, error) {
	if len(args) != 0 {
		return []byte(``), fmt.Errorf("could not craft call expected args len to be 1")
	}
	return c.crafter.CraftConstructor(bytecode, sig)
}

func makeCrafterContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope = tx.NewEnvelope()

	switch i {
	case 0:
		ctx.Set("errors", 0)
		ctx.Set("result", "")
	case 1:
		_ = ctx.Envelope.
			SetData(hexutil.MustDecode("0xa9059cbb"))
		ctx.Set("errors", 0)
		ctx.Set("result", "0xa9059cbb")
	case 2:
		_ = ctx.Envelope.
			SetMethodSignature("increment(uint256)").
			SetArgs([]string{"0xab"})
		ctx.Set("errors", 0)
		ctx.Set("result", "0x7cf5dab000000000000000000000000000000000000000000000000000000000000000ab")
	case 3:
		_ = ctx.Envelope.
			SetMethodSignature("testMethod()")
		ctx.Set("errors", 1)
		ctx.Set("result", "")
	case 4:
		_ = ctx.Envelope.
			SetContractName("known").
			SetMethodSignature("constructor()")
		ctx.Set("errors", 0)
		ctx.Set("result", "0x6080604052348015600f57600080fd5b5061010a8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80637cf5dab014602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b8060008082825401925050819055507f38ac789ed44572701765277c4d0970f2db1c1a571ed39e84358095ae4eaa54203382604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15056fea265627a7a72315820c084d653e3ba7607a5b05fb98edf3373a2b542aa6abdd9ae89cd4a407bb0a2b464736f6c63430005100032")
	case 5:
		// Invalid method signature
		_ = ctx.Envelope.
			SetContractName("known").
			SetMethodSignature("constructor(")
		ctx.Set("errors", 1)
		ctx.Set("result", "")
	case 6:
		// Invalid contract name
		_ = ctx.Envelope.
			SetContractName("unknown").
			SetMethodSignature("constructor()")
		ctx.Set("errors", 1)
		ctx.Set("result", "")
	case 7:
		// Invalid number of arguments for a constructor
		_ = ctx.Envelope.
			SetContractName("known").
			SetMethodSignature("constructor()").
			SetArgs([]string{"0xabcd"})
		ctx.Set("errors", 1)
		ctx.Set("result", "")
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
			Contract: &typesabi.Contract{
				Id: &typesabi.ContractId{
					Name: "known",
				},
				Abi:              `[{"constant":false,"inputs":[{"internalType":"uint256","name":"value","type":"uint256"}],"name":"increment","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"from","type":"address"},{"indexed":false,"internalType":"uint256","name":"by","type":"uint256"}],"name":"Incremented","type":"event"}]`,
				Bytecode:         "0x6080604052348015600f57600080fd5b5061010a8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80637cf5dab014602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b8060008082825401925050819055507f38ac789ed44572701765277c4d0970f2db1c1a571ed39e84358095ae4eaa54203382604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15056fea265627a7a72315820c084d653e3ba7607a5b05fb98edf3373a2b542aa6abdd9ae89cd4a407bb0a2b464736f6c63430005100032",
				DeployedBytecode: "0x6080604052348015600f57600080fd5b506004361060285760003560e01c80637cf5dab014602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b8060008082825401925050819055507f38ac789ed44572701765277c4d0970f2db1c1a571ed39e84358095ae4eaa54203382604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15056fea265627a7a72315820cc73213703da717157e9aa146473f2af6823d442bfb100062b58833aae34fa7b64736f6c63430005100032",
			},
		})
	assert.NoError(s.T(), err)
	s.Handler = Crafter(s.contractRegistry, &MockCrafter{})
}

func (s *CrafterTestSuite) TestCrafter() {
	var txctxs []*engine.TxContext
	for i := 0; i < 8; i++ {
		txctxs = append(txctxs, makeCrafterContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for i, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "%d/%d - Expected right count of errors", i, len(txctxs), txctx.Envelope.Args)
		assert.Equal(s.T(), txctx.Get("result").(string), txctx.Envelope.GetData(), "%d/%d - Expected correct payload", i, len(txctxs), txctx.Envelope.Args)
	}
}

func TestCrafter(t *testing.T) {
	suite.Run(t, new(CrafterTestSuite))
}
