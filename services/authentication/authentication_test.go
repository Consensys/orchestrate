package authentication

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockAuthCtxKey string

const key mockAuthCtxKey = "key"

type MockAuth struct {
	key    string
	called int
}

func (a *MockAuth) Check(ctx context.Context) (context.Context, error) {
	a.called++
	if k, ok := ctx.Value(key).(string); ok && k == a.key {
		return ctx, nil
	}
	return ctx, fmt.Errorf("invalid auth")
}

func TestCombineAuth(t *testing.T) {
	auth1 := &MockAuth{key: "key1"}
	auth2 := &MockAuth{key: "key2"}
	auth := CombineAuth(auth1, auth2)

	ctx := context.WithValue(context.Background(), key, "key1")
	_, err := auth.Check(ctx)
	assert.NoError(t, err, "#1 Auth should be valid")
	assert.Equal(t, 1, auth1.called, "#1 Auth1 should have been called the correct number of times")
	assert.Equal(t, 0, auth2.called, "#1 Auth2 should have been called the correct number of times")

	ctx = context.WithValue(context.Background(), key, "key2")
	_, err = auth.Check(ctx)
	assert.NoError(t, err, "#2 Auth should be valid")
	assert.Equal(t, 2, auth1.called, "#2 Auth1 should have been called the correct number of times")
	assert.Equal(t, 1, auth2.called, "#2 Auth2 should have been called the correct number of times")

	ctx = context.WithValue(context.Background(), key, "key-invalid")
	_, err = auth.Check(ctx)
	assert.Error(t, err, "#3 Auth should be invalid")
	assert.Equal(t, 3, auth1.called, "#3 Auth1 should have been called the correct number of times")
	assert.Equal(t, 2, auth2.called, "#3 Auth2 should have been called the correct number of times")
}
