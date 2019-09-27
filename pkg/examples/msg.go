package examples

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Msg is a dummy engine.Msg implementation
type Msg string

// Entrypoint is a dummy implementation of the method "Entrypoint of the dummy engine"
func (s Msg) Entrypoint() string 		{ return string(s) }
// Header is a dummy implementation of the method "Header"
func (s Msg) Header() engine.Header     { return nil }
// Key is a dummy implementation of the method "Key"
func (s Msg) Key() []byte        		{ return nil }
// Value is a dummy implementation of the method "Value"
func (s Msg) Value() []byte      		{ return nil }
