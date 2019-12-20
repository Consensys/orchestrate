package types

import (
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/stretchr/testify/assert"
)

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
				c.HTTP.Routers["testTenantId-testNode"] = &dynamic.Router{EntryPoints: []string{"http"}, Service: "testTenantId-testNode", Rule: "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`) || Path(`/testTenantId/testNode`)"}
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
				c.HTTP.Routers["testTenantId-testNode"] = &dynamic.Router{EntryPoints: []string{"http"}, Service: "testTenantId-testNode", Rule: "Path(`/0d60a85e-0b90-4482-a14c-108aea2557aa`) || Path(`/testTenantId/testNode`)"}
				c.HTTP.Routers["testTenantId-testNode2"] = &dynamic.Router{EntryPoints: []string{"http"}, Service: "testTenantId-testNode2", Rule: "Path(`/39240e9f-ae09-4e95-9fd0-a712035c8ad7`) || Path(`/testTenantId/testNode2`)"}
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
		assert.Equalf(t, expectedOutput.HTTP.Routers, output.HTTP.Routers, "Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), expectedOutput, output)
	}
}
