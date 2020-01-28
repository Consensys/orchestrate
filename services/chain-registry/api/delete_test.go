package api

import (
	"net/http"
)

var deleteChainByUUIDTests = []HTTPRouteTests{
	{
		name:                "TestDeleteChainByUUID200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/chains/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestDeleteChainByUUID404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/chains/0",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteChainByUUID500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/chains/1",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}

var deleteChainsByNameTests = []HTTPRouteTests{
	{
		name:                "TestDeleteChainByName200",
		store:               UseMockChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/testTenantID/chains/testChainName",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusOK,
		expectedContentType: expectedSuccessStatusContentType,
		expectedBody:        func() string { return expectedSuccessStatusBody },
	},
	{
		name:                "TestDeleteChainByName404",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/testTenantID/chains/notFoundError",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusNotFound,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedNotFoundErrorBody },
	},
	{
		name:                "TestDeleteChainByName500",
		store:               UseErrorChainRegistry,
		httpMethod:          http.MethodDelete,
		path:                "/testTenantID/chains/testChainName",
		body:                func() []byte { return nil },
		expectedStatusCode:  http.StatusInternalServerError,
		expectedContentType: expectedErrorStatusContentType,
		expectedBody:        func() string { return expectedInternalServerErrorBody },
	},
}
