package chains

import (
	"encoding/json"
	"net/http"
)

var postChainTests = []HTTPRouteTests{
	{
		name:       "TestPostChain200",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/chains",
		body: func() []byte {
			listenerDepth := uint64(1)
			listenerBlockPosition := int64(1)
			listenerBackOffDuration := "1s"

			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com"},
				Listener: &Listener{
					Depth:           &listenerDepth,
					BlockPosition:   &listenerBlockPosition,
					BackOffDuration: &listenerBackOffDuration,
				},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return "{\"uuid\":\"1\"}\n" },
	},
	{
		name:       "TestPostChain400WithTwiceSameURL",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com", "http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			return expectedInvalidBodyError
		},
	},
	{
		name:       "TestPostChain400WrongURL",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/chains",
		body: func() []byte {
			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInvalidBodyError },
	},
	{
		name:                "TestPostChain400WrongBody",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodPost,
		path:                "/testTenantID/chains",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedUnknownBodyError },
	},
	{
		name:       "TestPostChain500",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPost,
		path:       "/testTenantID/chains",
		body: func() []byte {
			listenerDepth := uint64(1)
			listenerBlockPosition := int64(1)
			listenerBackOffDuration := "1s"

			body, _ := json.Marshal(&PostRequest{
				Name: "testName",
				URLs: []string{"http://test.com"},
				Listener: &Listener{
					Depth:           &listenerDepth,
					BlockPosition:   &listenerBlockPosition,
					BackOffDuration: &listenerBackOffDuration,
				},
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
