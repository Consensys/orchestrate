package api

import (
	"encoding/json"
	"net/http"

	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

var postNodeTests = []HTTPRouteTests{
	{
		name:       "TestPostNode200",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/nodes",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				Name:                    "testName",
				URLs:                    []string{"http://test.com"},
				ListenerDepth:           1,
				ListenerBlockPosition:   1,
				ListenerFromBlock:       1,
				ListenerBackOffDuration: "1s",
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return "{\"id\":\"1\"}\n" },
	},
	{
		name:       "TestPostNode400WithTwiceSameURL",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/nodes",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com", "http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			return "{\"message\":\"FF000@chain-registry.store.api: cannot have twice the same url - got at least two times http://test.com\"}\n"
		},
	},
	{
		name:       "TestPostNode400WrongURL",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/nodes",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInvalidURLErrorBody },
	},
	{
		name:                "TestPostNode400WrongBody",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodPost,
		path:                "/testTenantID/nodes",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInvalidErrorBody },
	},
	{
		name:       "TestPostNode500",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/nodes",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				Name:                    "testName",
				URLs:                    []string{"http://test.com"},
				ListenerDepth:           1,
				ListenerBlockPosition:   1,
				ListenerFromBlock:       1,
				ListenerBackOffDuration: "1s",
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
