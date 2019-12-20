package api

import (
	"encoding/json"
	"net/http"
	"strings"

	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

var patchNodeByIDTests = []HTTPRouteTests{
	{
		name:       "TestPatchNodeByIDByID200",
		store:      &MockChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(patchNodeByIDPath, "{"+nodeIDPath+"}", "1"),
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
		store:      &MockChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(patchNodeByIDPath, "{"+nodeIDPath+"}", "1"),
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
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodPatch,
		path:                strings.ReplaceAll(patchNodeByIDPath, "{"+nodeIDPath+"}", "1"),
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
		store:      &ErrorChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(patchNodeByIDPath, "{"+nodeIDPath+"}", "0"),
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
		store:      &ErrorChainRegistry{},
		httpMethod: http.MethodPatch,
		path:       strings.ReplaceAll(patchNodeByIDPath, "{"+nodeIDPath+"}", "1"),
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
