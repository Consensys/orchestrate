package api

import (
	"net/http"
)

var deleteNodeByIDTests = []HTTPRouteTests{
	{
		name:                "TestDeleteNodeByID200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/nodes/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestDeleteNodeByID404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/nodes/0",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteNodeByID500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/nodes/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var deleteNodesByNameTests = []HTTPRouteTests{
	{
		name:                "TestDeleteNodeByName200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/testTenantID/nodes/testNodeName",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestDeleteNodeByName404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/testTenantID/nodes/notFoundError",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteNodeByName500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/testTenantID/nodes/testNodeName",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
