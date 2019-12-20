package api

import (
	"net/http"
)

var getNodesTests = []HTTPRouteTests{
	{
		name:                "TestGetNodes200",
		store:               &MockChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                getNodesPath,
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusSliceBody },
	},
	{
		name:                "TestGetNodeByID500",
		store:               &ErrorChainRegistry{},
		httpMethod:          http.MethodGet,
		path:                getNodesPath,
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
