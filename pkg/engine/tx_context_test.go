package engine

import (
	"context"
	"errors"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/error"
)

var testKey = "test"

func newPipeline(s string) HandlerFunc {
	return func(txctx *TxContext) {
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), s))
	}
}

var errTest = errors.New("test Error")

func newErrorHandler(s string) HandlerFunc {
	return func(txctx *TxContext) {
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), s))
		_ = txctx.Error(errTest)
	}
}

func newAborter() HandlerFunc {
	return func(txctx *TxContext) {
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), "abort"))
		_ = txctx.AbortWithError(errTest)
	}
}

func newMiddleware(s string) HandlerFunc {
	return func(txctx *TxContext) {
		sA := fmt.Sprintf("%v-before", s)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), sA))

		txctx.Next()

		sB := fmt.Sprintf("%v-after", s)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), sB))
	}
}

func TestApplyHandlers(t *testing.T) {
	txctx := NewTxContext()

	var (
		pA   = newPipeline("pA")
		mA   = newMiddleware("mA")
		hErr = newErrorHandler("err")
		mB   = newMiddleware("mB")
		pB   = newPipeline("pB")
		a    = newAborter()
		pC   = newPipeline("pC")
		mC   = newMiddleware("mC")
	)
	// Initialize context
	txctx.Prepare(nil, nil)
	txctx.Set(testKey, []string{})

	// Handle context
	txctx.applyHandlers([]HandlerFunc{pA, mA, hErr, mB, pB, a, pC, mC}...)

	res := txctx.Get(testKey).([]string)
	expected := []string{"pA", "mA-before", "err", "mB-before", "pB", "abort", "mB-after", "mA-after"}

	assert.Equal(t, expected, res, "Call order on handlers should be correct")
	assert.Len(t, txctx.Envelope.Errors, 2, "Error count should be correct")
}

func TestCtxError(t *testing.T) {
	e := fmt.Errorf("test Error")

	txctx := NewTxContext()
	_ = txctx.Error(e).ExtendComponent("bar")

	assert.Len(t, txctx.Envelope.Errors, 1, "Error count should be correct")

	e = ierror.New(5, "test Error").ExtendComponent("foo")
	_ = txctx.Error(e)

	assert.Len(t, txctx.Envelope.GetErrors(), 2, "Error count should be correct")
	assert.Equal(t, `["FF000@bar: test Error" "00005@foo: test Error"]`, txctx.Envelope.Error(), "Error message should be correct")
}

func TestLogger(t *testing.T) {
	logHandler := func(txctx *TxContext) { txctx.Logger.Info("Test") }
	txctx := NewTxContext()
	txctx.Prepare(log.NewEntry(log.StandardLogger()), nil)
	txctx.applyHandlers(logHandler)
}

type testingKey string

func TestWithContext(t *testing.T) {
	txctx := NewTxContext()
	txctx.Prepare(log.NewEntry(log.StandardLogger()), nil)

	// Update go context attached to TxContext
	txctx.WithContext(context.WithValue(context.Background(), testingKey("test-key"), "test-value"))

	// Check if go-context has been properly attached
	assert.Equal(t, "test-value", txctx.Context().Value(testingKey("test-key")).(string), "Go context should have been attached")
}

func TestCombineHandlers(t *testing.T) {
	var (
		pA   = newPipeline("pA")
		mA   = newMiddleware("mA")
		hErr = newErrorHandler("err")
		mB   = newMiddleware("mB")
		pB   = newPipeline("pB")
		a    = newAborter()
		pC   = newPipeline("pC")
		mC   = newMiddleware("mC")
	)

	// Create combined handler
	combinedHandler := CombineHandlers([]HandlerFunc{pA, mA, hErr, mB, pB, a, pC, mC}...)

	// Initialize context and apply combinedHandler
	txctx := NewTxContext()
	txctx.Prepare(nil, nil)
	txctx.Set(testKey, []string{})
	txctx.applyHandlers(combinedHandler)

	res := txctx.Get(testKey).([]string)
	expected := []string{"pA", "mA-before", "err", "mB-before", "pB", "abort", "mB-after", "mA-after"}

	assert.Equal(t, expected, res, "Call order on handlers should be correct")
	assert.Len(t, txctx.Envelope.Errors, 2, "Error count should be correct")
}

func TestCombineHandlersNested(t *testing.T) {
	var (
		pA   = newPipeline("pA")
		mA   = newMiddleware("mA")
		hErr = newErrorHandler("err")
		mB   = newMiddleware("mB")
		pB   = newPipeline("pB")
		a    = newAborter()
		pC   = newPipeline("pC")
		mC   = newMiddleware("mC")
	)

	// Create combined handler
	combinedHandler := CombineHandlers([]HandlerFunc{pA, mA, hErr, mB, pB, a, pC, mC}...)

	// Initialize context and apply combinedHandler
	txctx := NewTxContext()
	txctx.Prepare(nil, nil)
	txctx.Set(testKey, []string{})
	txctx.applyHandlers(combinedHandler)

	res := txctx.Get(testKey).([]string)
	expected := []string{"pA", "mA-before", "err", "mB-before", "pB", "abort", "mB-after", "mA-after"}

	assert.Equal(t, expected, res, "Call order on handlers should be correct")
	assert.Len(t, txctx.Envelope.Errors, 2, "Error count should be correct")
}

func TestForkedHandler(t *testing.T) {
	var (
		pA = newPipeline("pA")
		pB = newPipeline("pB")
		pC = newPipeline("pC")
		pD = newPipeline("pD")
		pE = newPipeline("pE")
		mA = newMiddleware("mA")
		mB = newMiddleware("mB")
		mC = newMiddleware("mC")
		a  = newAborter()
	)

	// Create combined handler
	handler1 := CombineHandlers([]HandlerFunc{mB, pB}...)
	handler2 := CombineHandlers([]HandlerFunc{pC, mC, a, pD}...)

	// Declare a fork
	forkedHandler := func(txctx *TxContext) {
		switch txctx.Get("fork").(string) {
		case "1":
			handler1(txctx)
		case "2":
			handler2(txctx)
		}
	}

	// Initialize context and test on fork 1
	txctx := NewTxContext()
	txctx.Prepare(nil, nil)
	txctx.Set(testKey, []string{})
	txctx.Set("fork", "1")
	txctx.applyHandlers(mA, pA, forkedHandler, pE)
	res := txctx.Get(testKey).([]string)
	expected := []string{"mA-before", "pA", "mB-before", "pB", "mB-after", "pE", "mA-after"}

	assert.Equal(t, expected, res, "Called handlers on fork 1 should be correctly ordered")

	// Initialize context and test on fork 2
	txctx = NewTxContext()
	txctx.Prepare(nil, nil)
	txctx.Set(testKey, []string{})
	txctx.Set("fork", "2")
	txctx.applyHandlers(mA, pA, forkedHandler, pE)
	res = txctx.Get(testKey).([]string)
	expected = []string{"mA-before", "pA", "pC", "mC-before", "abort", "mC-after", "mA-after"}

	assert.Equal(t, expected, res, "Called handlers on fork 2 should be correctly ordered")
}
