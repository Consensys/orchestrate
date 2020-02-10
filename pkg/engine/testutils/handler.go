package testutils

import (
	"sync"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

// HandlerTestSuite is an utility suite to test an handler
type HandlerTestSuite struct {
	suite.Suite

	Handler engine.HandlerFunc
}

// Handle execute handler concurrently on every input Builder
func (s *HandlerTestSuite) Handle(txctxs []*engine.TxContext) {
	// Execute handler on every Builder concurrently
	wg := &sync.WaitGroup{}
	for _, txctx := range txctxs {
		wg.Add(1)
		go func(txctx *engine.TxContext) {
			s.Handler(txctx)
			wg.Done()
		}(txctx)
	}
	// Wait for all Builder to have complete execution
	wg.Wait()
}
