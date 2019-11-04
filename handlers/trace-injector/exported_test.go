package trainjector

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, handler, "Global handler should have been set")
}

func TestSetGlobalHandler(t *testing.T) {
	myHandler := new(engine.HandlerFunc)

	Init(context.Background())

	SetGlobalHandler(*myHandler)

	handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	myHandlerName := runtime.FuncForPC(reflect.ValueOf(myHandler).Pointer()).Name()
	assert.Exactly(t, handlerName, myHandlerName)
}

func TestGlobalHandler(t *testing.T) {
	Init(context.Background())

	myHandler := GlobalHandler()

	handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	myHandlerName := runtime.FuncForPC(reflect.ValueOf(myHandler).Pointer()).Name()

	assert.Exactly(t, handlerName, myHandlerName)
}
