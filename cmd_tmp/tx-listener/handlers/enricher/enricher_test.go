package enricher

import (
	"fmt"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	crc "gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/mock"
)

var testsNum = 2

type EnricherTestSuite struct {
	testutils.HandlerTestSuite
	contractRegistry *crc.ContractRegistryClient
}

func (s *EnricherTestSuite) SetupSuite() {
	blocks := make(map[string][]*ethTypes.Block)
	mec := mock.NewClient(blocks)

	crc.New()

	s.contractRegistry = crc.New()
	s.Handler = Enricher(s.contractRegistry, mec)
}

func makeEnricherContext(i int) *engine.TxContext {
	ctx := engine.NewTxContext()
	ctx.Reset()
	ctx.Logger = log.NewEntry(log.StandardLogger())
	ctx.Envelope.Receipt = &ethereum.Receipt{}

	switch i % 2 {
	case 0:
		ctx.Set("errors", 0)
	case 1:
		ctx.Set("errors", 0)
		ctx.Envelope.Receipt.ContractAddress = ethereum.HexToAccount("0xd71400daD07d70C976D6AAFC241aF1EA183a7236")
	default:
		panic(fmt.Sprintf("No test case with number %d", i))
	}
	return ctx
}

func (s *EnricherTestSuite) TestEnricher() {
	var txctxs []*engine.TxContext
	for i := 0; i < testsNum-1; i++ {
		txctxs = append(txctxs, makeEnricherContext(i))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, txctx := range txctxs {
		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors", txctx.Envelope.Args)
		for _, err := range txctx.Envelope.Errors {
			assert.Equal(s.T(), txctx.Get("error.code").(uint64), err.GetCode(), "Error code be correct")
		}
	}
}

func TestEnricher(t *testing.T) {
	suite.Run(t, new(EnricherTestSuite))
}
