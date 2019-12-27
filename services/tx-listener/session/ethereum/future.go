package ethereum

// Future is an element used to start a task and retrieve its result later
type Future struct {
	res chan interface{}
	err chan error
}

// NewFuture creates a new future
func NewFuture(fn func() (interface{}, error)) *Future {
	future := &Future{
		res: make(chan interface{}, 1),
		err: make(chan error, 1),
	}

	go func() {
		res, err := fn()
		if err != nil {
			future.err <- err
		} else {
			future.res <- res
		}
	}()

	return future
}

// Err return a Error channel
func (f *Future) Err() <-chan error {
	return f.err
}

// Result return result channel
func (f *Future) Result() <-chan interface{} {
	return f.res
}

// Close future
func (f *Future) Close() {
	close(f.res)
	close(f.err)
}
