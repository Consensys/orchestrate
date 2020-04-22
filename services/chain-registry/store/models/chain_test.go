// +build unit

package models

import (
	"testing"

	genuuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func TestIsValidChain(t *testing.T) {
	testSet := []struct {
		chain   *Chain
		isValid bool
	}{
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			true,
		},
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
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
			UUID:                    genuuid.NewV4().String(),
			TenantID:                "test",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			URLs:                    []string{},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
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
			UUID:                  genuuid.NewV4().String(),
			Name:                  "test",
			TenantID:              "test",
			URLs:                  []string{"http://test.com", "http://test.net"},
			ListenerDepth:         &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:  &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock: &(&struct{ x uint64 }{1}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
			URLs:                    []string{"http://test.com", "http://test.net"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"200"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
			URLs:                    []string{"%%%"},
			ListenerDepth:           &(&struct{ x uint64 }{1}).x,
			ListenerCurrentBlock:    &(&struct{ x uint64 }{1}).x,
			ListenerStartingBlock:   &(&struct{ x uint64 }{1}).x,
			ListenerBackOffDuration: &(&struct{ x string }{"2s"}).x,
		},
			false,
		},
		{&Chain{
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
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
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
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
			UUID:                    genuuid.NewV4().String(),
			Name:                    "test",
			TenantID:                "test",
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
	assert.Equal(t, multitenancy.DefaultTenantIDName, chain.TenantID, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerDepth, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerCurrentBlock, "Should not be empty")
	assert.Equal(t, uint64(0), *chain.ListenerStartingBlock, "Should not be empty")
	assert.Equal(t, "1s", *chain.ListenerBackOffDuration, "Should not be empty")
}
