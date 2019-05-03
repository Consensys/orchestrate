package nonce

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, handler, "Global handler should have been set")
}
