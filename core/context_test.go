package core

import (
	"errors"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

var testKey = "test"

func newHandler(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), s)
	}
}

var errTest = errors.New("Test Error")

func newErrorHandler(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), s)
		ctx.Error(errTest)
	}
}

func newAborter(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), s)
		ctx.AbortWithError(errTest)
	}
}

func newMiddleware(s string, t *testing.T) HandlerFunc {
	return func(ctx *Context) {
		sA := fmt.Sprintf("%v-before", s)
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), sA)

		ctx.Next()

		sB := fmt.Sprintf("%v-after", s)
		t.Logf("At %v, index=%v", s, ctx.index)
		ctx.Keys[testKey] = append(ctx.Keys[testKey].([]string), sB)
	}
}

func TestNext(t *testing.T) {
	ctx := NewContext()

	var (
		hA   = newHandler("hA", t)
		mA   = newMiddleware("mA", t)
		hErr = newErrorHandler("err", t)
		mB   = newMiddleware("mB", t)
		hB   = newHandler("hB", t)
		a    = newAborter("abort", t)
		hC   = newHandler("hC", t)
		mC   = newMiddleware("mC", t)
	)
	// Initialize context
	ctx.Prepare([]HandlerFunc{hA, mA, hErr, mB, hB, a, hC, mC}, nil, nil)
	ctx.Keys[testKey] = []string{}

	// Handle context
	ctx.Next()

	res := ctx.Keys[testKey].([]string)
	expected := []string{"hA", "mA-before", "err", "mB-before", "hB", "abort", "mB-after", "mA-after"}

	assert.Equal(t, expected, res, "Call order on handlers should be correct")
	assert.Len(t, ctx.T.Errors, 2, "Error count should be correct")
}

func TestCtxError(t *testing.T) {
	err := fmt.Errorf("Test Error")

	ctx := NewContext()
	ctx.Error(err)

	assert.Len(t, ctx.T.Errors, 1, "Error count should be correct")

	err = &common.Error{Message: "Test Error", Type: 5}
	ctx.Error(err)

	assert.Len(t, ctx.T.Errors, 2, "Error count should be correct")
	assert.Equal(t, `2 error(s): ["Error #0: Test Error" "Error #5: Test Error"]`, ctx.T.Error(), "Error message should be correct")
}

func TestLogger(t *testing.T) {
	logHandler := func(ctx *Context) { ctx.Logger.Info("Test") }

	ctx := NewContext()
	ctx.Prepare([]HandlerFunc{logHandler}, log.NewEntry(log.StandardLogger()), nil)

	ctx.Next()
}
