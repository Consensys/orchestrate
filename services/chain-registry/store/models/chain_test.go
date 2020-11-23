// +build unit

package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func TestIsValidChain(t *testing.T) {
	testSet := []struct {
		chain   *Chain
		isValid bool
	}{
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			true,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
			PrivateTxManagers: []*PrivateTxManagerModel{
				&PrivateTxManagerModel{
					URL:  "http://test.com/private",
					Type: utils.TesseraChainType,
				},
			},
		},
			true,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			ChainID:                 "888",
			URLs:                    []string{},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                  uuid.Must(uuid.NewV4()).String(),
			Name:                  "test",
			TenantID:              "test",
			ChainID:               "888",
			URLs:                  []string{"http://test.com", "http://test.net"},
			ListenerDepth:         &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:  &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"200"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"%%%"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
			PrivateTxManagers: []*PrivateTxManagerModel{
				&PrivateTxManagerModel{
					URL:  "^&$$",
					Type: utils.TesseraChainType,
				},
			},
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
			PrivateTxManagers: []*PrivateTxManagerModel{
				&PrivateTxManagerModel{
					URL:  "http://test.com/private",
					Type: "TTTT",
				},
			},
		},
			false,
		},
		{&Chain{
			UUID:                    uuid.Must(uuid.NewV4()).String(),
			Name:                    "test",
			TenantID:                "test",
			ChainID:                 "888",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
			PrivateTxManagers: []*PrivateTxManagerModel{
				&PrivateTxManagerModel{
					URL: "http://test.com/private",
				},
			},
		},
			false,
		},
	}

	for _, test := range testSet {
		assert.Equal(t, test.isValid, test.chain.IsValid())
	}
}

func TestChain_SetDefaultIfNil(t *testing.T) {
	chain := Chain{}
	chain.SetDefault()

	assert.NotNil(t, chain.UUID, "Should not be empty")
	assert.Equal(t, multitenancy.DefaultTenant, chain.TenantID, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerDepth, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerCurrentBlock, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerStartingBlock, "Should not be empty")
	assert.Equal(t, "1s", *chain.ListenerBackOffDuration, "Should not be empty")
}
