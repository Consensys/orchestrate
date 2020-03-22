// +build unit

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

func TestIsValidChain(t *testing.T) {
	testSet := []struct {
		chain   *Chain
		isValid bool
	}{
		{&Chain{
			Name:                    "test",
			TenantID:                "test",
			URLs:                    []string{"test.com", "test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			true,
		},
		{&Chain{
			TenantID:                "test",
			URLs:                    []string{"test.com", "test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			Name:                    "test",
			URLs:                    []string{"test.com", "test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			Name:                    "test",
			TenantID:                "test",
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			Name:                  "test",
			TenantID:              "test",
			URLs:                  []string{"test.com", "test.net"},
			ListenerDepth:         &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:  &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
		},
			false,
		},
	}

	for _, test := range testSet {
		assert.Equal(t, test.chain.IsValid(), test.isValid)
	}
}

func TestChain_SetDefaultIfNil(t *testing.T) {
	chain := Chain{}
	chain.SetDefault()

	assert.NotNil(t, chain.UUID, "Should not be empty")
	assert.Equal(t, multitenancy.DefaultTenantIDName, chain.TenantID, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerDepth, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerCurrentBlock, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerStartingBlock, "Should not be empty")
	assert.Equal(t, "1s", *chain.ListenerBackOffDuration, "Should not be empty")
}
