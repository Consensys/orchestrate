package api

import (
	"encoding/json"
	"net/http"
	"strings"

	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

var patchNodeByNameTests = []HTTPRouteTests{
	{
		name:       "TestPatchNodeByName200",
		store:      &MockChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(strings.ReplaceAll(patchNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "nodeName"),
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
		store:      &MockChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(strings.ReplaceAll(patchNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "nodeName"),
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
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodPatch,
		path:                strings.ReplaceAll(strings.ReplaceAll(patchNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "nodeName"),
		body:                func() []byte { return []byte(`{"unknownField":"error"}`) },
		expectedStatusCode:  http.StatusBadRequest,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInvalidErrorBody },
	},
	{
		name:       "TestPatchNodeByName404",
		store:      &ErrorChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(strings.ReplaceAll(patchNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "notFoundError"),
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
		store:      &ErrorChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(strings.ReplaceAll(patchNodeByNamePath, "{"+tenantIDPath+"}", "testTenantID"), "{"+nodeNamePath+"}", "testNodeName"),
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
