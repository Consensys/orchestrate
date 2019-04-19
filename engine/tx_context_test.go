package engine

import (
	"context"
	"errors"
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

var testKey = "test"

func newHandler(s string, t *testing.T) HandlerFunc {
	return func(txctx *TxContext) {
		t.Logf("At %v, index=%v", s, txctx.index)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), s))
	}
}

var errTest = errors.New("Test Error")

func newErrorHandler(s string, t *testing.T) HandlerFunc {
	return func(txctx *TxContext) {
		t.Logf("At %v, index=%v", s, txctx.index)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), s))
		txctx.Error(errTest)
	}
}

func newAborter(s string, t *testing.T) HandlerFunc {
	return func(txctx *TxContext) {
		t.Logf("At %v, index=%v", s, txctx.index)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), s))
		txctx.AbortWithError(errTest)
	}
}

func newMiddleware(s string, t *testing.T) HandlerFunc {
	return func(txctx *TxContext) {
		sA := fmt.Sprintf("%v-before", s)
		t.Logf("At %v, index=%v", s, txctx.index)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), sA))

		txctx.Next()

		sB := fmt.Sprintf("%v-after", s)
		t.Logf("At %v, index=%v", s, txctx.index)
		txctx.Set(testKey, append(txctx.Get(testKey).([]string), sB))
	}
}

func TestNext(t *testing.T) {
	txctx := NewTxContext()

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
	txctx.Prepare([]HandlerFunc{hA, mA, hErr, mB, hB, a, hC, mC}, nil, nil)
	txctx.Set(testKey, []string{})

	// Handle context
	txctx.Next()

	res := txctx.Get(testKey).([]string)
	expected := []string{"hA", "mA-before", "err", "mB-before", "hB", "abort", "mB-after", "mA-after"}

	assert.Equal(t, expected, res, "Call order on handlers should be correct")
	assert.Len(t, txctx.Envelope.Errors, 2, "Error count should be correct")
}

func TestCtxError(t *testing.T) {
	err := fmt.Errorf("Test Error")

	txctx := NewTxContext()
	txctx.Error(err)

	assert.Len(t, txctx.Envelope.Errors, 1, "Error count should be correct")

	err = &common.Error{Message: "Test Error", Type: 5}
	txctx.Error(err)

	assert.Len(t, txctx.Envelope.Errors, 2, "Error count should be correct")
	assert.Equal(t, `2 error(s): ["Error #0: Test Error" "Error #5: Test Error"]`, txctx.Envelope.Error(), "Error message should be correct")
}

func TestLogger(t *testing.T) {
	logHandler := func(txctx *TxContext) { txctx.Logger.Info("Test") }
	txctx := NewTxContext()
	txctx.Prepare([]HandlerFunc{logHandler}, log.NewEntry(log.StandardLogger()), nil)
	txctx.Next()
}

type testingKey string

func TestWithContext(t *testing.T) {
	logHandler := func(txctx *TxContext) { txctx.Logger.Info("Test") }
	txctx := NewTxContext()
	txctx.Prepare([]HandlerFunc{logHandler}, log.NewEntry(log.StandardLogger()), nil)

	// Update go context attached to TxContext
	txctx.WithContext(context.WithValue(context.Background(), testingKey("test-key"), "test-value"))

	// Check if go-context has been properly attached
	assert.Equal(t, "test-value", txctx.Context().Value(testingKey("test-key")).(string), "Go context should have been attached")
}
