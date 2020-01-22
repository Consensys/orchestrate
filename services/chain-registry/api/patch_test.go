package api

import (
	"encoding/json"
	"net/http"

	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

var patchNodeByNameTests = []HTTPRouteTests{
	{
		name:       "TestPatchNodeByName200",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/testNodeName",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:       "TestPatchNodeByName400WithWrongURL",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/testNodeName",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			data, _ := json.Marshal(apiError{Message: "FF000@chain-registry.store.api: parse test.com: invalid URI for request"})
			return string(data) + "\n"
		},
	},
	{
		name:                "TestPatchNodeByName400WrongBody",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodPatch,
		path:                "/testTenantID/nodes/testNodeName",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInvalidErrorBody },
	},
	{
		name:       "TestPatchNodeByName404",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/notFoundError",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:       "TestPatchNodeByName500",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/testNodeName",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var patchNodeByIDTests = []HTTPRouteTests{
	{
		name:       "TestPatchNodeByIDByID200",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/1",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:       "TestPatchNodeByID400WithWrongURL",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/1",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			data, _ := json.Marshal(apiError{Message: "FF000@chain-registry.store.api: parse test.com: invalid URI for request"})
			return string(data) + "\n"
		},
	},
	{
		name:                "TestPatchNodeByID400WrongBody",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodPatch,
		path:                "/nodes/1",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			data, _ := json.Marshal(apiError{Message: "FF000@chain-registry.store.api: json: unknown field \"unknownField\""})
			return string(data) + "\n"
		},
	},
	{
		name:       "TestPatchNodeByID404",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/0",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:       "TestPatchNodeByID500",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/1",
		body: func() []byte {
			body, _ := json.Marshal(&models.Node{
				URLs: []string{"http://test.com"},
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var patchBlockPositionByIDTests = []HTTPRouteTests{
	{
		name:       "TestPatchBlockPositionByIDByID200",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/1/block-position",
		body: func() []byte {
			body, _ := json.Marshal(&PatchBlockPositionRequest{
				BlockPosition: 10,
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestPatchBlockPositionByID400WrongBody",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodPatch,
		path:                "/nodes/1/block-position",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody: func() string {
			data, _ := json.Marshal(apiError{Message: "FF000@chain-registry.store.api: json: unknown field \"unknownField\""})
			return string(data) + "\n"
		},
	},
	{
		name:       "TestPatchBlockPositionByID404",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/0/block-position",
		body: func() []byte {
			body, _ := json.Marshal(&PatchBlockPositionRequest{
				BlockPosition: 10,
			})
			return body
		},
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:       "TestPatchBlockPositionByID500",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/nodes/1/block-position",
		body: func() []byte {
			body, _ := json.Marshal(&PatchBlockPositionRequest{
				BlockPosition: 10,
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var patchBlockNumberByNameTests = []HTTPRouteTests{
	{
		name:       "TestPatchBlockNumberByName200",
		store:      UseMockChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/testNodeName/block-position",
		body: func() []byte {
			body, _ := json.Marshal(&PatchBlockPositionRequest{
				BlockPosition: 10,
			})
			return body
		},
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestPatchBlockNumberByName400WrongBody",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodPatch,
		path:                "/testTenantID/nodes/testNodeName/block-position",
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInvalidErrorBody },
	},
	{
		name:       "TestPatchBlockNumberByName404",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/notFoundError/block-position",
		body: func() []byte {
			body, _ := json.Marshal(&PatchBlockPositionRequest{
				BlockPosition: 10,
			})
			return body
		},
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:       "TestPatchBlockNumberByName500",
		store:      UseErrorChainRegistry,
		httpMethod: http.MethodPatch,
		path:       "/testTenantID/nodes/testNodeName/block-position",
		body: func() []byte {
			body, _ := json.Marshal(&PatchBlockPositionRequest{
				BlockPosition: 10,
			})
			return body
		},
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
