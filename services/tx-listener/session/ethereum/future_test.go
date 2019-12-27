package ethereum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuture(t *testing.T) {
	// Test with no error
	future := NewFuture(func() (interface{}, error) { return "test", nil })
	res := <-future.Result()
	assert.Equal(t, "test", res.(string), "Result should be correct")
	future.Close()

	// Test with error
	future = NewFuture(func() (interface{}, error) { return nil, fmt.Errorf("test-error") })
	err := <-future.Err()
	assert.Equal(t, "test-error", err.Error(), "Result should be correct")
	future.Close()
}
