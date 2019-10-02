package grpcclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

func TestDialContext(t *testing.T) {
	conn, err := DialContext(context.Background(), "unknown-target")
	assert.NotNil(t, err, "Dial should error")
	assert.True(t, errors.IsConnectionError(err), "Error should be a gRPC connection error")
	assert.Nil(t, conn, "Connection should be nil")
}
