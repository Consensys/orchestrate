package generic

import (
	"testing"
)

func TestGeneric(t *testing.T) {
	_ = NewTCP()
	_ = NewHTTP()
	_ = NewGRPCServer()
}
