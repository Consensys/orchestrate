package listener

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalListener(), "Global should have been set")

	var l TxListener
	SetGlobalListener(l)
	assert.Nil(t, GlobalListener(), "Global should be reset to nil")
}
