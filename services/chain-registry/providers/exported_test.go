package providers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, ProviderAggregator(), "Global ProviderAggregator should have been set")
}
