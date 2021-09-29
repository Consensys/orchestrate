// +build unit

package dynamic

import (
	"context"
	"testing"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	builder := NewBuilder()
	srv := &dynamic.Service{
		HealthCheck: &dynamic.HealthCheck{},
	}

	h, err := builder.Build(context.Background(), "", srv, nil)
	require.NoError(t, err)
	assert.NotNil(t, h)
}
