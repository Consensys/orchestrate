package types

import (
	"math"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestIsValidNode(t *testing.T) {
	testSet := []struct {
		node    *Node
		isValid bool
	}{
		{&Node{
			Name:                    "test",
			TenantID:                "test",
			URLs:                    []string{"test.com", "test.net"},
			ListenerDepth:           1,
			ListenerBlockPosition:   1,
			ListenerFromBlock:       1,
			ListenerBackOffDuration: "2s",
		},
			true,
		},
		{&Node{
			TenantID:                "test",
			URLs:                    []string{"test.com", "test.net"},
			ListenerDepth:           1,
			ListenerBlockPosition:   1,
			ListenerFromBlock:       1,
			ListenerBackOffDuration: "2s",
		},
			false,
		},
		{&Node{
			Name:                    "test",
			URLs:                    []string{"test.com", "test.net"},
			ListenerDepth:           1,
			ListenerBlockPosition:   1,
			ListenerFromBlock:       1,
			ListenerBackOffDuration: "2s",
		},
			false,
		},
		{&Node{
			Name:                    "test",
			TenantID:                "test",
			ListenerDepth:           1,
			ListenerBlockPosition:   1,
			ListenerFromBlock:       1,
			ListenerBackOffDuration: "2s",
		},
			false,
		},
		{&Node{
			Name:                  "test",
			TenantID:              "test",
			URLs:                  []string{"test.com", "test.net"},
			ListenerDepth:         1,
			ListenerBlockPosition: 1,
			ListenerFromBlock:     1,
		},
			false,
		},
	}

	for _, test := range testSet {
		assert.Equal(t, test.node.IsValid(), test.isValid)
	}
}

func TestBuildConfiguration(t *testing.T) {

	testSet := []struct {
		nodes          []*Node
		expectedOutput func(*dynamic.Configuration) *dynamic.Configuration
		expectedError  bool
	}{
		{
			[]*Node{
				{ID: "0d60a85e-0b90-4482-a14c-108aea2557aa", Name: "testNode", TenantID: "testTenantId", URLs: []string{"http://testURL1.com", "http://testURL2.com"}},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["testTenantId-testNode"] = &dynamic.Router{EntryPoints: []string{"orchestrate"}, Priority: math.MaxInt32, Service: "testTenantId-testNode", Rule: "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`) || Path(`/{tenantID:testTenantId}/testNode`)", Middlewares: []string{"orchestrate-auth"}}
				c.HTTP.Services["testTenantId-testNode"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{
					{Scheme: "http", URL: "http//testURL1.com"},
					{Scheme: "http", URL: "http://testURL2.com"},
				}}}
				return c
			},
			false,
		},
		{
			[]*Node{
				{ID: "0d60a85e-0b90-4482-a14c-108aea2557aa", Name: "testNode", TenantID: "testTenantId", URLs: []string{"http://testURL1.com", "http://testURL2.com"}},
				{ID: "39240e9f-ae09-4e95-9fd0-a712035c8ad7", Name: "testNode2", TenantID: "testTenantId", URLs: []string{"http://testURL10.com", "http://testURL20.com"}},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["testTenantId-testNode"] = &dynamic.Router{EntryPoints: []string{"orchestrate"}, Priority: math.MaxInt32, Service: "testTenantId-testNode", Rule: "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`) || Path(`/{tenantID:testTenantId}/testNode`)", Middlewares: []string{"orchestrate-auth"}}
				c.HTTP.Routers["testTenantId-testNode2"] = &dynamic.Router{EntryPoints: []string{"orchestrate"}, Priority: math.MaxInt32, Service: "testTenantId-testNode2", Rule: "Path(`/39240e9f-ae09-4e95-9fd0-a712035c8ad7`) || Path(`/{tenantID:testTenantId}/testNode2`)", Middlewares: []string{"orchestrate-auth"}}
				c.HTTP.Services["testTenantId-testNode"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{
					{Scheme: "http", URL: "http//testURL1.com"},
					{Scheme: "http", URL: "http://testURL2.com"},
				}}}
				c.HTTP.Services["testTenantId-testNode2"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{
					{Scheme: "http", URL: "http//testURL10.com"},
					{Scheme: "http", URL: "http://testURL20.com"},
				}}}
				return c
			},
			false,
		},
		{
			[]*Node{
				{ID: "0d60a85e-0b90-4482-a14c-108aea2557aa", Name: "testNode", TenantID: "testTenantId", URLs: []string{":/*testURL1@com", "http://testURL2.com"}},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				return c
			},
			true,
		},
	}

	for i, test := range testSet {
		output, err := BuildConfiguration(test.nodes)
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
