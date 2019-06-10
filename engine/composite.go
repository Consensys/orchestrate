
package engine

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
)

// MatcherFunc is an abstract type for a class of function 
// used to determine if a handler is appropriate to use in a Composite
type MatcherFunc func(txctx *TxContext) bool

// MatchAll is a common MatcherFunc that always returns true
func MatchAll(txctx *TxContext) bool { return true }

// MatchableHandler is a sugar type to manipulate a handler and its matcher in the same object
type MatchableHandler struct {
	Handler		HandlerFunc
	Matcher 	MatcherFunc
}

// NewMatchableHandler creates an explicit matchable handler from a given handler and matcher
func NewMatchableHandler(handler HandlerFunc, matcher MatcherFunc) MatchableHandler {
	return MatchableHandler{
		Handler:	handler,
		Matcher:	matcher,
	}
}

// CompositeHandler is an object listing multiple handlers and a matching function
// The executed handler is the first one to be matched. Tests starts at rank 0.
type CompositeHandler struct {
	// chain of handlers to be to be tested for execution
	matchableHandlers 	[]MatchableHandler
}

// NewCompositeHandler creates a new composite handler. Accept zero argument calls
func NewCompositeHandler(matchableHandlers ...MatchableHandler) CompositeHandler {
	return CompositeHandler{ matchableHandlers: matchableHandlers }
}

// Registers append a new handler and matcher to the composite
func (c *CompositeHandler) Registers(handler HandlerFunc, matcher MatcherFunc) {
	h := NewMatchableHandler(handler, matcher) 
	c.matchableHandlers = append(
		c.matchableHandlers,
		h,
	)
}

// BuildFirstMatchFunc returns a handler executing the first handler to match the context 
func (c *CompositeHandler) BuildFirstMatchFunc() HandlerFunc {
	return func(txctx *TxContext) {
		for _, matchable := range c.matchableHandlers {
			if matchable.Matcher(txctx) {
				matchable.Handler(txctx)
				return
			}
		}
	}
}

// BuildAllMatchSeqFunc returns a handler executing all the matches in a sequential way
func (c *CompositeHandler) BuildAllMatchSeqFunc() HandlerFunc {
	return func(txctx *TxContext) {
		for _, matchable := range c.matchableHandlers {
			if matchable.Matcher(txctx) {
				matchable.Handler(txctx)
			}
		}
	}
}

// BuildAllMatchParallelFunc returns a handler executing all the matches in a sequential way
func (c *CompositeHandler) BuildAllMatchParallelFunc() HandlerFunc {
	return func(txctx *TxContext) {
		matched := []func(){}

		for _, matchable := range c.matchableHandlers {
			if matchable.Matcher(txctx) {
				matched = append(matched, func() { matchable.Handler(txctx)})
			}
		}
		common.InParallel(matched...)
	}
}

