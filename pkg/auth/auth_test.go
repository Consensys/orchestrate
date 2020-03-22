// +build unit

package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCtxKey string

const testKey mockCtxKey = "key"

type MockChecker struct {
	key    string
	called int
}

func (c *MockChecker) Check(ctx context.Context) (context.Context, error) {
	c.called++
	if k, ok := ctx.Value(testKey).(string); ok && k == c.key {
		return ctx, nil
	}
	return ctx, fmt.Errorf("invalid auth")
}

func TestCombineCheckers(t *testing.T) {
	checker1 := &MockChecker{key: "key1"}
	checker2 := &MockChecker{key: "key2"}
	checker := CombineCheckers(checker1, checker2)

	ctx := context.WithValue(context.Background(), testKey, "key1")
	_, err := checker.Check(ctx)
	assert.NoError(t, err, "#1 checker should be valid")
	assert.Equal(t, 1, checker1.called, "#1 checker1 should have been called the correct number of times")
	assert.Equal(t, 0, checker2.called, "#1 checker2 should have been called the correct number of times")

	ctx = context.WithValue(context.Background(), testKey, "key2")
	_, err = checker.Check(ctx)
	assert.NoError(t, err, "#2 Auth should be valid")
	assert.Equal(t, 2, checker1.called, "#2 checker1 should have been called the correct number of times")
	assert.Equal(t, 1, checker2.called, "#2 checker2 should have been called the correct number of times")

	ctx = context.WithValue(context.Background(), testKey, "key-invalid")
	_, err = checker.Check(ctx)
	assert.Error(t, err, "#3 Auth should be invalid")
	assert.Equal(t, 3, checker1.called, "#3 checker1 should have been called the correct number of times")
	assert.Equal(t, 2, checker2.called, "#3 checker2 should have been called the correct number of times")
}
