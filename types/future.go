package types

// Future is an element used to start a task and retrieve its result later
type Future struct {
	res chan interface{}
	err chan error
}

// NewFuture creates a new future
func NewFuture() *Future {
	return &Future{
		res: make(chan interface{}),
		err: make(chan error),
	}
}

// Close future
func (f *Future) Close() {
	close(f.res)
	close(f.err)
}

// Err return a Error channel
func (f *Future) Err() chan error {
	return f.err
}

// Result return result channel
func (f *Future) Result() chan interface{} {
	return f.res
}
