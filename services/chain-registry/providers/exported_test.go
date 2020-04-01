// +build unit

package providers

import (
	"context"
	"testing"

	"github.com/containous/traefik/v2/pkg/provider/aggregator"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, ProviderAggregator(), "Global ProviderAggregator should have been set")

	var p *aggregator.ProviderAggregator
	SetGlobalProviderAggregator(p)
	assert.Nil(t, ProviderAggregator(), "Global should be reset to nil")

}
