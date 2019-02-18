package types

import (
	"errors"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
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
	if len(res) != len(expected) {
		t.Errorf("Context: expected %v execution but got %v", len(expected), len(res))
	}

	for i, s := range expected {
		if s != res[i] {
			t.Errorf("Context: expected %q at %v but got %q", s, i, res[i])
		}
	}

	if len(ctx.T.Errors) != 2 {
		t.Errorf("Context: expected 2 errors but got %v", len(ctx.T.Errors))
	}
}

func TestLogger(t *testing.T) {
	logHandler := func(ctx *Context) { ctx.Logger.Info("Test") }

	ctx := NewContext()
	ctx.Prepare([]HandlerFunc{logHandler}, log.NewEntry(log.StandardLogger()), nil)

	ctx.Next()
}
