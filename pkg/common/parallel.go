package common

import (
	"sync"
)

// InParallel runs provided functions in parallel and waits until
// all of them are finished
func InParallel(funcs ...func()) {
	wg := sync.WaitGroup{}
	wg.Add(len(funcs))

	for _, f := range funcs {
		go func(funcToCall func()) {
			funcToCall()
			wg.Done()
		}(f)
	}

	wg.Wait()
}
