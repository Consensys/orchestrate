package examples

// Msg is a dummy engine.Msg implementation
type Msg string

// Entrypoint is a dummy implementation of the method "Entrypoint of the dummy engine"
func (s Msg) Entrypoint() string { return "" }
