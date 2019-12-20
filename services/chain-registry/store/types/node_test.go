package types

import (
	"reflect"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
)

func TestBuildConfiguration(t *testing.T) {

	testSet := []struct {
		nodes          []*Node
		expectedOutput func(*dynamic.Configuration) *dynamic.Configuration
		expectedError  bool
	}{
		{
			[]*Node{
				{ID: 1, Name: "testNode", TenantID: "testTenantId", URLs: []string{"http://testURL1.com", "http://testURL2.com"}},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["testTenantId-testNode"] = &dynamic.Router{EntryPoints: []string{"http"}, Service: "testTenantId-testNode", Rule: "Path(`/1`) || Path(`/testTenantId/testNode`)"}
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
				{ID: 1, Name: "testNode", TenantID: "testTenantId", URLs: []string{"http://testURL1.com", "http://testURL2.com"}},
				{ID: 10, Name: "testNode10", TenantID: "testTenantId", URLs: []string{"http://testURL10.com", "http://testURL20.com"}},
			},
			func(c *dynamic.Configuration) *dynamic.Configuration {
				c.HTTP.Routers["testTenantId-testNode"] = &dynamic.Router{EntryPoints: []string{"http"}, Service: "testTenantId-testNode", Rule: "Path(`/1`) || Path(`/testTenantId/testNode`)"}
				c.HTTP.Routers["testTenantId-testNode10"] = &dynamic.Router{EntryPoints: []string{"http"}, Service: "testTenantId-testNode10", Rule: "Path(`/10`) || Path(`/testTenantId/testNode10`)"}
				c.HTTP.Services["testTenantId-testNode"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{
					{Scheme: "http", URL: "http//testURL1.com"},
					{Scheme: "http", URL: "http://testURL2.com"},
				}}}
				c.HTTP.Services["testTenantId-testNode10"] = &dynamic.Service{LoadBalancer: &dynamic.ServersLoadBalancer{Servers: []dynamic.Server{
					{Scheme: "http", URL: "http//testURL10.com"},
					{Scheme: "http", URL: "http://testURL20.com"},
				}}}
				return c
			},
			false,
		},
		{
			[]*Node{
				{ID: 1, Name: "testNode", TenantID: "testTenantId", URLs: []string{":/*testURL1@com", "http://testURL2.com"}},
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
		eq := reflect.DeepEqual(expectedOutput.HTTP.Routers, output.HTTP.Routers)
		if !eq {
			t.Errorf("Chain-registry - Store (%d/%d): expected %v but got %v", i+1, len(testSet), expectedOutput, output)
		}

	}

}
