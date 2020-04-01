// +build unit

package types

import (
	"math"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
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

func TestBuildConfiguration(t *testing.T) {
	testSet := []struct {
		chains         []*Chain
		expectedOutput func(*dynamic.Configuration) *dynamic.Configuration
		expectedError  bool
	}{
		{
			[]*Chain{
				{
					UUID:     "0d60a85e-0b90-4482-a14c-108aea2557aa",
					Name:     "testChain",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL1.com",
						"http://testURL2.com",
					},
				},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["testTenantId-testChain"] = &dynamic.Router{
					EntryPoints: []string{"orchestrate"},
					Priority:    math.MaxInt32,
					Service:     "testTenantId-testChain",
					Rule:        "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`)",
					Middlewares: []string{
						"orchestrate-auth",
						"strip-path@internal",
						"orchestrate-ratelimit",
					},
				}
				c.HTTP.Services["testTenantId-testChain"] = &dynamic.Service{
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{Scheme: "http", URL: "http//testURL1.com"},
							{Scheme: "http", URL: "http://testURL2.com"},
						}}}
				return c
			},
			false,
		},
		{
			[]*Chain{
				{
					UUID:     "0d60a85e-0b90-4482-a14c-108aea2557aa",
					Name:     "testChain",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL1.com",
						"http://testURL2.com",
					},
				},
				{
					UUID:     "39240e9f-ae09-4e95-9fd0-a712035c8ad7",
					Name:     "testChain2",
					TenantID: "testTenantId",
					URLs: []string{
						"http://testURL10.com",
						"http://testURL20.com",
					},
				},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["testTenantId-testChain"] = &dynamic.Router{
					EntryPoints: []string{"orchestrate"},
					Priority:    math.MaxInt32,
					Service:     "testTenantId-testChain",
					Rule:        "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`)",
					Middlewares: []string{
						"orchestrate-auth",
						"strip-path@internal",
						"orchestrate-ratelimit",
					},
				}
				c.HTTP.Routers["testTenantId-testChain2"] = &dynamic.Router{
					EntryPoints: []string{"orchestrate"},
					Priority:    math.MaxInt32,
					Service:     "testTenantId-testChain2",
					Rule:        "Path(`/39240e9f-ae09-4e95-9fd0-a712035c8ad7`)",
					Middlewares: []string{
						"orchestrate-auth",
						"strip-path@internal",
						"orchestrate-ratelimit",
					},
				}
				c.HTTP.Services["testTenantId-testChain"] = &dynamic.Service{
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{Scheme: "http", URL: "http//testURL1.com"},
							{Scheme: "http", URL: "http://testURL2.com"},
						}}}
				c.HTTP.Services["testTenantId-testChain2"] = &dynamic.Service{
					LoadBalancer: &dynamic.ServersLoadBalancer{
						Servers: []dynamic.Server{
							{Scheme: "http", URL: "http//testURL10.com"},
							{Scheme: "http", URL: "http://testURL20.com"},
						}}}
				return c
			},
			false,
		},
		{
			[]*Chain{
				{UUID: "0d60a85e-0b90-4482-a14c-108aea2557aa", Name: "testChain", TenantID: "testTenantId", URLs: []string{":/*testURL1@com", "http://testURL2.com"}},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				return c
			},
			true,
		},
	}

	for i, test := range testSet {
		output, err := BuildConfiguration(test.chains)
		if (err == nil && test.expectedError) || (err != nil && !test.expectedError) {
			t.Errorf("Chain-registry - Store: Expecting the following error %v but got %v", test.expectedError, err)
			return
		} else if err != nil && test.expectedError {
			continue
		}

		expectedOutput := test.expectedOutput(NewConfig())

		t.Log()
		assert.Equal(t, expectedOutput.HTTP.Routers, output.HTTP.Routers, "Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), expectedOutput, output)
	}
}
