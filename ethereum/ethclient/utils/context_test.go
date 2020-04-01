// +build unit

package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryNotFoundError(t *testing.T) {
	ctx := RetryNotFoundError(context.Background(), true)
	assert.True(t, ShouldRetryNotFoundError(ctx), "Flag should have been set")
}
